package commands

import "fmt"

// VersionCommand prints the version output matching mxt.
func VersionCommand(version string) {
	fmt.Printf("mxt v%s\n", version)
}
