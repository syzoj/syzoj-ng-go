package binlog

import (
	"bytes"
	"context"
	"encoding/binary"
	"io"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/siddontang/go-mysql/canal"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/replication"
	"github.com/sirupsen/logrus"
)

var log = logrus.StandardLogger()

type Subscriber struct {
	canal   *canal.Canal
	handler EventPosHandler
	table   string
}

type handler struct {
	canal.DummyEventHandler
	*Subscriber
}

type EventHandler interface {
	OnEvent(pos []byte, data []byte) error
}

type EventPosHandler interface {
	EventHandler
	GetPos() ([]byte, error)
	OnRotate(pos []byte) error
}

func NewSubscriber(dbName string, table string) (*Subscriber, error) {
	cfg := canal.NewDefaultConfig()
	cfg.Addr = os.Getenv(dbName + "_MYSQL_ADDR")
	cfg.User = os.Getenv(dbName + "_MYSQL_USER")
	cfg.Password = os.Getenv(dbName + "_MYSQL_PASSWORD")
	cfg.Dump.TableDB = os.Getenv(dbName + "_MYSQL_DATABASE")
	cfg.Dump.Tables = []string{table}
	c, err := canal.NewCanal(cfg)
	if err != nil {
		return nil, err
	}
	s := &Subscriber{}
	s.canal = c
	s.table = table
	return s, nil
}

func (s *Subscriber) Run(ctx context.Context, h EventPosHandler) (err error) {
	s.handler = h
	s.canal.SetEventHandler(&handler{Subscriber: s})
	pos, err := s.handler.GetPos()
	if err != nil {
		return err
	}
	go func() {
		<-ctx.Done()
		s.canal.Close()
	}()
	defer func() {
		// Use Error() because canal wraps the error with errors.WithStack
		if err.Error() == "context canceled" {
			// Ignore context cancel error
			err = nil
		}
	}()
	if pos == nil {
		return s.canal.Run()
	} else {
		var mpos mysql.Position
		mpos.Pos = binary.LittleEndian.Uint32(pos)
		mpos.Name = string(pos[4:])
		return s.canal.RunFrom(mpos)
	}
}

type posFileHandler struct {
	handler EventHandler
	posFile string
	m       sync.Mutex
	pos     []byte
	posSema chan struct{}
}

func (s *Subscriber) RunPosFile(ctx context.Context, handler EventHandler, posFile string) error {
	var wg sync.WaitGroup
	h := &posFileHandler{handler: handler, posFile: posFile}
	h.posSema = make(chan struct{}, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _ = range h.posSema {
			h.m.Lock()
			curPos := h.pos
			h.m.Unlock()
			h.doSavePos(curPos)
		}
	}()
	err := s.Run(ctx, h)
	close(h.posSema)
	wg.Wait()
	return err
}

func (h *posFileHandler) GetPos() ([]byte, error) {
	data, err := ioutil.ReadFile(h.posFile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Warning("Position file not found, may lose some updates")
			return nil, nil
		}
		return nil, err
	}
	return data, nil
}

func (h *posFileHandler) OnEvent(pos []byte, data []byte) error {
	if err := h.handler.OnEvent(pos, data); err != nil {
		return err
	}
	h.setPos(pos)
	return nil
}

func (h *posFileHandler) OnRotate(pos []byte) error {
	h.setPos(pos)
	return nil
}

func (h *posFileHandler) setPos(pos []byte) {
	h.m.Lock()
	h.pos = pos
	select {
	case h.posSema <- struct{}{}:
	default:
	}
	h.m.Unlock()
}

func (h *posFileHandler) doSavePos(pos []byte) {
	defer time.Sleep(time.Second)
	f, err := os.OpenFile(h.posFile+".tmp", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		log.WithError(err).Error("Failed to save pos")
		return
	}
	data := bytes.NewBuffer(pos)
	_, err = io.Copy(f, data)
	if err != nil {
		log.WithError(err).Error("Failed to save pos")
		return
	}
	if err = f.Sync(); err != nil {
		log.WithError(err).Error("Failed to save pos")
		return
	}
	f.Close()
	err = os.Rename(h.posFile+".tmp", h.posFile)
	if err != nil {
		log.WithError(err).Error("Failed to save pos")
		return
	}
	return
}

func encodePos(pos mysql.Position) []byte {
	bpos := make([]byte, 4+len(pos.Name))
	binary.LittleEndian.PutUint32(bpos, pos.Pos)
	copy(bpos[4:], pos.Name)
	return bpos
}

func (h *handler) OnRow(e *canal.RowsEvent) error {
	// insert: 1 e.Rows entry per row
	// update: 2 e.Rows entry per row, older and newer rows respectively
	if e.Action == "insert" && e.Table.Name == h.table {
		for _, row := range e.Rows {
			spos := h.canal.SyncedPosition()
			spos.Pos = e.Header.LogPos
			err := h.handler.OnEvent(encodePos(spos), row[0].([]byte))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// We catch only OnRotate events, not OnPosSynced.
// OnPosSynced events are sometimes emitted at shutdown with an outdated position.
// This causes our position to go backwards.
func (h *handler) OnRotate(ev *replication.RotateEvent) error {
	pos := mysql.Position{}
	pos.Pos = uint32(ev.Position)
	pos.Name = string(ev.NextLogName)
	return h.handler.OnRotate(encodePos(pos))
}
