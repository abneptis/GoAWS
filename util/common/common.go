package common

import (
	"os"
)

type FunctionModule struct {
	FlagFunc func()
	Setup    func() os.Error
	Calls    map[string]UserCall
}

func (self FunctionModule) Names() (out []string) {
	for k, _ := range self.Calls {
		out = append(out, k)
	}
	return
}

func NewFunctionModule(f func()) *FunctionModule {
	return &FunctionModule{
		FlagFunc: f,
		Calls:    make(map[string]UserCall),
	}
}

var Modules map[string]*FunctionModule

// NOT THREAD SAFE!
// DO NOT USE OUTSIDE OF init()!
func AddModule(mod string, ffunc func()) {
	Modules[mod] = NewFunctionModule(ffunc)
}

// Common functionality to ease sub-modules

type UserCall func([]string) os.Error


func init() {
	Modules = make(map[string]*FunctionModule)
}
