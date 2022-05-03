//go:build mage
// +build mage

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/zhiminwen/magetool/sshkit"

	"github.com/magefile/mage/mg"
	sh "github.com/magefile/mage/sh"
)

const (
	BUILD_DIR          string = "./bin"
	BUILD_BINARY       string = "onetimecode-server"
	BUILD_BINARY_LOCAL string = "onetimecode-server-debug.exe"
	DEPLOY_TARGET      string = "oli@opal5.opalstack.com:apps/onetimecode/onetimecode-server.new"
	DEPLOY_DIR         string = "oli@opal5.opalstack.com:apps/onetimecode/"
	SSH_CONNECT        string = "oli@opal5.opalstack.com"
)

var (
	buildVersion string
	fullCommit   string
	buildTime    string
)

func setBuildVariables() {
	versionCmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	versionCmdOutput, err := versionCmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	buildVersion = string(versionCmdOutput)

	versionCmd = exec.Command("git", "rev-parse", "HEAD")
	versionCmdOutput, err = versionCmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	fullCommit = string(versionCmdOutput)

	buildTime = time.Now().Format("2006-01-02_15:04:05")
}

// restart the newly deployed webserver binary on the production server
func RestartService() error {

	fmt.Println("Restarting webserver on production server...")

	productionServer, err := sshkit.NewSSHClient("opal5.opalstack.com", "22", "oli", "", "C:\\Users\\jakoubek\\.ssh\\id_rsa")
	if err != nil {
		return err
	}

	err = productionServer.Execute("cd apps/onetimecode && ./restart")
	if err != nil {
		return err
	}

	err = productionServer.Close()
	if err != nil {
		return err
	}

	fmt.Println("Restart finished.")
	return nil

}

// deploy to production server
func Deploy() error {

	mg.Deps(Build)

	binaryPath := path.Join(BUILD_DIR, BUILD_BINARY)

	fmt.Println("Deploying to production server...")

	fmt.Println("Copying webserver binary...")
	err := sh.Run("scp", binaryPath, DEPLOY_TARGET)
	if err != nil {
		return err
	}

	fmt.Println("Deployment finished.")
	return nil
}

// build the binary
func Build() error {

	mg.Deps(Prepare)

	buildPath := path.Join(BUILD_DIR, BUILD_BINARY)

	setBuildVariables()

	fmt.Printf("Building %s...\n", buildPath)

	return sh.RunWith(map[string]string{"GOOS": "linux"}, "go", "build", "-ldflags", "-s -X main.buildVersion="+buildVersion+" -X main.fullCommit="+fullCommit+" -X main.buildTime="+buildTime, "-o", buildPath, ".")

}

// building and running locally
func Debugrun() error {

	mg.Deps(Prepare)

	buildPath := path.Join(BUILD_DIR, BUILD_BINARY_LOCAL)

	setBuildVariables()

	fmt.Printf("Building and running locally %s...\n", buildPath)

	sh.RunWith(map[string]string{"GOOS": "windows"}, "go", "build", "-ldflags", "-s -X main.buildVersion="+buildVersion+" -X main.fullCommit="+fullCommit+" -X main.buildTime="+buildTime+" -X main.isDebugMode=true", "-o", buildPath, "./cmd/server/")
	//sh.RunWith(map[string]string{"GOOS": "windows"}, "go", "build", "-o", buildPath, "./cmd/server/")

	//cmd := exec.Command("go", "build", "-o", buildPath, "./cmd/server/")
	//cmd.Run()

	err := sh.Run(buildPath)

	//cmd = exec.Command(buildPath)

	//return cmd.Run()
	return err
}

// Prepare directory for builds
func Prepare() {
	fmt.Printf("Prepare %s directory...\n", BUILD_DIR)
	if err := os.Mkdir(BUILD_DIR, os.ModePerm); err != nil {
		log.Printf("Creating %s directory didn't work: ", BUILD_DIR, err.Error())
	}
}

// Clean up after yourself
func Clean() {
	fmt.Println("Cleaning...")
	os.RemoveAll("bin")
}
