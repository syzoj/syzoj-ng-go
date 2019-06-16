package model

type Problem struct {
	Id        string   `json:"id"`
	Title     string   `json:"title"`
	Statement string   `json:"statement"`
	Tags      []string `json:"tags"`
	Name      *string  `json:"name"`
}

type ProblemDoc struct {
	Title     string   `json:"title"`
	Statement string   `json:"statement"`
	Tags      []string `json:"tags"`
}
