package main

import (
	"context"
	"fmt"

	"dagger/itsfilmnoir/internal/dagger"
)

func (m *Itsfilmnoir) createBaseContainer() *dagger.Container {
	return dag.Container().From("homebrew/brew:4.4.1")
}

func (m *Itsfilmnoir) installTools(container *dagger.Container) *dagger.Container {
	return container.WithExec([]string{"brew", "install", "skopeo", "yq", "jq", "go"})
}

func (m *Itsfilmnoir) CreateVersionUpdaterContainer() *dagger.Container {
	baseContainer := m.createBaseContainer()
	containerWithTools := m.installTools(baseContainer)
	return containerWithTools
}

func (m *Itsfilmnoir) UpdateVersions(ctx context.Context, source *dagger.Directory) (string, error) {
	updaterContainer := m.CreateVersionUpdaterContainer()

	updaterProgram := `
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
)

func main() {
	if err := updateVersions(); err != nil {
		fmt.Fprintf(os.Stderr, "Error updating versions: %v\n", err)
		os.Exit(1)
	}
}

func updateVersions() error {
	updates := []struct {
		file    string
		image   string
		filter  string
		pattern string
	}{
		{"/src/dagger/aws.go", "amazon/aws-cli", "^2", "amazon/aws-cli:[0-9]+\\.[0-9]+\\.[0-9]+"},
		{"/src/dagger/prettier.go", "node", "lts-alpine", "node:lts-alpine[0-9]+\\.[0-9]+\\.[0-9]+"},
		{"/src/dagger/gofumpt.go", "homebrew/brew", "^[0-9]", "homebrew/brew:[0-9]+\\.[0-9]+\\.[0-9]+"},
	}

	for _, update := range updates {
		newVersion, err := getLatestVersion(update.image, update.filter)
		if err != nil {
			return fmt.Errorf("error getting latest version for %s: %w", update.image, err)
		}

		if err := updateVersion(update.file, update.pattern, newVersion); err != nil {
			return fmt.Errorf("error updating version in %s: %w", update.file, err)
		}

		fmt.Printf("Updated %s to version %s\n", update.image, newVersion)
	}

	daggerVersion, err := getGithubLatestRelease("dagger/dagger")
	if err != nil {
		return fmt.Errorf("error getting latest Dagger version: %w", err)
	}

	if err := updateDaggerVersion(daggerVersion); err != nil {
		return fmt.Errorf("error updating Dagger version: %w", err)
	}

	fmt.Printf("Updated Dagger to version %s\n", daggerVersion)

	return nil
}

func getLatestVersion(image, filter string) (string, error) {
	url := fmt.Sprintf("https://registry.hub.docker.com/v2/repositories/%s/tags", image)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		Results []struct {
			Name string ` + "`json:\"name\"`" + `
		} ` + "`json:\"results\"`" + `
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	var versions semver.Collection
	for _, tag := range result.Results {
	    fmt.Printf("%s", tag.Name)
		if matched, _ := regexp.MatchString(filter, tag.Name); matched {
			if v, err := semver.NewVersion(tag.Name); err == nil {
				versions = append(versions, v)
			}
		}
	}

	if len(versions) == 0 {
		return "", fmt.Errorf("no matching versions found")
	}

	sort.Sort(sort.Reverse(versions))
	return versions[0].String(), nil
}

func updateVersion(file, pattern, newVersion string) error {
	content, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	re := regexp.MustCompile(pattern)
	updatedContent := re.ReplaceAllString(string(content), strings.Split(pattern, ":")[0]+":"+newVersion)

	return os.WriteFile(file, []byte(updatedContent), 0644)
}

func getGithubLatestRelease(repo string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		TagName string ` + "`json:\"tag_name\"`" + `
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	return strings.TrimPrefix(result.TagName, "v"), nil
}

func updateDaggerVersion(version string) error {
	cmd := fmt.Sprintf("go get dagger.io/dagger@v%s", version)
	return executeCommand(cmd)
}

func executeCommand(command string) error {
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
`

	output, err := updaterContainer.
		WithMountedDirectory("/src", source).
		WithWorkdir("/tmp/updater").
		WithNewFile("updater.go", updaterProgram).
		WithExec([]string{"sudo", "chmod", "-R", "a+rwx", "."}).
		WithExec([]string{"go", "mod", "init", "updater"}).
		WithExec([]string{"go", "get", "github.com/Masterminds/semver/v3@latest"}).
		Terminal().
		WithExec([]string{"go", "run", "updater.go"}).
		Stdout(ctx)

	if err != nil {
		return "", fmt.Errorf("error running updater: %w", err)
	}

	return output, nil
}
