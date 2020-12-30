// +build mage

package main

import (
	"fmt"
	"path"
	"os/exec"

	_"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"github.com/magefile/mage/sh"
)

const (
	BUILD_DIR     string = "./bin"
	BUILD_BINARY  string = "onetimecode"
)

var (
	buildVersion string
)

// Default target to run when none is specified
// If not set, running mage will list available targets
var Default = Build

// Build executable
func Build() error {
	fmt.Println("Building...")

	versionCmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	version, err := versionCmd.Output()
	if err != nil {
		return err
	}
	buildVersion = string(version)

	cmd := exec.Command("go", "build", "-ldflags", "-X main.version="+buildVersion, "-o", path.Join(BUILD_DIR, BUILD_BINARY), ".")
	return cmd.Run()
}

// Clean up build directory
func Clean() {
	fmt.Println("Cleaning...")
	sh.Rm("bin")
}
