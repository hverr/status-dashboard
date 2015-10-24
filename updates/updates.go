package updates

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"text/template"

	"github.com/hverr/go-updater"
	"github.com/hverr/status-dashboard/version"
)

const (
	RepositoryOwner = "hverr"
	RepositoryName  = "status-dashboard"
	AssetTemplate   = "{{.BaseName}}_{{.OS}}_{{.Arch}}"
)

func CheckForUpdates(baseName string) {
	u := appUpdater(baseName)

	fmt.Println("Checking for updates...")
	latest, err := u.Check()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		fmt.Println("Failed to check for updates.")
		os.Exit(1)
	}

	if latest == nil {
		fmt.Println("This version is up to date.")
		os.Exit(0)
	}

	fmt.Printf("A new version is availeble: %v (%v):\n", latest.Name(), latest.Identifier())
	fmt.Println(latest.Information())
	fmt.Println("")

	if version.HasVersionInformation() {
		fmt.Print("Current version: ")
		version.PrintVersionInformation(os.Stdout)
	} else {
		fmt.Println("No version is available for the current binary.")
	}
}

func UpdateApp(baseName string) {
	u := appUpdater(baseName)

	fmt.Println("Checking for updates...")
	latest, err := u.Check()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		fmt.Println("Failed to check for updates.")
		os.Exit(1)
	}

	if latest == nil {
		fmt.Println("This version is up to date.")
		os.Exit(0)
	}

	fmt.Printf("Updating to: %v (%v):\n", latest.Name(), latest.Identifier())
	fmt.Println(latest.Information())
	fmt.Println("")

	err = u.UpdateTo(latest)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		fmt.Println("An error ocurred while updating.")
		os.Exit(1)
	}

	fmt.Printf("Successfully updated the binary at %v\n", os.Args[0])
}

func appUpdater(baseName string) updater.Updater {
	return updater.Updater{
		App: updater.NewGitHub("hverr", "status-dashboard", nil),
		CurrentReleaseIdentifier: version.Commit,
		WriterForAsset:           assetWriter(baseName),
	}
}

func assetWriter(baseName string) func(updater.Asset) (updater.AbortWriter, error) {
	return func(a updater.Asset) (updater.AbortWriter, error) {
		if a.Name() != assetName(baseName) {
			return nil, nil
		}
		return updater.NewDelayedFile(os.Args[0]), nil
	}
}

func assetName(baseName string) string {
	data := struct {
		BaseName string
		OS       string
		Arch     string
	}{
		BaseName: baseName,
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
	}

	t, _ := template.New("asset").Parse(AssetTemplate)

	b := bytes.NewBuffer(nil)
	t.Execute(b, &data)

	return b.String()
}
