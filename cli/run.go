package cli

import (
	"context"
	"html/template"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/gobuffalo/here"
)

func Run(ctx context.Context, args []string) error {
	h := here.New()

	info, err := h.Current()
	if err != nil {
		return err
	}

	k := struct {
		Pkg     string
		Command string
		Dir     string
	}{
		Pkg:     path.Join(info.ImportPath, "grifts"),
		Command: "grift",
		Dir:     info.Dir,
	}

	if s, ok := ctx.Value("command").(string); ok {
		k.Command = s
	}

	od := filepath.Join(info.Dir, ".grifter")
	out := filepath.Join(od, "main.go")

	os.MkdirAll(od, 0755)
	defer os.RemoveAll(od)
	defer func() {
		if err := recover(); err != nil {
			os.RemoveAll(od)
		}
	}()

	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()

	mpl, err := template.New("main.go").Parse(tmpl)
	if err != nil {
		return err
	}

	if err := mpl.Execute(f, k); err != nil {
		return err
	}

	cargs := []string{"run", "-tags", "sqlite", out}
	cargs = append(cargs, args...)

	c := exec.CommandContext(ctx, "go", cargs...)
	c.Stdin = Stdin(ctx)
	c.Stdout = Stdout(ctx)
	c.Stderr = Stderr(ctx)

	if err := c.Run(); err != nil {
		return err
	}
	return nil
}
