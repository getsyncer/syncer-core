package osstub

import (
	"fmt"
	"os"
)

type OsStub interface {
	Exit(code int)
	Getenv(key string) string
	Println(a ...interface{})
}

type osStub struct{}

func (o *osStub) Println(a ...interface{}) {
	fmt.Println(a...)
}

func (o *osStub) Exit(code int) {
	os.Exit(code)
}

func (o *osStub) Getenv(key string) string {
	return os.Getenv(key)
}

var _ OsStub = (*osStub)(nil)

type TestStub struct {
	Env      map[string]string
	Prints   []string
	ExitCode *int
}

func (t *TestStub) Println(a ...interface{}) {
	t.Prints = append(t.Prints, fmt.Sprintln(a...))
}

func (t *TestStub) Exit(code int) {
	t.ExitCode = &code
}

func (t *TestStub) Getenv(key string) string {
	return t.Env[key]
}

var _ OsStub = (*TestStub)(nil)
