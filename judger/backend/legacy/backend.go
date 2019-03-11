package legacy

import (
	"context"
	"errors"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/syzoj/syzoj-ng-go/judger/judger"
	"github.com/syzoj/syzoj-ng-go/judger/rpc"
)

var log = logrus.StandardLogger()

type ProblemConf struct {
	Subtasks         []*ProblemSubtask         `yaml:"subtasks"`
	InputFile        string                    `yaml:"inputFile"`
	OutputFile       string                    `yaml:"outputFile"`
	SpecialJudge     ProblemSpecialJudge       `yaml:"specialJudge"`
	ExtraSourceFiles []*ProblemExtraSourceFile `yaml:"extraSourceFiles"`
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

type lojJudgeContext struct {
	ctx      context.Context
	jctx     *judger.JudgeContext
	conf     ProblemConf
	data     *ProblemData
	path     string // problem data path
	tempDir  string
	hasConf  bool
	execPath string
	content  *SubmissionContent

	results  map[string]*SubmissionResultTestcase
	sumScore float64
}

var ErrUnknownLanguage = errors.New("Unknown language")
var ErrInvalidData = errors.New("Invalid test data")

func (lojBackend) JudgeSubmission(ctx context.Context, task *rpc.Task) error {
	submissionContent := new(SubmissionContent)
	if err := ptypes.UnmarshalAny(task.SubmissionContent, submissionContent); err != nil {
		return err
	}
	problemData := new(ProblemData)
	if err := ptypes.UnmarshalAny(task.ProblemData, problemData); err != nil {
		return err
	}

	c := new(lojJudgeContext)
	c.ctx = ctx
	c.jctx = judger.GetJudgeContext(ctx)
	c.data = problemData
	c.path = filepath.Join("data", *task.ProblemId.Id)
	c.content = submissionContent
	return c.judge()
}

func (c *lojJudgeContext) reportResult(result *SubmissionResult) {
	val, err := ptypes.MarshalAny(result)
	if err != nil {
		log.WithError(err).Error("Failed to marshal result")
		return
	}
	if err = c.jctx.ReportResult(val); err != nil {
		log.WithError(err).Error("Failed to report result")
		return
	}
}

func (c *lojJudgeContext) judge() error {
	confData, err := ioutil.ReadFile(filepath.Join(c.path, "problem.conf"))
	if err != nil && !os.IsNotExist(err) {
		return err
	} else if !os.IsNotExist(err) {
		c.hasConf = true
		err = yaml.Unmarshal(confData, &c.conf)
		if err != nil {
			return err
		}
	}

	if err := c.prepareTempDir(); err != nil {
		return err
	}
	if err := c.compile(); err != nil {
		return err
	}
	{
		var err error
		if c.hasConf {
			err = c.judgeWithConf()
		} else {
			err = c.judgeWithoutConf()
		}
		if err != nil {
			return err
		}
	}
	c.reportResult(&SubmissionResult{Score: proto.Float64(c.sumScore)})
	c.cleanup()
	return nil
}

func (c *lojJudgeContext) prepareTempDir() error {
	dir, err := ioutil.TempDir("", "SYZOJ")
	if err != nil {
		return err
	}
	c.tempDir = dir
	return nil
}

func (c *lojJudgeContext) compile() error {
	switch c.content.GetLanguage() {
	case "cpp":
		fpath := filepath.Join(c.tempDir, "code.cpp")
		if err := ioutil.WriteFile(fpath, []byte(c.content.GetCode()), 0644); err != nil {
			return err
		}
		c.execPath = filepath.Join(c.tempDir, "out")
		cmd := exec.Command("g++", fpath, "-O2", "-o", c.execPath)
		_, err := cmd.Output()
		if err != nil {
			// TODO: parse error output
			return err
		}
	default:
		return ErrUnknownLanguage
	}
	return nil
}

func (c *lojJudgeContext) judgeWithConf() error {
	log.WithField("path", c.path).Debug("Judging with conf")
	c.results = make(map[string]*SubmissionResultTestcase)
	for _, subtask := range c.conf.Subtasks {
		for _, n := range subtask.Cases {
			if err := c.judgeTestcase(n); err != nil {
				return err
			}
			if subtask.Type == "min" && c.results[n].GetScore() == 0 {
				break
			}
		}
		var score float64
		switch subtask.Type {
		case "sum":
			for _, n := range subtask.Cases {
				score += c.results[n].GetScore()
			}
			score /= float64(len(subtask.Cases))
		case "min":
			score = 1
			for _, n := range subtask.Cases {
				s := c.results[n].GetScore()
				if s < score {
					score = s
				}
			}
		case "mul":
			for _, n := range subtask.Cases {
				score *= c.results[n].GetScore()
			}
		default:
			return ErrInvalidData
		}
		score *= subtask.Score / 100
		c.sumScore += score
	}
	return nil
}

func (c *lojJudgeContext) judgeWithoutConf() error {
	log.WithField("path", c.path).Debug("Judging without conf")
	c.results = make(map[string]*SubmissionResultTestcase)
	dir, err := ioutil.ReadDir(c.path)
	if err != nil {
		return err
	}
	names := make(map[string]struct {
		a bool
		b bool
	})
	for _, info := range dir {
		name := info.Name()
		if strings.HasSuffix(name, ".in") {
			bname := name[:len(name)-3]
			x := names[bname]
			x.a = true
			names[bname] = x
		} else if strings.HasSuffix(name, ".out") || strings.HasSuffix(name, ".ans") {
			bname := name[:len(name)-4]
			x := names[bname]
			x.b = true
			names[bname] = x
		}
	}
	var score float64
	var cnt int
	for name, x := range names {
		if x.a && x.b {
			if err := c.judgeTestcase(name); err != nil {
				return err
			}
			result := c.results[name]
			cnt++
			score += result.GetScore()
		}
	}
	score /= float64(cnt)
	c.sumScore = score
	log.Debugf("Score is %f\n", score)
	return nil
}

func (c *lojJudgeContext) judgeTestcase(name string) error {
	result := &SubmissionResultTestcase{
		Score: proto.Float64(rand.Float64()),
	}
	c.results[name] = result
	return nil
}

func (c *lojJudgeContext) cleanup() error {
	return os.RemoveAll(c.tempDir)
}

func init() {
	judger.RegisterBackend("legacy", lojBackend{})
}
