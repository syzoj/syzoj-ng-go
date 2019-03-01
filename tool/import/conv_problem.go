package tool_import

import (
    "os"
    "path"
    "path/filepath"
    "strconv"
    "io/ioutil"
    "strings"
    "io"
    "sync"

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

func conv_problem(p *problem, orgpath string, newpath string) {
    var err error
    if !p.Type.Valid || p.Type.String != "traditional" {
        log.WithField("path", orgpath).Warning("Non-traditional problems not supported")
        return
    }
    if err = os.MkdirAll(newpath, 0755); err != nil {
        log.WithField("path", newpath).WithError(err).Error("Failed to make directory")
        return
    }
    if err = os.MkdirAll(path.Join(newpath, "data"), 0755); err != nil {
        log.WithField("path", path.Join(newpath, "data")).WithError(err).Error("Failed to make directory")
        return
    }
    log.WithField("id", p.Id).Info("Converting problem data")

    // Copy all files
    var wg sync.WaitGroup
    // ignore errors
    filepath.Walk(orgpath, func(p string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() {
            return nil
        }
        name := info.Name()
        orgfile, e := os.Open(p)
        if e != nil {
            log.WithField("path", p).Warning("Failed to open file")
            return nil
        }
        // Flatten the directory structure
        newf := path.Join(newpath, "data", name)
        newfile, e := os.Create(newf)
        if e != nil {
            log.WithField("path", newf).Warning("Failed to create file")
            orgfile.Close()
            return nil
        }
        wg.Add(1)
        go func() {
            defer wg.Done()
            defer orgfile.Close()
            defer newfile.Close()
            _, e = io.Copy(newfile, orgfile)
            if e != nil {
                log.WithField("from", p).WithField("to", newf).Warning("Failed to copy file")
            }
        }()
        return nil
    })
    wg.Wait()

    var newyml model.ProblemConf
    newyml.Type = proto.String("traditional")
    newyml.CasesGlobal = new(model.ProblemCase)
    if p.MemoryLimit.Valid {
        newyml.CasesGlobal.MemoryLimit = proto.String(strconv.FormatInt(p.MemoryLimit.Int64, 10) + "MB")
    }
    if p.TimeLimit.Valid {
        newyml.CasesGlobal.TimeLimit = proto.String(strconv.FormatInt(p.TimeLimit.Int64, 10) + "ms")
    }
    if p.FileIoInput.Valid && p.FileIoInput.String != "" {
        newyml.CasesGlobal.InputFile = proto.String(p.FileIoInput.String)
    }
    if p.FileIoOutput.Valid && p.FileIoOutput.String != "" {
        newyml.CasesGlobal.OutputFile = proto.String(p.FileIoOutput.String)
    }

    if data, err := ioutil.ReadFile(path.Join(orgpath, "data.yml")); err == nil {
        var yml ProblemYaml
        err = yaml.Unmarshal(data, &yml)
        if err != nil {
            log.WithField("path", orgpath).WithError(err).Error("Failed to parse data.yml")
            goto parseDone
        }
        if yml.InputFile != "" {
            file := "data/" + strings.Replace(yml.InputFile, "#", "{name}", 100)
            newyml.CasesGlobal.InputData = proto.String(file)
        }
        if yml.OutputFile != "" {
            file := "data/" + strings.Replace(yml.OutputFile, "#", "{name}", 100)
            newyml.CasesGlobal.AnswerData = proto.String(file)
        }

        testcases := make(map[string]*model.ProblemCase)
        for _, org_subtask := range yml.Subtasks {
            switch org_subtask.Type {
            case "min":
                new_subtask := new(model.ProblemSubtask)
                new_subtask.Score = proto.Float64(org_subtask.Score)
                for _, c := range org_subtask.Cases {
                    new_subtask.Testcases = append(new_subtask.Testcases, c)
                    if _, found := testcases[c]; !found {
                        new_case := &model.ProblemCase{Name: proto.String(c)}
                        testcases[c] = new_case
                    }
                }
                newyml.Subtasks = append(newyml.Subtasks, new_subtask)
            case "sum":
                single_score := org_subtask.Score / float64(len(org_subtask.Cases))
                for _, c := range org_subtask.Cases {
                    new_subtask := new(model.ProblemSubtask)
                    new_subtask.Score = proto.Float64(single_score)
                    new_subtask.Testcases = []string{c}
                    if _, found := testcases[c]; !found {
                        new_case := &model.ProblemCase{Name: proto.String(c)}
                        testcases[c] = new_case
                    }
                }
             case "mul":
                 log.WithField("problem_id", p.Id).Warning("Subtask with type mul is not supported")
                 return
             }
        }
        for _, new_case := range testcases {
            newyml.Cases = append(newyml.Cases, new_case)
        }
        if yml.SpecialJudge.Language != "" {
            log.WithField("path", orgpath).Warning("Special judge not yet supported")
        }
    }
    parseDone:

    data, err := yaml.Marshal(&newyml)
    if err != nil {
        log.WithField("path", orgpath).WithError(err).Error("Failed to convert data.yml")
        return
    }

    if err = ioutil.WriteFile(path.Join(newpath, "problem.yml"), data, 0644); err != nil {
        log.WithField("path", newpath).WithField("oldpath", orgpath).WithError(err).Error("Failed to write problem.yml")
        return
    }
}
