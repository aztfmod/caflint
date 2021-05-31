package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aztfmod/caflint/lint"
)

var cafLintconfig string = "./.caflint.hcl"

func main() {
	logger := log.New(os.Stderr, "", 0)

	lzPath := flag.String("lz", "", "path to the landing zone")
	configPath := flag.String("var-folder", "", "path to caf configs")
	showAll := flag.Bool("show-all", false, "show all available configurations")

	flag.Parse()

	var cafConfig lint.CafConfig
	if _, err := os.Stat(cafLintconfig); err == nil {
		config, err := lint.ParseCafLintFile(cafLintconfig)
		if err != nil {
			panic(err)
		}
		cafConfig = *config
	}

	var landingZonePath string = *lzPath
	var configurationPath string = *configPath
	if *lzPath == "" {
		landingZonePath = cafConfig.LandingZonePath
	}
	if *configPath == "" {
		configurationPath = cafConfig.ConfigPath
	}

	if *showAll {
		fmt.Printf("Landing Zone: %s\n", landingZonePath)
		fmt.Printf("Available Configurations:\n-------------------------\n")
		lint.ShowAll(logger, landingZonePath)
	} else {
		lint.CafLint(logger, landingZonePath, configurationPath)
	}
}
