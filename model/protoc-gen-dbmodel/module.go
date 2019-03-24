package main

import (
    "errors"
    "fmt"
    "sort"
    "strings"
    "text/template"
    "os"

    pgs "github.com/lyft/protoc-gen-star"
    pgsgo "github.com/lyft/protoc-gen-star/lang/go"
    "github.com/golang/protobuf/protoc-gen-go/descriptor"

    "github.com/syzoj/syzoj-ng-go/model/protoc-gen-dbmodel/dbmodel"
)

var tplOrm = template.Must(template.New("dbmodel_orm").Parse(`
package database

import (
    "context"
    "database/sql"

    "github.com/syzoj/syzoj-ng-go/model"
)

{{range .Tables}}
func (t *DatabaseTxn) Get{{.CapName}}(ctx context.Context, ref model.{{.CapName}}Ref) (*model.{{.CapName}}, error) {
    v := new(model.{{.CapName}})
    err := t.tx.QueryRowContext(ctx, "SELECT {{.SelList}} FROM {{.Name}} WHERE id=?", ref).Scan(&v.Id, {{.ScanList}})
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }
    return v, nil
}

func (t *DatabaseTxn) Update{{.CapName}}(ctx context.Context, ref model.{{.CapName}}Ref, v *model.{{.CapName}}) error {
    if v.Id == nil || v.GetId() != ref {
        panic("ref and v does not match")
    }
    _, err := t.tx.ExecContext(ctx, "UPDATE {{.Name}} SET {{.UpdateList}} WHERE id=?", {{.ArgList}}, v.Id)
    return err
}

func (t *DatabaseTxn) Insert{{.CapName}}(ctx context.Context, v *model.{{.CapName}}) error {
    if v.Id == nil {
        ref := model.New{{.CapName}}Ref()
        v.Id = &ref
    }
    _, err := t.tx.ExecContext(ctx, "INSERT INTO {{.Name}} (id, {{.InsList}}) VALUES ({{.InsValue}})", v.Id, {{.ArgList}})
    return err
}

func (t *DatabaseTxn) Delete{{.CapName}}(ctx context.Context, ref model.{{.CapName}}Ref) error {
    _, err := t.tx.ExecContext(ctx, "DELETE FROM {{.Name}} WHERE id=?", ref)
    return err
}
{{end}}
`))
var tplModel = template.Must(template.New("dbmodel_model").Parse(`
package model

import (
    "crypto/rand"
    "database/sql/driver"
    "encoding/base64"
    "errors"

    "github.com/golang/protobuf/proto"
)
var ErrInvalidType = errors.New("Can only scan []byte into protobuf message")

func newId() string {
    var b [12]byte
    if _, err := rand.Read(b[:]); err != nil {
        panic(err)
    }
    return base64.URLEncoding.EncodeToString(b[:])
}

{{range .Tables}}
type {{.CapName}}Ref string

func New{{.CapName}}Ref() {{.CapName}}Ref {
    return {{.CapName}}Ref(newId())
}
{{end}}

{{range .Messages}}
func (m *{{.}}) Value() (driver.Value, error) {
    return proto.Marshal(m)
}

func (m *{{.}}) Scan(v interface{}) error {
    if v == nil {
        return nil
    }
    if b, ok := v.([]byte); ok {
        return proto.Unmarshal(b, m)
    }
    return ErrInvalidType
}
{{end}}
`))
var tplSql = template.Must(template.New("dbmodel_sql").Parse(`{{range .Tables}}CREATE TABLE {{.Name}} (
{{.SqlFields}}
);

{{end}}
`))

type module struct {
    *pgs.ModuleBase
    ctx pgsgo.Context
}

func newModule() pgs.Module {
    return &module{ModuleBase: &pgs.ModuleBase{}}
}

func (m *module) InitContext(c pgs.BuildContext) {
    m.ModuleBase.InitContext(c)
    m.ctx = pgsgo.InitContext(c.Parameters())
}

func (m *module) Name() string {
    return "dbmodel"
}

func (m *module) Execute(targets map[string]pgs.File, packages map[string]pgs.Package) []pgs.Artifact {
    for _, f := range targets {
        v := makeVisitor(m)
        if err := pgs.Walk(v, f); err != nil {
            panic(err)
        }
        data := v.getData()
        m.AddCustomTemplateFile("dbmodel_orm.go", tplOrm, data, 0644)
        m.AddCustomTemplateFile("dbmodel_model.go", tplModel, data, 0644)
        m.AddCustomTemplateFile("dbmodel_sql.sql", tplSql, data, 0644)
    }
    return m.Artifacts()
}

