package model

// The interface to implement for submission results with a score.
type SubmissionScore interface {
	Model_GetScore() (float64, error)
}
