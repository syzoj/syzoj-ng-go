package main

import (
	"strings"
    "text/template"

	pgs "github.com/lyft/protoc-gen-star"
	pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

var tpl = template.Must(template.New("dbmodel").Parse(`
package database

import (
	"context"
	"database/sql"

	"github.com/syzoj/syzoj-ng-go/model"
)

{{range .}}
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
        m.AddCustomTemplateFile("dbmodel.go", tpl, v.getData(), 0644)
	}
	return m.Artifacts()
}

type visitor struct {
	pgs.Visitor
	pgs.DebuggerCommon
    t []tplTable
}
type tplTable struct {
    Name string
    CapName string
    SelList string
    UpdateList string
    ScanList string
    ArgList string
    InsList string
    InsValue string
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
    for i, f := range m.Fields() {
        insValue = append(insValue, "?")
        if i != 0 {
            selList = append(selList, f.Name().String())
            updateList = append(updateList, f.Name().String() + "=?")
            argList = append(argList, "v." + f.Name().UpperCamelCase().String())
            scanList = append(scanList, "&v." + f.Name().UpperCamelCase().String())
            insList = append(insList, f.Name().String())
        }
    }
    t.SelList = strings.Join(selList, ", ")
    t.UpdateList = strings.Join(updateList, ", ")
    t.ArgList = strings.Join(argList, ", ")
    t.ScanList = strings.Join(scanList, ", ")
    t.InsList = strings.Join(insList, ", ")
    t.InsValue = strings.Join(insValue, ", ")
    v.t = append(v.t, t)
    return v, nil
}
func (v *visitor) getData() interface{} {
    return v.t
}
