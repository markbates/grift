package grift

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"sync"
)

var taskList = map[string]Task{}
var descriptions = map[string]string{}
var lock = &sync.Mutex{}

type Task func(c *Context) error

func Add(name string, task Task) error {
	lock.Lock()
	defer lock.Unlock()
	if taskList[name] != nil {
		fn := taskList[name]
		taskList[name] = func(c *Context) error {
			err := fn(c)
			if err != nil {
				return err
			}
			return task(c)
		}
	} else {
		taskList[name] = task
	}
	return nil
}

func Set(name string, task Task) error {
	lock.Lock()
	defer lock.Unlock()
	taskList[name] = task
	return nil
}

func Rename(old string, new string) error {
	lock.Lock()
	defer lock.Unlock()
	if taskList[old] == nil {
		return fmt.Errorf("No task named %s defined!", old)
	}
	taskList[new] = taskList[old]
	delete(taskList, old)
	return nil
}

func Remove(name string) error {
	lock.Lock()
	defer lock.Unlock()
	delete(taskList, name)
	return nil
}

func Desc(name string, description string) error {
	lock.Lock()
	defer lock.Unlock()
	descriptions[name] = description
	return nil
}

func Run(name string, c *Context) error {
	if taskList[name] == nil {
		return fmt.Errorf("No task named %s defined!", name)
	}
	return taskList[name](c)
}

func TaskNames() []string {
	keys := []string{}
	for k := range taskList {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func Exec(args []string) error {
	name := "default"
	if len(args) >= 1 {
		name = args[0]
	}
	switch name {
	case "list":
		PrintTasks(os.Stdout)
	default:
		c := NewContext()
		c.Args = os.Args[2:]
		err := Run(name, c)
		if err != nil {
			log.Fatal(err)
		}
	}
	return nil
}

func PrintTasks(w io.Writer) {
	for _, k := range TaskNames() {
		m := fmt.Sprintf("grift run %s", k)
		desc := descriptions[k]
		if desc != "" {
			m = fmt.Sprintf("%s | %s", m, desc)
		}
		fmt.Fprintln(w, m)
	}
}
