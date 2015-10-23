package version

import (
	"fmt"
	"io"
)

var Commit string
var Dirty string

func IsDirty() bool {
	return Dirty == "yes"
}

func HasVersionInformation() bool {
	return Commit != ""
}

func PrintVersionInformation(w io.Writer) bool {
	if !HasVersionInformation() {
		fmt.Fprintln(w, "No version information available.")
		return false
	}

	if !IsDirty() {
		fmt.Fprintln(w, "Commit", Commit)
	} else {
		fmt.Fprintln(w, "Commit", Commit, "(dirty)")
	}
	return true
}
