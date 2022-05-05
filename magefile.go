//go:build mage
// +build mage

package main

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"github.com/magefile/mage/mg"
	sh "github.com/magefile/mage/sh"
	"github.com/zhiminwen/magetool/sshkit"
	"log"
	"os"
	"path"
	"time"
)

var (
	buildDir         string
	buildBinary      string
	buildBinaryLocal string
	deployTarget     string
	deployDir        string
	sshConnect       string
	sshPort          string
	sshUser          string
	sshKeyfile       string
	statsApiUrl      string
	secureKey        string

	buildTime string
)

var Default = Debugrun

// restart the newly deployed webserver binary on the production server
func RestartService() error {

	mg.Deps(LoadEnvironment)

	fmt.Println("Restarting webserver on production server...")

	productionServer, err := sshkit.NewSSHClient(sshConnect, sshPort, sshUser, "", sshKeyfile)
	if err != nil {
		return err
	}

	err = productionServer.Execute("cd apps/onetimecode && ./update")
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

	binaryPath := path.Join(buildDir, buildBinary)

	fmt.Println("Deploying to production server...")

	fmt.Println("Copying webserver binary...")
	err := sh.Run("scp", binaryPath, deployTarget)
	if err != nil {
		return err
	}

	fmt.Println("Deployment finished.")
	return nil
}

// build the binary
func Build() error {

	mg.Deps(Prepare)

	buildPath := path.Join(buildDir, buildBinary)

	fmt.Printf("Building %s...\n", buildPath)

	return sh.RunWith(map[string]string{"GOOS": "linux"}, "go", "build", "-ldflags", "-s -X main.buildTime="+buildTime, "-o", buildPath, "./cmd/server/")

}

// build the binary locally
func LocalBuild() error {

	mg.Deps(Prepare)

	buildPath := path.Join(buildDir, buildBinaryLocal)

	fmt.Printf("Building locally %s...\n", buildPath)

	return sh.RunWith(map[string]string{"GOOS": "windows"}, "go", "build", "-ldflags", "-s -X main.buildTime="+buildTime, "-o", buildPath, "./cmd/server/")

}

// building and running locally
func Debugrun() error {

	mg.Deps(Prepare)

	mg.Deps(LocalBuild)

	buildPath := path.Join(buildDir, buildBinaryLocal)

	fmt.Printf("Running locally %s...\n", buildPath)

	err := sh.Run(buildPath, "-env", "staging", "-securekey", secureKey)

	return err
}

func LoadEnvironment() {
	fmt.Println("Loading environment variables...")
	buildDir = os.Getenv("BUILD_DIR")
	buildBinary = os.Getenv("BUILD_BINARY")
	buildBinaryLocal = os.Getenv("BUILD_BINARY_LOCAL")
	deployTarget = os.Getenv("DEPLOY_TARGET")
	deployDir = os.Getenv("DEPLOY_DIR")
	sshConnect = os.Getenv("SSH_CONNECT")
	sshPort = os.Getenv("SSH_PORT")
	sshUser = os.Getenv("SSH_USER")
	sshKeyfile = os.Getenv("SSH_KEYFILE")
	statsApiUrl = os.Getenv("STATS_API_URL")
	secureKey = os.Getenv("SECUREKEY")
	buildTime = time.Now().Format("2006-01-02_15:04:05")
}

// Prepare directory for builds
func Prepare() {
	mg.Deps(LoadEnvironment)
	fmt.Printf("Prepare %s directory...\n", buildDir)
	if err := os.Mkdir(buildDir, os.ModePerm); err != nil {
		log.Printf("Creating %s directory didn't work: ", buildDir, err.Error())
	}
}

// Clean up after yourself
func Clean() {
	mg.Deps(LoadEnvironment)
	fmt.Println("Cleaning...")
	os.RemoveAll(buildDir)
}
