package lint

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type Variable struct {
	Name string
}

type Config struct {
	File   string
	Expr   hcl.Expression
	Line   int
	Column int
}

type CafConfig struct {
	LandingZonePath string `hcl:"landingZonePath,optional"`
	ConfigPath      string `hcl:"configPath,optional"`
}

var cafConfigSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{Type: "config"},
	},
}
var mapSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{Type: "variable", LabelNames: []string{"id"}},
	},
}

func loadLandingZoneVariables(logger *log.Logger, landingZonePath string) ([]Variable, *LintError) {
	var variables []Variable
	files, err := ioutil.ReadDir(landingZonePath)
	if err != nil {
		return nil, NewLintError(FILE_OR_FOLDER_NOT_FOUND, "cannot read lz path\nmessage:%s", err)
	}

	foundVariablesConfig := false
	for _, f := range files {
		name := f.Name()
		filePath := path.Join(landingZonePath, name)
		extension := filepath.Ext(filePath)

		if strings.HasPrefix(name, "variables.") && extension == ".tf" {
			foundVariablesConfig = true
			vars, err := parseVariableFile(logger, filePath)
			if err != nil {
				return nil, NewLintError(INVALID_VARIABLE_FILE_SYNTAX, "%v", err)
			}
			variables = append(variables, vars...)
		}
	}
	if !foundVariablesConfig {
		return nil, NewLintError(NO_VARIABLE_FILES_FOUND, "landing zone error: no variables.*.tf files found in path %s", landingZonePath)
	}
	return variables, nil
}

func parseVariableFile(logger *log.Logger, variablesFilePath string) ([]Variable, error) {
	variablefileBytes, err := ioutil.ReadFile(variablesFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %s", err)
	}

	file, diags := hclsyntax.ParseConfig(variablefileBytes, variablesFilePath, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, fmt.Errorf("parse variable file: %v", diags.Errs())
	}

	content, diags := file.Body.Content(mapSchema)
	if diags.HasErrors() {
		return nil, fmt.Errorf("content: %v", diags.Errs())
	}

	var variables []Variable = make([]Variable, 0)
	for _, block := range content.Blocks {
		switch block.Type {
		case "variable":
			v := Variable{}
			v.Name = block.Labels[0]
			variables = append(variables, v)
		}
	}
	return variables, nil
}

func LoadConfigs(logger *log.Logger, configFilePath string) (map[string]Config, *LintError) {
	var configurations map[string]Config
	configurations = make(map[string]Config)

	files, err := ioutil.ReadDir(configFilePath)
	if err != nil {
		return nil, NewLintError(FILE_OR_FOLDER_NOT_FOUND, "cannot read config path\nmessage:%s", err)
	}

	foundVars := false
	for _, f := range files {
		name := f.Name()
		filePath := path.Join(configFilePath, name)
		extension := filepath.Ext(filePath)

		if extension == ".tfvars" {
			foundVars = true
			configurations, err = parseConfigFile(filePath, configurations)
			if err != nil {
				return nil, NewLintError(INVALID_TFVARS_SYNTAX, "cannot read config path\nmessage:%s", err)
			}
		}
	}
	if !foundVars {
		return nil, NewLintError(NO_TFVARS_FOUND, "No .tfvars found in path %s", configFilePath)
	}
	return configurations, nil
}

func parseConfigFile(filename string, configurations map[string]Config) (map[string]Config, error) {
	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %s", err)
	}

	file, diags := hclsyntax.ParseConfig(fileBytes, filename, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, fmt.Errorf("%v", diags.Errs())
	}
	attrs, hclDiags := file.Body.JustAttributes()
	if hclDiags.HasErrors() {
		return nil, fmt.Errorf("parse config: %v", hclDiags.Errs())
	}

	for name, attr := range attrs {
		configurations[name] = Config{
			File:   filename,
			Expr:   attr.Expr,
			Line:   attr.Range.Start.Line,
			Column: attr.Range.Start.Column,
		}
	}
	return configurations, nil
}

func ParseCafLintFile(filename string) (*CafConfig, error) {
	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("file Not found %s", filename)
	}

	file, diags := hclsyntax.ParseConfig(fileBytes, filename, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, fmt.Errorf("%v", diags.Errs())
	}

	content, diags := file.Body.Content(cafConfigSchema)
	if diags.HasErrors() {
		return nil, fmt.Errorf("%v", diags.Errs())
	}
	ctx := &hcl.EvalContext{}

	for _, block := range content.Blocks {
		switch block.Type {
		case "config":
			config := CafConfig{
				LandingZonePath: "",
				ConfigPath:      "",
			}
			diags := gohcl.DecodeBody(block.Body, ctx, &config)
			if diags != nil {
				return nil, fmt.Errorf("%v", diags.Errs())
			}
			return &config, nil
		}
	}
	return nil, fmt.Errorf("Config not found")
}
