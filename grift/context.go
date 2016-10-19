package grift

type Context struct {
	Args    []string
	Verbose bool
	data    map[string]interface{}
}

func (c *Context) Get(key string) interface{} {
	return c.data[key]
}

func (c *Context) Set(key string, val interface{}) {
	c.data[key] = val
}

func NewContext() *Context {
	return &Context{
		Args: []string{},
		data: map[string]interface{}{},
	}
}
