package common

import (
	"fmt"
	"os"
)

// printError prints error message to stderr and exits with code 1
func PrintError(message string, err error) {
	fmt.Fprintf(os.Stderr, "YAMLtecture\n")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n%v\n", message, err)
	} else {
		fmt.Fprintf(os.Stderr, "Error:\n%s\n", message)
	}
	os.Exit(1)
}
