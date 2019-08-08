package testdata

type TestdataInfo struct {
	Cases         map[string]*Testcase     `json:"cases"`
	SpecialJudges map[string]*SpecialJudge `json:"special_judges"`
	Languages     map[string]*Language     `json:"languages"`
	Subtasks      []*Subtask               `json:"subtasks"`
}

type Testcase struct {
	MemoryLimit  int64  `json:"memory_limit,omitempty"`
	TimeLimit    int64  `json:"time_limit,omitempty"`
	Input        string `json:"input,omitempty"`
	Output       string `json:"output,omitempty"`
	Answer       string `json:"answer,omitempty"`
	SpecialJudge string `json:"special_judge,omitempty"`
}

type SpecialJudge struct {
	Language string `json:"language,omitempty"`
	FileName string `json:"filename,omitempty"`
}

type Language struct {
}

type Subtask struct {
	Cases []string `json:"cases"`
	Score float64  `json:"score"`
}
