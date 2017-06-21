package cmd

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

const exePath = ".grifter/main.go"

var once = &sync.Once{}

type grifter struct {
	GriftsPackagePath  string
	CommandName        string
	Verbose            bool
	GriftsAbsolutePath string
}

func hasGriftDir(path string) bool {
	stat, err := os.Stat(filepath.Join(path, "grifts"))
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return false
	}

	if !stat.IsDir() {
		return false
	}

	return true

}

func newGrifter(name string) (*grifter, error) {
	g := &grifter{
		CommandName: name,
	}

	path, err := os.Getwd()
	if err != nil {
		return g, errors.WithStack(err)
	}

	for !strings.HasSuffix(path, "/src") && path != "/" {
		if hasGriftDir(path) {
			break
		}

		path = filepath.Dir(path)
	}

	p := strings.SplitN(path, "/src/", 2)
	if len(p) == 1 {
		return g, errors.Errorf("There is no directory named 'grifts'. Run '%s init' or switch to the appropriate directory", name)
	}
	g.GriftsAbsolutePath = filepath.ToSlash(filepath.Join(path, "grifts"))
	g.GriftsPackagePath = filepath.ToSlash(filepath.Join(p[1], "grifts"))
	return g, nil
}

func (g *grifter) Setup() error {
	t, err := template.New("main").Parse(mainTmpl)
	if err != nil {
		return errors.WithStack(err)
	}

	err = os.MkdirAll(filepath.Dir(exePath), 0755)
	if err != nil {
		return errors.WithStack(err)
	}
	f, err := os.Create(exePath)
	if err != nil {
		return errors.WithStack(err)
	}

	err = t.Execute(f, g)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (g *grifter) TearDown() error {
	return os.RemoveAll(filepath.Dir(exePath))
}
