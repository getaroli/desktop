// Package install runs the phased installer backend from a managed checkout.
package install

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/getaroli/desktop/cmd/aroli/internal/sys"
)

type installOptions struct {
	dryRun, yes, copy, resume bool
	lang                      string
	phases                    []string
}

// Install parses install flags and hands over to install.sh in the checkout.
func Install(args []string) error {
	fs := flag.NewFlagSet("install", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	var opts installOptions
	var repo string
	fs.BoolVar(&opts.dryRun, "dry-run", false, "não faz alterações")
	fs.BoolVar(&opts.dryRun, "n", false, "não faz alterações")
	fs.BoolVar(&opts.yes, "yes", false, "não pede confirmação")
	fs.BoolVar(&opts.yes, "y", false, "não pede confirmação")
	fs.BoolVar(&opts.copy, "copy", false, "copia em vez de criar links")
	fs.BoolVar(&opts.resume, "resume", false, "retoma a última instalação interrompida")
	fs.Bool("link", false, "cria links (padrão)")
	fs.StringVar(&opts.lang, "lang", sys.DefaultLang, "pt-BR, en ou es")
	fs.StringVar(&repo, "repo", "", "checkout do aroli")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if opts.lang != "pt-BR" && opts.lang != "en" && opts.lang != "es" {
		return errors.New("--lang aceita apenas pt-BR, en ou es")
	}
	opts.phases = fs.Args()
	if !opts.dryRun {
		if err := installSelf(); err != nil {
			return err
		}
	}
	path, cleanup, err := sys.EnsureRepositoryWithCleanup(repo, opts.dryRun)
	if err != nil {
		return err
	}
	defer cleanup()
	backend := []string{"--lang", opts.lang}
	if opts.dryRun {
		backend = append(backend, "--dry-run")
	}
	if opts.yes {
		backend = append(backend, "--yes")
	}
	if opts.copy {
		backend = append(backend, "--copy")
	}
	if opts.resume {
		backend = append(backend, "--resume")
	}
	backend = append(backend, opts.phases...)
	fmt.Printf("\n%s será instalado a partir de:\n  %s\n\n", sys.Project, path)
	return sys.Command(path, "bash", append([]string{"install.sh"}, backend...)...).Run()
}

// installSelf keeps ~/.local/bin/aroli under the CLI's control. The dotfile
// deployment intentionally skips that filename so an upgrade cannot replace
// the binary with the historical shell dispatcher.
func installSelf() error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	if runtime.GOOS == "windows" {
		return nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	dest := filepath.Join(home, ".local", "bin", "aroli")
	if sameFile(exe, dest) {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	src, err := os.Open(exe)
	if err != nil {
		return err
	}
	defer src.Close()
	tmp, err := os.CreateTemp(filepath.Dir(dest), ".aroli-")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	if _, err = io.Copy(tmp, src); err == nil {
		err = tmp.Chmod(0o755)
	}
	if closeErr := tmp.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return err
	}
	if err := os.Rename(tmpName, dest); err != nil {
		return err
	}
	return nil
}

func sameFile(a, b string) bool {
	aInfo, aErr := os.Stat(a)
	bInfo, bErr := os.Stat(b)
	return aErr == nil && bErr == nil && os.SameFile(aInfo, bInfo)
}
