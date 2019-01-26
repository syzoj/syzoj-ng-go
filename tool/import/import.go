package tool_import

import (
    "context"
    "database/sql"
    "strconv"
    "strings"
    "regexp"

    "github.com/mongodb/mongo-go-driver/mongo"
    "github.com/mongodb/mongo-go-driver/bson"
    "github.com/mongodb/mongo-go-driver/bson/primitive"
    "github.com/sirupsen/logrus"
)
var log = logrus.StandardLogger()

type importer struct {
    mongodb *mongo.Database
    problems chan *problem
    db *sql.DB
}
type problem struct {
    Id string
    Title string
    Description string
    InputFormat string
    OutputFormat string
    Example string
    LimitAndHint string
    Count int
}

func ImportMySQL(mongodb *mongo.Client, mysql *sql.DB) {
    i := &importer{
        mongodb: mongodb.Database("syzoj"),
        problems: make(chan *problem),
        db: mysql,
    }
    i.work()
}

func (i *importer) work() {
    log.Info("Importing problems")
    go i.getProblems()
    i.writeProblems()
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
            bson.D{{"_id", problemId}, {"title", p.Title}, {"statement", contents}},
        ); err != nil {
            log.WithField("id", p.Id).Info("Error insering problem: ", err.Error())
            err = nil
        }
    }
}
