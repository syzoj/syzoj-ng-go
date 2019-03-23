package server

import (
	"github.com/syzoj/syzoj-ng-go/hook"
)

var HookBuilder = hook.NewHookBuilder()

func init() {
	HookBuilder.DefineHook("register")
}
