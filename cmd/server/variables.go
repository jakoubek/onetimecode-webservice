package main

import "github.com/jakoubek/onetimecode-webservice/internal/vcs"

var (
	version     = vcs.Version()
	buildTime   string
	isDebugMode string = "false"
)
