package lint

import (
	"fmt"
	"log"
	"sort"

	"github.com/TwinProduction/go-color"
)

var lintErrors []string

func printError(logger *log.Logger, message string) {
	logger.Printf(color.Ize(color.Red, "%s\n"), message)
}

func configInVariables(variables []Variable, name string) bool {
	for _, variable := range variables {
		if variable.Name == name {
			return true
		}
	}
	return false
}

func ShowAll(logger *log.Logger, exiter *Exiter, landingZonePath string) {
	variables, err := loadLandingZoneVariables(logger, landingZonePath)
	if err != nil {
		printError(logger, fmt.Sprintf("%s", err))
		exiter.Exit(FILE_OR_FOLDER_NOT_FOUND)
	}
	var options []string
	options = make([]string, 0)
	for _, variable := range variables {
		options = append(options, variable.Name)
	}

	sort.Strings(options)
	for _, option := range options {
		logger.Println(option)
	}

}

func CafLint(logger *log.Logger, exiter *Exiter, landingZonePath string, configPath string) bool {
	lintErrors = make([]string, 0)
	variables, err := loadLandingZoneVariables(logger, landingZonePath)
	if err != nil {
		printError(logger, fmt.Sprintf("%s", err))
		exiter.Exit(err.StatusCode)
	} else {
		configs, lintError := LoadConfigs(logger, configPath)
		if lintError != nil {
			printError(logger, fmt.Sprintf("Parse Error: Invalid Configuration %s\n", lintError))
			exiter.Exit(lintError.StatusCode)
		}

		for name, config := range configs {
			found := configInVariables(variables, name)
			if !found {
				lintErrors = append(lintErrors, fmt.Sprintf("%s is not a valid configuration. %s (line: %d col: %d)", name, config.File, config.Line, config.Column))
			}
		}

		errorLength := len(lintErrors)
		if errorLength > 0 {
			printError(logger, fmt.Sprintf("Lint failed: %d error%s found", errorLength, pluralString(errorLength)))
			for _, error := range lintErrors {
				printError(logger, error)
			}
			exiter.Exit(1)
		}
	}

	return true
}

func pluralString(size int) string {
	if size > 1 {
		return "s"
	}
	return ""
}
