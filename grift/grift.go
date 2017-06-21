package grift

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

var CommandName = "grift"
var griftList = map[string]Grift{}
var descriptions = map[string]string{}
var lock = &sync.Mutex{}
var namespace string
var maxNameLen int

type Grift func(c *Context) error

// Create a namespace. All tasks within the
// namespace will share the same prefix.
func Namespace(name string, s func()) error {
	defer func() {
		namespace = ""
	}()

	namespace = applyNamespace(name)
	s()
	return nil
}

func applyNamespace(name string) string {
	if namespace != "" {
		if strings.HasPrefix(name, ":") {
			return name[1:]
		}
		if name == "default" {
			return name
		}
		return fmt.Sprintf("%s:%s", namespace, name)

	}

	return name
}

// Add a grift. If there is already a grift
// with the given name the two grifts will
// be bundled together.
func Add(name string, grift Grift) error {
	lock.Lock()
	defer lock.Unlock()

	name = applyNamespace(name)
	if len(name) > maxNameLen {
		maxNameLen = len(name)
	}

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

// Set a grift. This is similar to `Add` but it will
// overwrite an existing grift with the same name.
func Set(name string, grift Grift) error {
	lock.Lock()
	defer lock.Unlock()
	name = applyNamespace(name)
	griftList[name] = grift
	return nil
}

// Rename a grift. Useful if you want to re-define
// an existing grift, but don't want to write over
// the original.
func Rename(old string, new string) error {
	lock.Lock()
	defer lock.Unlock()

	old = applyNamespace(old)
	new = applyNamespace(new)

	if griftList[old] == nil {
		return fmt.Errorf("No task named %s defined!", old)
	}
	griftList[new] = griftList[old]
	delete(griftList, old)
	return nil
}

// Remove a grift. Not incredibly useful, but here for
// completeness.
func Remove(name string) error {
	lock.Lock()
	defer lock.Unlock()

	name = applyNamespace(name)

	delete(griftList, name)
	delete(descriptions, name)
	return nil
}

// Desc sets a helpful descriptive text for a grift.
// This description will be shown when `grift list`
// is run.
func Desc(name string, description string) error {
	lock.Lock()
	defer lock.Unlock()

	name = applyNamespace(name)

	descriptions[name] = description
	return nil
}

// Run a grift. This allows for the chaining for grifts.
// One grift can Run another grift and so on.
func Run(name string, c *Context) error {
	name = applyNamespace(name)

	if griftList[name] == nil {
		return fmt.Errorf("No task named '%s' defined!", name)
	}
	if c.Verbose {
		defer func(start time.Time) {
			log.Printf("Completed task %s in %s\n", name, time.Now().Sub(start))
		}(time.Now())
		log.Printf("Starting task %s\n", name)
	}
	return griftList[name](c)
}

// List of the names of the defined grifts.
func List() []string {
	keys := []string{}
	for k := range griftList {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Exec the grift stack. This is the main "entry point" to
// the grift system.
func Exec(args []string, verbose bool) error {
	name := "default"
	if len(args) >= 1 {
		name = args[0]
	}
	switch name {
	case "list":
		PrintGrifts(os.Stdout)
	default:
		c := NewContext(name)
		c.Verbose = verbose
		if len(args) >= 1 {
			c.Args = args[1:]
		}
		return Run(name, c)
	}
	return nil
}

// PrintGrifts to the screen, nice, sorted, and with descriptions,
// should they exist.
func PrintGrifts(w io.Writer) {
	for _, k := range List() {
		m := fmt.Sprintf("%s %s", CommandName, k)
		desc := descriptions[k]
		if desc != "" {
			m = fmt.Sprintf("%s | %s", m, desc)
		}
		fmt.Fprintln(w, m)
	}
}
