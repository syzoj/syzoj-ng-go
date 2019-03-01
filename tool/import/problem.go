package tool_import

import (
    "database/sql"
    "strconv"
    "regexp"
    "strings"
    "path"
    "context"

    "github.com/golang/protobuf/proto"

    "github.com/syzoj/syzoj-ng-go/app/model"
)

type problem struct {
	Id           string
	Title        string
	Description  string
	InputFormat  string
	OutputFormat string
	Example      string
	LimitAndHint string
    Type sql.NullString
    MemoryLimit sql.NullInt64
    TimeLimit sql.NullInt64
    FileIoInput sql.NullString
    FileIoOutput sql.NullString
	Count        int
}

func (i *importer) getProblems(problems chan *problem) {
	var err error
	var rows *sql.Rows
	if rows, err = i.db.Query("SELECT id, title, user_id, description, input_format, output_format, example, limit_and_hint, time_limit, memory_limit, additional_file_id, is_public, file_io_input_name, file_io_output_name, type FROM problem"); err != nil {
		log.Fatal("Error importing problems from MySQL: ", err.Error())
	}
	for rows.Next() {
		var id int
		p := new(problem)
		var d interface{}
		err = rows.Scan(&id, &p.Title, &d, &p.Description, &p.InputFormat, &p.OutputFormat, &p.Example, &p.LimitAndHint, &p.MemoryLimit, &p.TimeLimit, &d, &d, &p.FileIoInput, &p.FileIoOutput, &p.Type)
		if err != nil {
			log.WithField("id", id).Info("Error reading problem: ", err.Error())
			err = nil
		}
		p.Id = strconv.Itoa(id)
		problems <- p
	}
	close(problems)
}

// LEGACY
var inlineMathRegexp = regexp.MustCompile("\\$(.+?)\\$")
var mathRegexp = regexp.MustCompile("(?m)\\$\\$(.+?)\\$\\$")

func convertMath(s string) string {
	s2 := inlineMathRegexp.ReplaceAll([]byte(s), []byte("<math inline>$1</math>"))
	return string(mathRegexp.ReplaceAll(s2, []byte("<math>$1</math>")))
}

func (i *importer) writeProblems(problems chan *problem) {
	var err error
	for p := range problems {
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
		problemModel := new(model.Problem)
		problemModel.Id = model.NewObjectIDProto()
		problemModel.Title = proto.String(p.Title)
		problemModel.Statement = proto.String(contents)
		problemModel.ShortName = proto.String(p.Id)
		if _, err = i.mongodb.Collection("problem").InsertOne(context.Background(), problemModel); err != nil {
			log.WithField("id", p.Id).Info("Error inserting problem: ", err.Error())
			err = nil
		}

        name := problemModel.Id.GetId()
        if i.oldDataPath != "" {
            conv_problem(p, path.Join(i.oldDataPath, p.Id), path.Join(i.newDataPath, name))
        }
	}
}
