package utils

import "golang.org/x/exp/slices"

type HookEntry[TFunc any] struct {
	key   string
	hFunc TFunc
}

func (h *HookEntry[TFunc]) Key() string {
	return h.key
}

func (h *HookEntry[TFunc]) GetFunc() TFunc {
	return h.hFunc
}

func newHookEntry[TFunc any](key string, hFunc TFunc) *HookEntry[TFunc] {
	return &HookEntry[TFunc]{
		key:   key,
		hFunc: hFunc,
	}
}

type Hook[TFunc any] struct {
	hooks []*HookEntry[TFunc]
}

func NewHook[TFunc any]() *Hook[TFunc] {
	hook := &Hook[TFunc]{}
	hook.Clear()
	return hook
}
func (h *Hook[TFunc]) Add(key string, hFunc TFunc) {
	h.hooks = append(h.hooks, newHookEntry(key, hFunc))
}

func (h *Hook[TFunc]) Insert(key string, hFunc TFunc, position int) {
	if position < 0 {
		position = 0
	}
	if position > len(h.hooks) {
		position = len(h.hooks)
	}
	h.hooks = slices.Insert(h.hooks, position, newHookEntry(key, hFunc))
}

func (h *Hook[TFunc]) GetEntries() []*HookEntry[TFunc] {
	return h.hooks
}

func (h *Hook[TFunc]) GetEntry(key string) []*HookEntry[TFunc] {
	entries := make([]*HookEntry[TFunc], 0)
	for _, hook := range h.hooks {
		if hook != nil && hook.key == key {
			entries = append(entries, hook)
		}
	}
	return entries
}

func (h *Hook[TFunc]) Remove(position int) {
	h.hooks = slices.Delete(h.hooks, position, position+1)
}

func (h *Hook[TFunc]) RemoveByKey(key string) {
	entries := make([]*HookEntry[TFunc], len(h.hooks))
	i := 0
	for _, hook := range h.hooks {
		if hook != nil && hook.key != key {
			entries[i] = hook
			i++
		}
	}
	h.hooks = entries[:i]
}

func (h *Hook[TFunc]) Clear() {
	h.hooks = make([]*HookEntry[TFunc], 0)
}
