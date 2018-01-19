package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/chzyer/readline"
)

// EnvVarTable is the global environment variable table.
type EnvVarTable struct {
	mux *sync.Mutex
	et  map[string]string
}

var (
	evt *EnvVarTable
)

// set sets an environment variable
func (et *EnvVarTable) set(envar, value string) {
	et.mux.Lock()
	et.et[envar] = value
	et.mux.Unlock()
}

// get returns the value of an environment variable
func (et *EnvVarTable) get(envar string) string {
	et.mux.Lock()
	val, ok := et.et[envar]
	et.mux.Unlock()
	if !ok {
		return ""
	}
	return val
}

// unset removes an environment variable
func (et *EnvVarTable) unset(envar string) {
	et.mux.Lock()
	delete(et.et, envar)
	et.mux.Unlock()
}

// init sets default env vars and loads some
// such as $PATH, $HOME, etc. from parent shell.
func (et *EnvVarTable) init() {
	// set defaults:
	et.set("SERVICE_PORT", "80")
	et.set("BINARY_IMAGE", "alpine:3.7")
	et.set("NODE_IMAGE", "node:9.4-alpine")
	et.set("PYTHON_IMAGE", "python:3.6-alpine3.7")
	et.set("RUBY_IMAGE", "ruby:2.5-alpine3.7")
	// load from parent shell, if present:
	val, ok := os.LookupEnv("KUBECTL_BINARY")
	if ok {
		et.set("KUBECTL_BINARY", val)
	}
	val, ok = os.LookupEnv("PATH")
	if ok {
		et.set("PATH", val)
	}
	val, ok = os.LookupEnv("HOME")
	if ok {
		et.set("HOME", val)
	}
}

func setprompt(rl *readline.Instance, context string) {
	namespace := "default"
	// res, err := kubectl("run", )
	// if err != nil {
	// 	warn("Can't determine namespace")
	// }
	rl.SetPrompt(fmt.Sprintf("[\033[32m%s\033[0m::\033[36m%s\033[0m]$ ", context, namespace))
}
