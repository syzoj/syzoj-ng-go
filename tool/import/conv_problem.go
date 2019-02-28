package tool_import

import (
    "os"
    "path"
    "io/ioutil"
    "strings"

    "gopkg.in/yaml.v2"
    "github.com/golang/protobuf/proto"

    "github.com/syzoj/syzoj-ng-go/app/model"
)

type ProblemYaml struct {
    Subtasks []ProblemYamlSubtask `yaml:"subtasks"`
    InputFile string `yaml:"inputFile"`
    OutputFile string `yaml:"outputFile"`
    SpecialJudge ProblemSpecialJudge `yaml:"specialJudge"`
    ExtraSourceFiles []ProblemExtraSourceFileInfo `yaml:"extraSourceFiles"`
}

type ProblemYamlSubtask struct {
    Score float64 `yaml:"score"`
    Type string `yaml:"type"`
    Cases []string `yaml:"cases"`
}

type ProblemSpecialJudge struct {
    Language string `yaml:"language"`
    FileName string `yaml:"fileName"`
}

type ProblemExtraSourceFileInfo struct {
    Language string `yaml:"language"`
    Files []ProblemExtraSourceFile `yaml:"files"`
}

type ProblemExtraSourceFile struct {
    Name string `yaml:"name"`
    Dest string `yaml:"dest"`
}

func conv_problem(orgpath string, newpath string) {
    var err error
    var data []byte
    if data, err = ioutil.ReadFile(path.Join(orgpath, "data.yml")); err != nil {
        log.WithField("path", orgpath).WithError(err).Warning("Failed to read data.yml")
        return
    }
    if err = os.MkdirAll(newpath, 0755); err != nil {
        log.WithField("path", newpath).WithError(err).Error("Failed to make directory")
        return
    }
    var yml ProblemYaml
    err = yaml.Unmarshal(data, &yml)
    if err != nil {
        log.WithField("path", orgpath).WithError(err).Error("Failed to parse data.yml")
        return
    }

    var newyml model.ProblemConf
    newyml.Type = proto.String("traditional")
    newyml.CasesGlobal = new(model.ProblemCase)
    if yml.InputFile != "" {
        file := strings.Replace(yml.InputFile, "#", "{name}", 100)
        newyml.CasesGlobal.InputData = proto.String(file)
    }
    if yml.OutputFile != "" {
        file := strings.Replace(yml.InputFile, "#", "{name}", 100)
        newyml.CasesGlobal.AnswerData = proto.String(file)
    }

    testcases := make(map[string]*model.ProblemCase)
    for _, org_subtask := range yml.Subtasks {
        new_subtask := new(model.ProblemSubtask)
        new_subtask.Score = proto.Float64(org_subtask.Score)
        for _, c := range org_subtask.Cases {
            new_subtask.Cases = append(new_subtask.Cases, c)
            if _, found := testcases[c]; !found {
                new_case := &model.ProblemCase{Name: proto.String(c)}
                testcases[c] = new_case
            }
        }
        newyml.Subtasks = append(newyml.Subtasks, new_subtask)
    }
    for _, new_case := range testcases {
        newyml.Cases = append(newyml.Cases, new_case)
    }
    if yml.SpecialJudge.Language != "" {
        log.WithField("path", orgpath).Warning("Special judge not yet supported")
    }

    data, err = yaml.Marshal(&newyml)
    if err != nil {
        log.WithField("path", orgpath).WithError(err).Error("Failed to convert data.yml")
        return
    }

    if err = ioutil.WriteFile(path.Join(newpath, "problem.yml"), data, 0644); err != nil {
        log.WithField("path", newpath).WithField("oldpath", orgpath).WithError(err).Error("Failed to write problem.yml")
        return
    }
}
