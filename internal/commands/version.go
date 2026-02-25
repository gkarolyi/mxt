package commands

import "fmt"

// VersionCommand prints the version output matching muxtree.
func VersionCommand(version string) {
	fmt.Printf("muxtree v%s\n", version)
}
