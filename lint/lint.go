package lint

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/TwinProduction/go-color"
)

var lintErrors []string

func printError(message string) {
	fmt.Printf(color.Ize(color.Red, "%s\n"), message)
}

func configInVariables(variables []Variable, name string) bool {
	for _, variable := range variables {
		if variable.Name == name {
			return true
		}
	}
	return false
}

func ShowAll(logger *log.Logger, landingZonePath string) {
	variables, err := loadLandingZoneVariables(logger, landingZonePath)
	if err != nil {
		printError(fmt.Sprintf("%s", err))
		os.Exit(-1)
	}
	var options []string
	options = make([]string, 0)
	for _, variable := range variables {
		options = append(options, variable.Name)
	}

	sort.Strings(options)
	for _, option := range options {
		fmt.Println(option)
	}

}

func CafLint(logger *log.Logger, landingZonePath string, configPath string) bool {
	lintErrors = make([]string, 0)
	variables, err := loadLandingZoneVariables(logger, landingZonePath)
	if err != nil {
		printError(fmt.Sprintf("%s", err))
		os.Exit(-1)
	}

	var configs map[string]Config
	configs, err = LoadConfigs(logger, configPath)
	if err != nil {
		printError(fmt.Sprintf("Parse Error: Invalid Configuration %s\n", err))
		os.Exit(1)
	}

	for name, config := range configs {
		found := configInVariables(variables, name)
		if !found {
			lintErrors = append(lintErrors, fmt.Sprintf("%s is not a valid configuration. %s (line: %d col: %d)", name, config.File, config.Line, config.Column))
		}
	}

	errorLength := len(lintErrors)
	if errorLength > 0 {
		printError(fmt.Sprintf("Lint failed: %d error%s found", errorLength, pluralString(errorLength)))
		for _, error := range lintErrors {
			printError(error)
		}
		os.Exit(1)
	}

	return true
}

func pluralString(size int) string {
	if size > 1 {
		return "s"
	}
	return ""
}
