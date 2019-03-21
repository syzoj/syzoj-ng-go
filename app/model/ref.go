package model

type UserRef string

func (u UserRef) Test() string {
	return string(u)
}
