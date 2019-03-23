package hook

import (
	"context"
)

type HookBuilder struct {
	hooks map[string]*hook
}

type PreHook func(context.Context, interface{}, func(context.Context, interface{}))
type PostHook func(context.Context, interface{})

type hook struct {
	preHook  []PreHook
	postHook []PostHook
}

func NewHookBuilder() *HookBuilder {
	builder := new(HookBuilder)
	builder.hooks = make(map[string]*hook)
	return builder
}

func (b *HookBuilder) DefineHook(key string) *HookBuilder {
	_, found := b.hooks[key]
	if found {
		panic("HookBuilder: DefineHook: Duplicate key")
	}
	b.hooks[key] = new(hook)
	return b
}

func (b *HookBuilder) PreHook(key string, h PreHook) *HookBuilder {
	hook, found := b.hooks[key]
	if !found {
		panic("HookBuilder: Hook: Key not found")
	}
	hook.preHook = append(hook.preHook, h)
	return b
}

func (b *HookBuilder) PostHook(key string, h PostHook) *HookBuilder {
	hook, found := b.hooks[key]
	if !found {
		panic("HookBuilder: Hook: Key not found")
	}
	hook.postHook = append(hook.postHook, h)
	return b
}

func (b *HookBuilder) Build() *HookInvoker {
	invoker := new(HookInvoker)
	for key, h := range b.hooks {
		ih := new(hook)
		ih.preHook = make([]PreHook, len(h.preHook))
		copy(ih.preHook, h.preHook)
		ih.postHook = make([]PostHook, len(h.postHook))
		copy(ih.postHook, h.postHook)
		invoker.hooks[key] = ih
	}
	return invoker
}

type HookInvoker struct {
	hooks map[string]*hook
}

type preHookFunc struct {
	hooks  []PreHook
	fn     func(context.Context, interface{})
	called bool
}

func (p *preHookFunc) invoke(ctx context.Context, val interface{}) {
	if p.called {
		panic("HookInvoker: invoke: fn called twice")
	} else {
		p.called = true
	}
	if len(p.hooks) == 0 {
		p.fn(ctx, val)
	} else {
		f := &preHookFunc{hooks: p.hooks[1:], fn: p.fn}
		p.hooks[0](ctx, val, f.invoke)
	}
}

func (h *HookInvoker) Invoke(ctx context.Context, key string, val interface{}, fn func(context.Context, interface{})) {
	hook, found := h.hooks[key]
	if !found {
		panic("HookInvoker: Invoke: Key not found")
	}
	f := &preHookFunc{hooks: hook.preHook, fn: fn}
	f.invoke(ctx, val)
}
