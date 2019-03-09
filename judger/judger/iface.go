package judger

import (
	"context"

	"github.com/syzoj/syzoj-ng-go/judger/rpc"
)

type ProblemBackend interface {
	JudgeSubmission(context.Context, *rpc.Task) error
}

var backends = make(map[string]ProblemBackend)

func RegisterBackend(name string, backend ProblemBackend) {
	if _, found := backends[name]; found {
		panic("Duplicate backend name: " + name)
	} else if name == "" {
		panic("Empty backend name")
	}
	backends[name] = backend
}
