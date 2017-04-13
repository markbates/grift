package grift

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
)

// Context used to pass information between grifts
type Context struct {
	context.Context
	Name    string
	Args    []string
	Verbose bool
	data    map[interface{}]interface{}
}

// Get a piece of data from the Context.
func (c *Context) Get(key string) interface{} {
	warningMsg := "Context#Get is deprecated. Please use Context#Value instead."

	_, file, no, ok := runtime.Caller(1)
	if ok {
		file = filepath.Base(file)
		warningMsg = fmt.Sprintf("Context#Get is deprecated. Please use Context#Value instead. Called from grifts/%s:%d", file, no)
	}
	log.Println(warningMsg)

	return c.data[key]
}

func (c *Context) Value(key interface{}) interface{} {
	return c.data[key]
}

// Set a piece of data onto the Context.
func (c *Context) Set(key string, val interface{}) {
	c.data[key] = val
}

// NewContext builds and returns a new default Context.
func NewContext(name string) *Context {
	return &Context{
		Context: context.Background(),
		Name:    name,
		Args:    []string{},
		data:    map[interface{}]interface{}{},
	}
}
