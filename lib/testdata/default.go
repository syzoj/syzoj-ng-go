package testdata

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type defaultTestcase struct {
	Name   string
	Input  string
	Output string
}

func ParseDefault(path string) (*TestdataInfo, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	list := make(map[string]struct{})
	for _, file := range files {
		if !file.IsDir() {
			list[file.Name()] = struct{}{}
		}
	}
	var cases []defaultTestcase
	for name, _ := range list {
		if strings.HasSuffix(name, ".in") {
			caseName := name[:len(name)-3]
			var outFile string
			outFile1 := caseName + ".out"
			if _, exists := list[outFile1]; exists {
				outFile = outFile1
			} else {
				outFile2 := caseName + ".ans"
				if _, exists := list[outFile2]; exists {
					outFile = outFile2
				}
			}
			if outFile != "" {
				cases = append(cases, defaultTestcase{
					Name:   caseName,
					Input:  caseName + ".in",
					Output: outFile,
				})
			}
		}
	}
	if len(cases) == 0 {
		return nil, fmt.Errorf("cannot parse test data")
	}
	info := &TestdataInfo{}
	info.Cases = make(map[string]*Testcase, len(cases))
	caseScore := 100. / float64(len(cases))
	for _, testcase := range cases {
		info.Cases[testcase.Name] = &Testcase{
			Input:  testcase.Input,
			Output: testcase.Output,
		}
		info.Subtasks = append(info.Subtasks, &Subtask{
			Cases: []string{testcase.Name},
			Score: caseScore,
		})
	}
	return info, nil
}
