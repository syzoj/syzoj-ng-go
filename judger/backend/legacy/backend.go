package legacy

import (
	"context"
	"io/ioutil"
	"path/filepath"

	"github.com/golang/protobuf/ptypes"
	"gopkg.in/yaml.v2"

	"github.com/syzoj/syzoj-ng-go/judger/judger"
	"github.com/syzoj/syzoj-ng-go/judger/rpc"
)

type ProblemConf struct {
	Subtasks         []ProblemSubtask         `yaml:"subtasks"`
	InputFile        string                   `yaml:"inputFile"`
	OutputFile       string                   `yaml:"outputFile"`
	SpecialJudge     ProblemSpecialJudge      `yaml:"specialJudge"`
	ExtraSourceFiles []ProblemExtraSourceFile `yaml:"extraSourceFiles"`
}

type ProblemSubtask struct {
	Score float64  `yaml:"score"`
	Type  string   `yaml:"type"` // sum, min, mul
	Cases []string `yaml:"cases"`
}

type ProblemSpecialJudge struct {
	Language string `yaml:"language"`
	Code     string `yaml:"code"`
}

type ProblemExtraSourceFile struct {
	Language string                        `yaml:"language"`
	Files    []ProblemExtraSourceFileEntry `yaml:"files"`
}

type ProblemExtraSourceFileEntry struct {
	Name string `yaml:"name"`
	Dest string `yaml:"dest"`
}

type lojBackend struct{}

func (lojBackend) JudgeSubmission(ctx context.Context, task *rpc.Task) error {
	jctx := judger.GetJudgeContext(ctx)
	problemData := new(ProblemData)
	if err := ptypes.UnmarshalAny(task.ProblemData, problemData); err != nil {
		return err
	}

	confData, err := ioutil.ReadFile(filepath.Join("data", *task.ProblemId.Id, "problem.conf"))
	if err != nil {
		return err
	}
	var conf ProblemConf
	err = yaml.Unmarshal(confData, &conf)
	if err != nil {
		return err
	}
}

func init() {
	judger.RegisterBackend("legacy", lojBackend{})
}
