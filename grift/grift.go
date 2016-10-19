package grift

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"sync"
	"time"
)

var griftList = map[string]Grift{}
var descriptions = map[string]string{}
var lock = &sync.Mutex{}

type Grift func(c *Context) error

func Add(name string, grift Grift) error {
	lock.Lock()
	defer lock.Unlock()
	if griftList[name] != nil {
		fn := griftList[name]
		griftList[name] = func(c *Context) error {
			err := fn(c)
			if err != nil {
				return err
			}
			return grift(c)
		}
	} else {
		griftList[name] = grift
	}
	return nil
}

func Set(name string, grift Grift) error {
	lock.Lock()
	defer lock.Unlock()
	griftList[name] = grift
	return nil
}

func Rename(old string, new string) error {
	lock.Lock()
	defer lock.Unlock()
	if griftList[old] == nil {
		return fmt.Errorf("No grift named %s defined!", old)
	}
	griftList[new] = griftList[old]
	delete(griftList, old)
	return nil
}

func Remove(name string) error {
	lock.Lock()
	defer lock.Unlock()
	delete(griftList, name)
	return nil
}

func Desc(name string, description string) error {
	lock.Lock()
	defer lock.Unlock()
	descriptions[name] = description
	return nil
}

func Run(name string, c *Context) error {
	if griftList[name] == nil {
		return fmt.Errorf("No grift named %s defined!", name)
	}
	if c.Verbose {
		defer func(start time.Time) {
			log.Printf("Completed grift %s in %s\n", name, time.Now().Sub(start))
		}(time.Now())
		log.Printf("Starting grift %s\n", name)
	}
	return griftList[name](c)
}

func GriftNames() []string {
	keys := []string{}
	for k := range griftList {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func Exec(args []string, verbose bool) error {
	name := "default"
	if len(args) >= 1 {
		name = args[0]
	}
	switch name {
	case "list":
		PrintGrifts(os.Stdout)
	default:
		c := NewContext()
		c.Verbose = verbose
		if len(args) >= 1 {
			c.Args = args[1:]
		}
		err := Run(name, c)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

func PrintGrifts(w io.Writer) {
	for _, k := range GriftNames() {
		m := fmt.Sprintf("grift %s", k)
		desc := descriptions[k]
		if desc != "" {
			m = fmt.Sprintf("%s | %s", m, desc)
		}
		fmt.Fprintln(w, m)
	}
}