type visitor struct {
    pgs.Visitor
    pgs.DebuggerCommon
    d tplData
}
type tplData struct {
    Tables   []tplTable
    Messages []string
}
type tplTable struct {
    Name       string
    CapName    string
    SelList    string
    UpdateList string
    ScanList   string
    ArgList    string
    InsList    string
    InsValue   string
    SqlFields string
}

func makeVisitor(d pgs.DebuggerCommon) *visitor {
    return &visitor{
        Visitor:        pgs.NilVisitor(),
        DebuggerCommon: d,
    }
}

func (v *visitor) VisitPackage(pgs.Package) (pgs.Visitor, error) { return v, nil }
func (v *visitor) VisitFile(pgs.File) (pgs.Visitor, error)       { return v, nil }
func (v *visitor) VisitMessage(m pgs.Message) (pgs.Visitor, error) {
    var t tplTable
    t.Name = m.Name().LowerSnakeCase().String()
    t.CapName = m.Name().String()
    var selList []string
    var updateList []string
    var argList []string
    var scanList []string
    var insList []string
    var insValue []string
    var sqlFields []string
    for i, f := range m.Fields() {
        insValue = append(insValue, "?")
        if i == 0 && f.Name().String() != "id" {
            return nil, errors.New("The first field of a database model must be named \"id\"")
        }
        if f.Type().IsMap() || f.Type().IsRepeated() {
            return nil, errors.New("Map or repeated fields in a database model is not allowed")
        }
        if i != 0 {
            selList = append(selList, f.Name().String())
            updateList = append(updateList, f.Name().String()+"=?")
            argList = append(argList, "v."+f.Name().UpperCamelCase().String())
            scanList = append(scanList, "&v."+f.Name().UpperCamelCase().String())
            insList = append(insList, f.Name().String())
        }
        if m := f.Type().Embed(); m != nil {
            v.d.Messages = append(v.d.Messages, m.Name().String())
        }
        var sql string
        if ok, _ := f.Extension(dbmodel.E_Sql, &sql); ok {
        } else {
            if i == 0 {
                sql = "id VARCHAR(16) PRIMARY KEY"
            } else {
                t := f.Type().ProtoType()
                if t.IsInt() {
                    sql = fmt.Sprintf("%s BIGINT", f.Name().String())
                } else {
                    switch t.Proto() {
                    case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
                        sql = fmt.Sprintf("%s DOUBLE", f.Name().String())
                    case descriptor.FieldDescriptorProto_TYPE_FLOAT:
                        sql = fmt.Sprintf("%s FLOAT", f.Name().String())
                    case descriptor.FieldDescriptorProto_TYPE_STRING:
                        sql = fmt.Sprintf("%s VARCHAR(255)", f.Name().String())
                    case descriptor.FieldDescriptorProto_TYPE_BYTES, descriptor.FieldDescriptorProto_TYPE_MESSAGE:
                        sql = fmt.Sprintf("%s BLOB", f.Name().String())
                    default:
                        return nil, errors.New(fmt.Sprintf("Cannot generate SQL statement for %s.%s", m.Name().String(), f.Name().String()))
                    }
                }
            }
        }
        sqlFields = append(sqlFields, "  " + sql)
        fmt.Sprintln(os.Stderr, sql)
    }
    t.SelList = strings.Join(selList, ", ")
    t.UpdateList = strings.Join(updateList, ", ")
    t.ArgList = strings.Join(argList, ", ")
    t.ScanList = strings.Join(scanList, ", ")
    t.InsList = strings.Join(insList, ", ")
    t.InsValue = strings.Join(insValue, ", ")
    t.SqlFields = strings.Join(sqlFields, ",\n")
    v.d.Tables = append(v.d.Tables, t)
    return v, nil
}
func (v *visitor) getData() interface{} {
    sort.Strings(v.d.Messages)
    var i, j int
    for i = 0; i < len(v.d.Messages); i++ {
        if i == len(v.d.Messages)-1 || v.d.Messages[i] != v.d.Messages[i+1] {
            v.d.Messages[j] = v.d.Messages[i]
            j++
        }
    }
    v.d.Messages = v.d.Messages[:j]
    return v.d
}
