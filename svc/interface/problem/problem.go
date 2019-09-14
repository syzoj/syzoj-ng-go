// The problem service.
package problem

import (
	"fmt"
	"bytes"
	"context"
	"database/sql"
	"errors"
	"strings"
	"html"

	"github.com/jmoiron/sqlx"
	"github.com/beevik/etree"
	"github.com/sirupsen/logrus"
	"github.com/minio/minio-go"
	lredis "github.com/syzoj/syzoj-ng-go/lib/redis"
)

var log = logrus.StandardLogger()
var ErrNotFound = errors.New("Not found")
var ErrXMLTooBig = errors.New("XML too large or too complex")
var ErrXMLEmpty = errors.New("XML empty")

// ProblemService is a service that manages problems.
type ProblemService struct {
	Db *sqlx.DB
	RedisCache *lredis.PoolWrapper
	Minio *minio.Client
	RedisPrefix string
	TestdataBucket string
	XMLSizeLimit int
	XMLNodeLimit int
}

// Creates a problem service with default settings.
func DefaultProblemService(db *sqlx.DB, redis *lredis.PoolWrapper, minio *minio.Client, testdataBucket string) *ProblemService {
	return &ProblemService{
		Db: db,
		RedisCache: redis,
		Minio: minio,
		RedisPrefix: "problem-cache:",
		TestdataBucket: testdataBucket,
		XMLSizeLimit: 1024 * 1024,
		XMLNodeLimit: 10000,
	}
}

// Warner is an interface that receives warnings generated during parsing.
type Warner interface {
	Warningf(format string, args ...interface{})
}

// WarningList is a simple implementation for Warner.
type WarningList []string

func (w *WarningList) Warningf(format string, args ...interface{}) {
	*w = append(*w, fmt.Sprintf(format, args...))
}

func (w *WarningList) Error() string {
	return strings.Join(*w, "\n")
}

func (w *WarningList) renderHTML(buf *bytes.Buffer) {
	buf.WriteString("<p>")
	buf.WriteString(html.EscapeString(w.Error()))
	buf.WriteString("</p>")
}

// Inserts a problem into DB.
func (s *ProblemService) CreateProblem(ctx context.Context, problemId string, body []byte) error {
	if len(body) > s.XMLSizeLimit {
		return ErrXMLTooBig
	}
	// Verify XML and check node count.
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(body); err != nil {
		return err
	}
	if doc.Root() == nil {
		return ErrXMLEmpty
	}
	cnt := countNodes(doc.Root())
	if cnt > s.XMLNodeLimit {
		return ErrXMLTooBig
	}
	body, err := doc.WriteToBytes()
	if err != nil {
		return err
	}
	const SQLCreateProblem = "INSERT INTO `problems` (`uid`, `body`) VALUES (?, ?)"
	_, err = s.Db.ExecContext(ctx, SQLCreateProblem, problemId, body)
	return err
}

// Updates a problem in DB.
func (s *ProblemService) UpdateProblem(ctx context.Context, problemId string, body []byte) error {
	if len(body) > s.XMLSizeLimit {
		return ErrXMLTooBig
	}
	// Verify XML and check node count.
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(body); err != nil {
		return err
	}
	if doc.Root() == nil {
		return ErrXMLEmpty
	}
	cnt := countNodes(doc.Root())
	if cnt > s.XMLNodeLimit {
		return ErrXMLTooBig
	}
	body, err := doc.WriteToBytes()
	if err != nil {
		return err
	}
	const SQLCreateProblem = "UPDATE `problems` SET `body`=? WHERE `uid`=?"
	_, err = s.Db.ExecContext(ctx, SQLCreateProblem, body, problemId)
	return err
}

// Gets and parses the problem's XML body.
func (s *ProblemService) getProblemXML(ctx context.Context, problemId string) (*etree.Document, error) {
	var body []byte
	const SQLGetProblem = "SELECT `body` FROM `problems` WHERE `uid`=?"
	if err := s.Db.QueryRowContext(ctx, SQLGetProblem, problemId).Scan(&body); err != nil {
		if err == sql.ErrNoRows {
			err = ErrNotFound
		}
		return nil, err
	}
	doc := etree.NewDocument()
	err := doc.ReadFromBytes(body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// Check for cache first before parsing.
func (s *ProblemService) doCache(ctx context.Context, problemId string, name string, f func() ([]byte, error)) ([]byte, error) {
	return lredis.WithCache(ctx, s.RedisCache, s.RedisPrefix + problemId + ":" + name, f)
}

// Gets problem statement for problem in HTML.
// Note that a partially rendered HTML might also be returned even if error is not nil.
func (s *ProblemService) GetStatementHTML(ctx context.Context, problemId string) ([]byte, error) {
	return s.doCache(ctx, problemId, "stmt", func() ([]byte, error) {
		doc, err := s.getProblemXML(ctx, problemId)
		if err != nil {
			return nil, err
		}
		root := doc.Root()
		buf := &bytes.Buffer{}
		w := new(WarningList)
		if root != nil && root.Tag == "Problem" {
			stmt := root.SelectElement("Statement")
			if stmt != nil {
				renderStmt(buf, w, stmt)
			} else {
				w.Warningf("Statement node doesn't exist")
			}
		} else {
			w.Warningf("Root node doesn't exist or is not <Problem>")
		}
		w.renderHTML(buf)
		return buf.Bytes(), nil
	})
}

func countNodes(tok etree.Token) int {
	switch obj := tok.(type) {
	case *etree.Element:
		cnt := 1
		for _, ch := range obj.Child {
			cnt += countNodes(ch)
		}
		return cnt
	default:
		return 1
	}
}
