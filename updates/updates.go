package updates

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"text/template"

	"github.com/hverr/go-updater"
	"github.com/hverr/status-dashboard/version"
	"github.com/kardianos/osext"
)

const (
	RepositoryOwner = "hverr"
	RepositoryName  = "status-dashboard"
	AssetTemplate   = "{{.BaseName}}_{{.OS}}_{{.Arch}}"
)

func CheckForUpdates(baseName string) {
	u := appUpdater(baseName, nil)

	fmt.Println("[+] Checking for updates...")
	latest, err := u.Check()
	if err != nil {
		fmt.Println("[-] Error:", err)
		fmt.Println("[-] Failed to check for updates.")
		os.Exit(1)
	}

	if latest == nil {
		fmt.Println("[*] This version is up to date.")
		os.Exit(0)
	}

	fmt.Printf("[*] A new version is available: %v (%v):\n", latest.Name(), latest.Identifier())
	fmt.Println(latest.Information())
	fmt.Println("")

	if version.HasVersionInformation() {
		fmt.Print("[*] Current version: ")
		version.PrintVersionInformation(os.Stdout)
	} else {
		fmt.Println("[*] No version is available for the current binary.")
	}
}

func UpdateApp(baseName string) {
	path, err := osext.Executable()
	if err != nil || path == "" {
		if err != nil {
			fmt.Println("[-] Error:", err)
		}
		fmt.Println("[-] Failed to get the executable location.")
		os.Exit(1)
	}
	f := updater.NewDelayedFile(path)
	u := appUpdater(baseName, f)

	fmt.Println("[+] Checking for updates...")
	latest, err := u.Check()
	if err != nil {
		fmt.Println("[-] Error:", err)
		fmt.Println("[-] Failed to check for updates.")
		os.Exit(1)
	}

	if latest == nil {
		fmt.Println("[*] This version is up to date.")
		os.Exit(0)
	}

	fmt.Printf("[+] Downloading %v (%v):\n", latest.Name(), latest.Identifier())
	fmt.Println(latest.Information())
	fmt.Println("")

	err = u.UpdateTo(latest)
	if err != nil {
		fmt.Println("[-] Error:", err)
		fmt.Println("[-] An error ocurred while updating.")
		os.Exit(1)
	}

	fmt.Println("[+] Replacing", path)
	err = f.Close()
	if err != nil {
		fmt.Println("[-] Error:", err)
		fmt.Println("[-] An error ocurred while updating.")
		os.Exit(1)
	}

	fmt.Println("[+] Successfully updated the binary.")
}

func appUpdater(baseName string, f updater.AbortWriter) updater.Updater {
	return updater.Updater{
		App: updater.NewGitHub("hverr", "status-dashboard", nil),
		CurrentReleaseIdentifier: version.Commit,
		WriterForAsset:           assetWriter(baseName, f),
	}
}

func assetWriter(baseName string, f updater.AbortWriter) func(updater.Asset) (updater.AbortWriter, error) {
	return func(a updater.Asset) (updater.AbortWriter, error) {
		if a.Name() != assetName(baseName) {
			return nil, nil
		}
		return f, nil
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
