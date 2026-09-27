// Package selfupdate replaces only the CLI binary from a stable GitHub
// release after verifying its SHA-256 checksum. It never touches the rice.
package selfupdate

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/getaroli/desktop/cmd/aroli/internal/install"
)

const githubAPIRelease = "https://api.github.com/repos/getaroli/desktop/releases/latest"

type githubRelease struct {
	TagName string `json:"tag_name"`
}

var stableTagPattern = regexp.MustCompile(`^v(?:0|[1-9][0-9]*)\.(?:0|[1-9][0-9]*)\.(?:0|[1-9][0-9]*)$`)

// UpdateCLI implements `aroli cli update [--dry-run]`.
func UpdateCLI(version string, args []string) error {
	if len(args) == 0 || args[0] != "update" {
		return errors.New("use: aroli cli update [--dry-run]")
	}
	fs := flag.NewFlagSet("cli update", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	dryRun := fs.Bool("dry-run", false, "consulta a versão sem substituir o binário")
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}
	if len(fs.Args()) > 0 {
		return errors.New("use: aroli cli update [--dry-run]")
	}

	release, err := latestRelease(http.DefaultClient)
	if err != nil {
		return err
	}
	asset, err := cliAssetName(runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return err
	}
	base := "https://github.com/getaroli/desktop/releases/download/" + release.TagName + "/"
	fmt.Printf("CLI atual: %s\nDisponível: %s\n", version, strings.TrimPrefix(release.TagName, "v"))
	if *dryRun {
		fmt.Printf("seria baixado e verificado: %s%s\n", base, asset)
		return nil
	}

	checksums, err := download(http.DefaultClient, base+"SHA256SUMS.txt")
	if err != nil {
		return err
	}
	expected, err := checksumFor(asset, string(checksums))
	if err != nil {
		return err
	}
	binary, err := download(http.DefaultClient, base+asset)
	if err != nil {
		return err
	}
	actual := fmt.Sprintf("%x", sha256.Sum256(binary))
	if !strings.EqualFold(expected, actual) {
		return errors.New("a verificação de integridade falhou; a CLI não foi alterada")
	}
	if err := installDownloadedCLI(binary); err != nil {
		return err
	}
	fmt.Printf("Pronto! A CLI foi atualizada para %s.\n", strings.TrimPrefix(release.TagName, "v"))
	return nil
}

func latestRelease(client *http.Client) (githubRelease, error) {
	request, err := http.NewRequest(http.MethodGet, githubAPIRelease, nil)
	if err != nil {
		return githubRelease{}, err
	}
	request.Header.Set("Accept", "application/vnd.github+json")
	response, err := client.Do(request)
	if err != nil {
		return githubRelease{}, fmt.Errorf("não foi possível verificar atualizações da CLI: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return githubRelease{}, fmt.Errorf("GitHub respondeu %s ao consultar a CLI", response.Status)
	}
	var release githubRelease
	if err := json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(&release); err != nil {
		return githubRelease{}, err
	}
	if !stableTagPattern.MatchString(release.TagName) {
		return githubRelease{}, fmt.Errorf("tag de release inválida para atualização: %q", release.TagName)
	}
	return release, nil
}

func cliAssetName(goos, arch string) (string, error) {
	if goos != "linux" || (arch != "amd64" && arch != "arm64") {
		return "", fmt.Errorf("não há binário publicado para %s/%s", goos, arch)
	}
	return "aroli-linux-" + arch, nil
}

func download(client *http.Client, url string) ([]byte, error) {
	clientCopy := *client
	if clientCopy.Timeout == 0 {
		clientCopy.Timeout = 30 * time.Second
	}
	response, err := clientCopy.Get(url)
	if err != nil {
		return nil, fmt.Errorf("não foi possível baixar %s: %w", url, err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download falhou: %s", response.Status)
	}
	data, err := io.ReadAll(io.LimitReader(response.Body, 100<<20))
	if err != nil {
		return nil, err
	}
	return data, nil
}

func checksumFor(asset, contents string) (string, error) {
	for _, line := range strings.Split(contents, "\n") {
		fields := strings.Fields(line)
		if len(fields) >= 2 && strings.TrimPrefix(fields[len(fields)-1], "*") == asset {
			if len(fields[0]) != 64 || !isHex(fields[0]) {
				return "", fmt.Errorf("checksum inválido para %s", asset)
			}
			return fields[0], nil
		}
	}
	return "", fmt.Errorf("a release não publicou checksum para %s", asset)
}

func isHex(value string) bool {
	for _, r := range value {
		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')) {
			return false
		}
	}
	return true
}

func installDownloadedCLI(binary []byte) error {
	destination := cliInstallPath()
	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return err
	}
	temporary, err := os.CreateTemp(filepath.Dir(destination), ".aroli-download-")
	if err != nil {
		return err
	}
	name := temporary.Name()
	defer os.Remove(name)
	if _, err = temporary.Write(binary); err == nil {
		err = temporary.Chmod(0o755)
	}
	if closeErr := temporary.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return err
	}
	if err := os.Rename(name, destination); err != nil {
		return err
	}
	home, _ := os.UserHomeDir()
	if home != "" {
		return install.EnsureRiceShim(home)
	}
	return nil
}

func cliInstallPath() string {
	if executable, err := os.Executable(); err == nil {
		if resolved, err := filepath.EvalSymlinks(executable); err == nil {
			executable = resolved
		}
		if filepath.Base(executable) == "aroli" || filepath.Base(executable) == "aroli" {
			return executable
		}
	}
	home, err := os.UserHomeDir()
	if err == nil {
		return filepath.Join(home, ".local", "bin", "aroli")
	}
	return filepath.Join(".local", "bin", "aroli")
}
