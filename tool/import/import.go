package tool_import

import (
	"database/sql"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var log = logrus.StandardLogger()

type importer struct {
	mongodb     *mongo.Database
	db          *sql.DB
	oldDataPath string
	newDataPath string

	problemId map[int64]primitive.ObjectID
	userId    map[int64]primitive.ObjectID
}

func ImportMySQL(mongodb *mongo.Client, mysql *sql.DB, oldDataPath string, newDataPath string) {
	i := &importer{
		mongodb:     mongodb.Database("syzoj"),
		db:          mysql,
		oldDataPath: oldDataPath,
		newDataPath: newDataPath,
	}
	i.work()
}

func (i *importer) work() {
	log.Info("Importing problems")
	chanProblems := make(chan *problem)
	go i.getProblems(chanProblems)
	i.writeProblems(chanProblems)

	log.Info("Importing users")
	chanUsers := make(chan *user)
	go i.readUsers(chanUsers)
	i.writeUsers(chanUsers)

	log.Info("Importing submissions")
	chanSubmissions := make(chan *submission)
	go i.readSubmissions(chanSubmissions)
	i.writeSubmissions(chanSubmissions)
}
