package cmd

import (
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/pkg/errors"
)

var once = &sync.Once{}

type grifter struct {
	CurrentDir        string
	BuildPath         string
	GriftsPackagePath string
	ExePath           string
	CommandName       string
	Verbose           bool
}

func newGrifter(name string) (*grifter, error) {
	g := &grifter{
		CommandName: name,
	}

	pwd, err := os.Getwd()
	if err != nil {
		return g, err
	}
	g.CurrentDir = pwd

	stat, err := os.Stat(filepath.Join(pwd, "grifts"))
	if err != nil {
		if os.IsNotExist(err) {
			return g, errors.Errorf("There is no directory named 'grifts'. Run '%s init' or switch to the appropriate directory.", name)
		}
		return g, err
	}

	if !stat.IsDir() {
		return g, errors.New("There should be a directory named 'grifts', not a file.")
	}

	base := filepath.Base(g.CurrentDir)
	g.BuildPath = filepath.Join(os.Getenv("GOPATH"), "src", "grift.build", base)
	g.GriftsPackagePath = filepath.Join("grift.build", base, "grifts")
	return g, nil
}

func (g *grifter) Setup() error {
	err := os.MkdirAll(g.BuildPath, 0777)
	if err != nil {
		return err
	}

	return g.Build()
}

func (g *grifter) Build() error {
	err := g.copyGrifts()
	if err != nil {
		return err
	}

	t, err := template.New("main").Parse(mainTmpl)
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(g.BuildPath, "main.go"))
	if err != nil {
		return err
	}

	err = t.Execute(f, g)
	if err != nil {
		return err
	}

	g.ExePath = filepath.Join(g.BuildPath, "main.go")
	return nil
}

func (g *grifter) TearDown() error {
	return os.RemoveAll(g.BuildPath)
}

func (g *grifter) copyGrifts() error {
	var err error
	cp := exec.Command("cp", "-rv", filepath.Join(g.CurrentDir, "grifts"), g.BuildPath)
	cp.Stderr = os.Stderr
	err = cp.Run()
	if err != nil {
		return errors.WithStack(err)
	}

	once.Do(func() {
		vdir := filepath.Join(g.CurrentDir, "vendor")
		if _, err := os.Stat(vdir); err == nil {
			fmt.Println("Vendor directory found. Please be aware that this can slow down running of tasks.")
			bdir := filepath.Join(g.BuildPath, "vendor")
			err = os.RemoveAll(bdir)
			if err != nil {
				return
			}
			cp = exec.Command("ln", "-s", vdir, bdir)
			cp.Stderr = os.Stderr
			err = cp.Run()
		}
	})

	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
