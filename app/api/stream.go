package api

import (
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type StreamSender struct {
	body      io.Reader
	timer     *time.Timer
	closeChan chan struct{}
	closeOnce sync.Once
}

type StreamReceiver struct {
	srv   *ApiServer
	token string

	o sync.Once
	// notify is closed either when StreamReceiver is closed or when a sender is ready.
	notifyChan chan struct{}
	sender     *StreamSender
}

func (srv *ApiServer) HandleStream(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]
	srv.streamLock.Lock()
	if _, ok := srv.streamSender[token]; ok {
		srv.streamLock.Unlock()
		http.Error(w, "Conflict", 409)
		return
	}
	sender := &StreamSender{
		body:      r.Body,
		closeChan: make(chan struct{}),
		timer:     time.NewTimer(time.Second * 120),
	}
	srv.streamSender[token] = sender
	if _, ok := srv.streamReceiver[token]; ok {
		pairStream(srv.streamSender[token], srv.streamReceiver[token])
		delete(srv.streamSender, token)
		delete(srv.streamReceiver, token)
	}
	srv.streamLock.Unlock()
	select {
	case <-r.Context().Done():
	case <-sender.closeChan:
	case <-sender.timer.C:
	}
	srv.streamLock.Lock()
	if s2, _ := srv.streamSender[token]; s2 == sender {
		delete(srv.streamSender, token)
	}
	srv.streamLock.Unlock()
}

func (srv *ApiServer) GetStream(token string) *StreamReceiver {
	srv.streamLock.Lock()
	if _, ok := srv.streamReceiver[token]; ok {
		srv.streamLock.Unlock()
		return nil
	}
	receiver := &StreamReceiver{
		srv:        srv,
		token:      token,
		notifyChan: make(chan struct{}),
	}
	srv.streamReceiver[token] = receiver
	if _, ok := srv.streamSender[token]; ok {
		pairStream(srv.streamSender[token], srv.streamReceiver[token])
		delete(srv.streamSender, token)
		delete(srv.streamReceiver, token)
	}
	srv.streamLock.Unlock()
	return receiver
}

func pairStream(sender *StreamSender, receiver *StreamReceiver) {
	go func() {
		receiver.o.Do(func() {
			if sender.timer.Stop() {
				receiver.sender = sender
			}
			close(receiver.notifyChan)
		})
	}()
}

func (s *StreamReceiver) Close() error {
	s.srv.streamLock.Lock()
	s2, _ := s.srv.streamReceiver[s.token]
	if s2 == s {
		delete(s.srv.streamReceiver, s.token)
	}
	s.srv.streamLock.Unlock()

	s.o.Do(func() {
		close(s.notifyChan)
	})
	if s.sender != nil {
		s.sender.closeOnce.Do(func() {
			close(s.sender.closeChan)
		})
	}
	return nil
}

// Read implements the io.Reader interface.
func (s *StreamReceiver) Read(p []byte) (int, error) {
	<-s.notifyChan
	if s.sender == nil {
		return 0, io.ErrUnexpectedEOF
	}
	return s.sender.body.Read(p)
}

// WriteTo implements the io.WriterTo interface.
func (s *StreamReceiver) WriteTo(w io.Writer) (int64, error) {
	<-s.notifyChan
	if s.sender == nil {
		return 0, io.ErrUnexpectedEOF
	}
	return io.Copy(w, s.sender.body)
}
