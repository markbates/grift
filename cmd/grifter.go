package cmd

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gobuffalo/envy"
	"github.com/pkg/errors"
	"github.com/rogpeppe/go-internal/modfile"
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

	currentPath, err := os.Getwd()
	if err != nil {
		return g, errors.WithStack(err)
	}

	if envy.InGoPath() {
		// pwd is inside a GOPATH
		for !strings.HasSuffix(currentPath, "/src") && currentPath != "/" {
			if hasGriftDir(currentPath) {
				break
			}

			currentPath = filepath.Dir(currentPath)
		}
		p := strings.SplitN(currentPath, filepath.FromSlash("/src/"), 2)
		if len(p) == 1 {
			return g, errors.Errorf("There is no directory named 'grifts'. Run '%s init' or switch to the appropriate directory", name)
		}
		g.GriftsAbsolutePath = filepath.ToSlash(filepath.Join(currentPath, "grifts"))
		g.GriftsPackagePath = filepath.ToSlash(filepath.Join(p[1], "grifts"))
	} else {
		// is outside of gopath, dont loop to parent
		if !hasGriftDir(currentPath) {
			return g, errors.Errorf("There is no directory named 'grifts'. Run '%s init' or switch to the appropriate directory", name)
		}
		g.GriftsAbsolutePath = filepath.ToSlash(filepath.Join(currentPath, "grifts"))

		// check for go module to see if we can get go.mod
		if envy.Mods() {
			moddata, err := ioutil.ReadFile("go.mod")
			if err != nil {
				return g, errors.New("go.mod cannot be read or does not exist while go module is enabled.")
			}
			packagePath := modfile.ModulePath(moddata)
			if packagePath == "" {
				return g, errors.New("go.mod is malformed.")
			}
			g.GriftsPackagePath = fmt.Sprintf("%s/grifts", packagePath)
		} else {
			// no go module, infer package path from current directory
			g.GriftsPackagePath = filepath.ToSlash(filepath.Join(path.Base(currentPath), "grifts"))
		}

	}

	return g, nil
}

func (g *grifter) Setup() error {

	root := filepath.Dir(exePath)
	err := os.MkdirAll(root, 0755)
	if err != nil {
		return errors.WithStack(err)
	}

	tmpls := map[string]string{}
	tmpls[exePath] = mainTmpl
	if envy.Mods() {
		tmpls[filepath.Join(root, "go.mod")] = modTmpl
	}
	for k, v := range tmpls {
		t, err := template.New(k).Parse(v)
		if err != nil {
			return errors.WithStack(err)
		}

		f, err := os.Create(k)
		if err != nil {
			return errors.WithStack(err)
		}

		err = t.Execute(f, g)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (g *grifter) TearDown() error {
	return os.RemoveAll(filepath.Dir(exePath))
}
