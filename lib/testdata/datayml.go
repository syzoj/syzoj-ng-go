package testdata

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-yaml/yaml"
)

type dataYml struct {
	InputFile    string               `yaml:"inputFile"`
	OutputFile   string               `yaml:"outputFile"`
	UserOutput   string               `yaml:"userOutput"`
	Subtasks     []*dataYmlSubtask    `yaml:"subtasks"`
	SpecialJudge *dataYmlSpecialJudge `yaml:"specialJudge"`
}

type dataYmlSubtask struct {
	Score float64  `yaml:"score"`
	Type  string   `yaml:"type"` // sum, min, mul
	Cases []string `yaml:"cases"`
}

type dataYmlSpecialJudge struct {
	Language string `yaml:"language"`
	FileName string `yaml:"fileName"`
}

func ParseDataYml(path string) (*TestdataInfo, error) {
	filePath := filepath.Join(path, "data.yml")
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	fileReader := io.LimitReader(file, 100*1024)
	var data dataYml
	err = yaml.NewDecoder(fileReader).Decode(&data)
	if err != nil {
		return nil, err
	}

	info := &TestdataInfo{}
	info.SpecialJudges = make(map[string]*SpecialJudge)
	var spjName string
	if data.SpecialJudge != nil {
		spjName = "Spj"
		dataSpj := data.SpecialJudge
		info.SpecialJudges[spjName] = &SpecialJudge{
			Language: dataSpj.Language,
			FileName: dataSpj.FileName,
		}
	}
	info.Cases = make(map[string]*Testcase)
	for _, dataSubtask := range data.Subtasks {
		for _, dataCase := range dataSubtask.Cases {
			if _, exists := info.Cases[dataCase]; !exists {
				info.Cases[dataCase] = &Testcase{
					Input:        strings.ReplaceAll(data.InputFile, "#", dataCase),
					Output:       strings.ReplaceAll(data.OutputFile, "#", dataCase),
					Answer:       strings.ReplaceAll(data.UserOutput, "#", dataCase),
					SpecialJudge: spjName,
				}
			}
		}
		if dataSubtask.Type == "sum" && len(dataSubtask.Cases) != 0 {
			// Split into multiple subtasks
			score := dataSubtask.Score / float64(len(dataSubtask.Cases))
			for _, dataCase := range dataSubtask.Cases {
				infoSubtask := &Subtask{}
				infoSubtask.Cases = []string{dataCase}
				infoSubtask.Score = score
				info.Subtasks = append(info.Subtasks, infoSubtask)
			}
		} else {
			infoSubtask := &Subtask{}
			infoSubtask.Cases = dataSubtask.Cases
			infoSubtask.Score = dataSubtask.Score
			info.Subtasks = append(info.Subtasks, infoSubtask)
		}
	}
	return info, nil
}
