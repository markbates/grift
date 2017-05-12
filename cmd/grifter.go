package cmd

import (
	"html/template"
	"os"
	"path/filepath"
	"sync"

	"github.com/gobuffalo/envy"
	"github.com/pkg/errors"
)

const exePath = "grift_runner.go"

var once = &sync.Once{}

type grifter struct {
	GriftsPackagePath string
	CommandName       string
	Verbose           bool
}

func newGrifter(name string) (*grifter, error) {
	g := &grifter{
		CommandName: name,
	}

	pwd, err := os.Getwd()
	if err != nil {
		return g, errors.WithStack(err)
	}

	stat, err := os.Stat(filepath.Join(pwd, "grifts"))
	if err != nil {
		if os.IsNotExist(err) {
			return g, errors.Errorf("there is no directory named 'grifts'. Run '%s init' or switch to the appropriate directory", name)
		}
		return g, err
	}

	if !stat.IsDir() {
		return g, errors.New("there should be a directory named 'grifts', not a file")
	}

	g.GriftsPackagePath = filepath.Join(envy.CurrentPackage(), "grifts")
	return g, nil
}

func (g *grifter) Setup() error {
	t, err := template.New("main").Parse(mainTmpl)
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
	return os.RemoveAll(exePath)
}
