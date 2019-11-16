package grifts

import (
	"fmt"

	. "github.com/markbates/grift/grift"
)

var _ = Desc("hello", "Say Hello!")
var _ = Add("hello", func(c *Context) error {
	fmt.Println("Hello World!", c.Args)
	return nil
})

var _ = Namespace("db", func() {
	Add("seed", func(c *Context) error {
		fmt.Println("Seeding DB", c.Args)
		return nil
	})
})
