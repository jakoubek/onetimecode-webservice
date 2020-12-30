// +build mage

package main

import (
	"fmt"
	"os/exec"
	"path"

	"github.com/magefile/mage/sh"
)

const (
	BUILD_DIR     string = "./bin"
	BUILD_BINARY  string = "onetimecode"
	DEPLOY_TARGET string = "oli@opal5.opalstack.com:apps/onetimecode/onetimecode.new"
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

// Deploy to server
func Deploy() error {
	//mg.Deps(Build)
	fmt.Println("Deploying...")

	cmd := exec.Command("scp", path.Join(BUILD_DIR, BUILD_BINARY), DEPLOY_TARGET)
	err := cmd.Run()
	if err != nil {
		return err
	}
	fmt.Println("Copy to server ok.")

	return nil
}

// Clean up build directory
func Clean() {
	fmt.Println("Cleaning...")
	sh.Rm("bin")
}
