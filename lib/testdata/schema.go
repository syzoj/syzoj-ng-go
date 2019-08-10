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
	Input        *File  `json:"input,omitempty"`
	Output       *File  `json:"output,omitempty"`
	Answer       string `json:"answer,omitempty"`
	SpecialJudge string `json:"special_judge,omitempty"`
}

type SpecialJudge struct {
	Language string `json:"language,omitempty"`
	File     *File  `json:"file,omitempty"`
}

type Language struct {
}

type Subtask struct {
	Cases []string `json:"cases"`
	Score float64  `json:"score"`
}

type File struct {
	Name      string `json:"name"`
	Sha256Sum string `json:"sha256sum"`
}
