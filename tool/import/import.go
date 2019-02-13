package tool_import

import (
	"context"
	"database/sql"
	"encoding/hex"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/sirupsen/logrus"
)

var log = logrus.StandardLogger()

type importer struct {
	mongodb  *mongo.Database
	problems chan *problem
	users    chan *user
	db       *sql.DB
}
type problem struct {
	Id           string
	Title        string
	Description  string
	InputFormat  string
	OutputFormat string
	Example      string
	LimitAndHint string
	Count        int
}

func ImportMySQL(mongodb *mongo.Client, mysql *sql.DB) {
	i := &importer{
		mongodb:  mongodb.Database("syzoj"),
		problems: make(chan *problem),
		users:    make(chan *user),
		db:       mysql,
	}
	i.work()
}

func (i *importer) work() {
	log.Info("Importing problems")
	go i.getProblems()
	i.writeProblems()
	log.Info("Importing users")
	go i.readUsers()
	i.writeUsers()
}

func (i *importer) getProblems() {
	var err error
	var rows *sql.Rows
	if rows, err = i.db.Query("SELECT id, title, user_id, description, input_format, output_format, example, limit_and_hint, time_limit, memory_limit, additional_file_id, is_public, file_io_input_name, file_io_output_name, type FROM problem"); err != nil {
		log.Fatal("Error importing problems from MySQL: ", err.Error())
	}
	for rows.Next() {
		var id int
		p := new(problem)
		var d interface{}
		err = rows.Scan(&id, &p.Title, &d, &p.Description, &p.InputFormat, &p.OutputFormat, &p.Example, &p.LimitAndHint, &d, &d, &d, &d, &d, &d, &d)
		if err != nil {
			log.WithField("id", id).Info("Error reading problem: ", err.Error())
			err = nil
		}
		p.Id = strconv.Itoa(id)
		i.problems <- p
	}
	close(i.problems)
}

// LEGACY
var inlineMathRegexp = regexp.MustCompile("\\$(.+?)\\$")
var mathRegexp = regexp.MustCompile("(?m)\\$\\$(.+?)\\$\\$")

func convertMath(s string) string {
	s2 := inlineMathRegexp.ReplaceAll([]byte(s), []byte("<math inline>$1</math>"))
	return string(mathRegexp.ReplaceAll(s2, []byte("<math>$1</math>")))
}

func (i *importer) writeProblems() {
	var err error
	for p := range i.problems {
		problemId := primitive.NewObjectID()
		var content []string
		if p.Description != "" {
			content = append(content, "# 题目描述\n", p.Description)
		}
		if p.InputFormat != "" {
			content = append(content, "# 输入格式\n", p.InputFormat)
		}
		if p.OutputFormat != "" {
			content = append(content, "# 输出格式\n", p.OutputFormat)
		}
		if p.Example != "" {
			content = append(content, "# 样例\n", p.Example)
		}
		if p.LimitAndHint != "" {
			content = append(content, "# 数据范围与提示\n", p.LimitAndHint)
		}
		contents := strings.Join(content, "\n\n")
		if _, err = i.mongodb.Collection("problem").InsertOne(context.Background(),
			bson.D{{"_id", problemId}, {"title", p.Title}, {"statement", contents}, {"short_name", p.Id}},
		); err != nil {
			log.WithField("id", p.Id).Info("Error inserting problem: ", err.Error())
			err = nil
		}
	}
}

type user struct {
	UserName     string
	Password     string
	Email        string
	RegisterTime sql.NullInt64
}

func (i *importer) readUsers() {
	var err error
	var rows *sql.Rows
	if rows, err = i.db.Query("SELECT username, password, email, register_time FROM user"); err != nil {
		log.Fatal("Error importing users from MySQL: ", err.Error())
	}
	for rows.Next() {
		u := new(user)
		err = rows.Scan(&u.UserName, &u.Password, &u.Email, &u.RegisterTime)
		if err != nil {
			log.Error("Error reading user: ", err)
			err = nil
		}
		i.users <- u
	}
	close(i.users)
}

func (i *importer) writeUsers() {
	var err error
	for user := range i.users {
		var passmd5 []byte
		passmd5, err = hex.DecodeString(user.Password)
		if err != nil {
			log.WithField("username", user.UserName).Error("Error parsing password")
			err = nil
			continue
		}
		doc := bson.D{
			{"_id", primitive.NewObjectID()},
			{"username", user.UserName},
			{"email", user.Email},
			{"auth", bson.D{
				{"method", int64(2)},
				{"password", passmd5},
			}},
		}
		if user.RegisterTime.Valid {
			doc = append(doc, bson.E{"register_time", time.Unix(user.RegisterTime.Int64, 0)})
		}
		if _, err = i.mongodb.Collection("user").InsertOne(context.Background(), doc); err != nil {
			log.Error("Error inserting user: ", err)
			err = nil
		}
	}
}
