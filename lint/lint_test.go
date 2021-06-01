package lint

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigInVariables(t *testing.T) {
	//arrange

	var configInVariablesTest = []struct {
		testID         int
		testName       string
		variables      []Variable
		configName     string
		expectedResult bool
	}{
		{
			testID:         0,
			testName:       "Test config exists in lz variables.",
			configName:     "test1",
			expectedResult: true,
			variables:      []Variable{{Name: "test1"}, {Name: "test2"}, {Name: "test3"}},
		},
		{
			testID:         1,
			testName:       "Test config does not exist in lz variables.",
			configName:     "doesnotexist",
			expectedResult: false,
			variables:      []Variable{{Name: "test1"}, {Name: "test2"}, {Name: "test3"}},
		},
		{
			testID:         3,
			testName:       "Empty variables list should return false.",
			configName:     "doesnotexist",
			expectedResult: false,
			variables:      []Variable{},
		},
		{
			testID:         4,
			testName:       "Empty configName list should return false.",
			configName:     "",
			expectedResult: false,
			variables:      []Variable{},
		},
	}

	for _, test := range configInVariablesTest {
		testDisplay := fmt.Sprintf("%d - %s", test.testID, test.testName)
		t.Run(testDisplay, func(t *testing.T) {
			//act
			result := configInVariables(test.variables, test.configName)

			//assert
			assert.Equal(t, result, test.expectedResult)

		})
	}
}

func TestPluralizeString(t *testing.T) {
	//arrange

	var configInVariablesTest = []struct {
		testID         int
		testName       string
		size           int
		expectedResult string
	}{
		{
			testID:         0,
			testName:       "Test s is returned if size is more than 1.",
			size:           2,
			expectedResult: "s",
		},
		{
			testID:         1,
			testName:       "Test empty string is returned if size is 1.",
			size:           1,
			expectedResult: "",
		},
		{
			testID:         2,
			testName:       "Test empty string is returned if size is zero.",
			size:           0,
			expectedResult: "",
		},
		{
			testID:         2,
			testName:       "Test empty string is returned if size is less than zero.",
			size:           -10,
			expectedResult: "",
		},
	}

	for _, test := range configInVariablesTest {
		testDisplay := fmt.Sprintf("%d - %s", test.testID, test.testName)
		t.Run(testDisplay, func(t *testing.T) {
			//act
			result := pluralString(test.size)

			//assert
			assert.Equal(t, result, test.expectedResult)

		})
	}
}

func TestPrintError(t *testing.T) {
	//arrange
	buf := new(bytes.Buffer)
	logger := log.New(buf, "", 0)

	message := "hello world"
	expectedOut := `[31mhello world
[0m
` // this is indented weird, but needed to validate the output of printError. The non printable characters is how colorize sets it to red.

	//act
	printError(logger, message)

	//assert
	assert.Equal(t, expectedOut, buf.String())
}

func TestShowAllReturnsDefinedResources(t *testing.T) {
	//arrange
	exiter := NewExiter(func(i int) {})
	lzPath := "../testharness/tf"
	expected := `global_settings
landingzone
resource_groups
subscriptions
tfstate_container_name
tfstate_key
tfstate_resource_group_name
`
	buf := new(bytes.Buffer)
	logger := log.New(buf, "", 0)

	ShowAll(logger, exiter, lzPath)

	assert.Equal(t, expected, buf.String())
}

func TestShowAllReturnsCorrectNumberOfItems(t *testing.T) {
	//arrange
	exiter := NewExiter(func(i int) {})
	lzPath := "../testharness/tf"
	buf := new(bytes.Buffer)
	logger := log.New(buf, "", 0)

	ShowAll(logger, exiter, lzPath)
	result := buf.String()

	data := strings.Split(result, "\n")
	assert.Equal(t, 8, len(data))
}

func TestShowAllWithInvalidPathDisplaysCorrectErrorMessage(t *testing.T) {
	//arrange
	exiter := NewExiter(func(i int) {})
	lzPath := "fake_path"
	expected := `[31mcannot read lz path
message:open fake_path: no such file or directory
[0m
`
	buf := new(bytes.Buffer)
	logger := log.New(buf, "", 0)

	//act
	ShowAll(logger, exiter, lzPath)

	//assert
	assert.Equal(t, expected, buf.String())
}

func TestShowAllWithInvalidPathReturnsCorrectExitCode(t *testing.T) {
	//arrange
	exiter := NewExiter(func(i int) {})
	lzPath := "fake_path"

	buf := new(bytes.Buffer)
	logger := log.New(buf, "", 0)

	//act
	ShowAll(logger, exiter, lzPath)

	//assert
	assert.Equal(t, FILE_OR_FOLDER_NOT_FOUND, exiter.statusCode)
}

func TestCafLintWithValidConfig(t *testing.T) {
	//arrange
	exiter := NewExiter(func(i int) {})
	lzPath := "../testharness/tf"
	configPath := "../testharness/config/valid"

	buf := new(bytes.Buffer)
	logger := log.New(buf, "", 0)

	//act
	CafLint(logger, exiter, lzPath, configPath)

	//assert
	assert.Equal(t, "", buf.String())
	assert.Equal(t, SUCCESS, exiter.statusCode)
}

func TestCafLintWithInValidConfigValue(t *testing.T) {
	//arrange
	exiter := NewExiter(func(i int) {})
	lzPath := "../testharness/tf"
	configPath := "../testharness/config/invalid/invalidConfigName"
	expected := `[31mLint failed: 1 error found
[0m
[31mlandingzoneFake is not a valid configuration. ../testharness/config/invalid/invalidConfigName/one.tfvars (line: 1 col: 1)
[0m
`
	buf := new(bytes.Buffer)
	logger := log.New(buf, "", 0)

	//act
	CafLint(logger, exiter, lzPath, configPath)

	//assert
	assert.Equal(t, expected, buf.String())
	assert.Equal(t, LINT_ERROR, exiter.statusCode)
}

func TestCafLintWithInValidConfigPath(t *testing.T) {
	//arrange
	exiter := NewExiter(func(i int) {})
	lzPath := "../testharness/tf"
	configPath := "../testharness/wrongpath"
	expected := `[31mParse Error: Invalid Configuration cannot read config path
message:open ../testharness/wrongpath: no such file or directory

[0m
`
	buf := new(bytes.Buffer)
	logger := log.New(buf, "", 0)

	//act
	CafLint(logger, exiter, lzPath, configPath)

	//assert
	assert.Equal(t, expected, buf.String())
	assert.Equal(t, FILE_OR_FOLDER_NOT_FOUND, exiter.statusCode)
}

func TestCafLintWithInValidTFVarsFormat(t *testing.T) {
	//arrange
	exiter := NewExiter(func(i int) {})
	lzPath := "../testharness/tf"
	configPath := "../testharness/config/invalid/invalidTfVarsFile"
	expected := `[31mParse Error: Invalid Configuration cannot read config path
message:[../testharness/config/invalid/invalidTfVarsFile/one.tfvars:2,37-3,1: Invalid multi-line string; Quoted strings may not be split over multiple lines. To produce a multi-line string, either use the \n escape to represent a newline character or use the "heredoc" multi-line template syntax. ../testharness/config/invalid/invalidTfVarsFile/one.tfvars:2,34-36: Missing attribute separator; Expected a newline or comma to mark the beginning of the next attribute.]

[0m
`
	buf := new(bytes.Buffer)
	logger := log.New(buf, "", 0)

	//act
	CafLint(logger, exiter, lzPath, configPath)

	//assert
	assert.Equal(t, expected, buf.String())
	assert.Equal(t, INVALID_TFVARS_SYNTAX, exiter.statusCode)
}

func TestCafLintConfigFolderWithNoTfVars(t *testing.T) {
	//arrange
	exiter := NewExiter(func(i int) {})
	lzPath := "../testharness/tf"
	configPath := "../testharness/config/invalid/NoTfVarFiles"
	expected := `[31mParse Error: Invalid Configuration No .tfvars found in path ../testharness/config/invalid/NoTfVarFiles

[0m
`
	buf := new(bytes.Buffer)
	logger := log.New(buf, "", 0)

	//act
	CafLint(logger, exiter, lzPath, configPath)

	//assert
	assert.Equal(t, expected, buf.String())
	assert.Equal(t, NO_TFVARS_FOUND, exiter.statusCode)
}

func TestCafLintWithInvalidLandingZonePath(t *testing.T) {
	//arrange
	exiter := NewExiter(func(i int) {})
	lzPath := "../testharness/wrongpath"
	configPath := "../testharness/config/valid"
	expected := `[31mcannot read lz path
message:open ../testharness/wrongpath: no such file or directory
[0m
`
	buf := new(bytes.Buffer)
	logger := log.New(buf, "", 0)

	//act
	CafLint(logger, exiter, lzPath, configPath)

	//assert
	assert.Equal(t, expected, buf.String())
	assert.Equal(t, FILE_OR_FOLDER_NOT_FOUND, exiter.statusCode)
}

func TestCafLintLZFolderWithNoTFFiles(t *testing.T) {
	//arrange
	exiter := NewExiter(func(i int) {})
	lzPath := "../testharness/config/valid"
	configPath := "../testharness/config/valid"
	expected := `[31mlanding zone error: no variables.*.tf files found in path ../testharness/config/valid
[0m
`
	buf := new(bytes.Buffer)
	logger := log.New(buf, "", 0)

	//act
	CafLint(logger, exiter, lzPath, configPath)

	//assert
	assert.Equal(t, expected, buf.String())
	assert.Equal(t, NO_VARIABLE_FILES_FOUND, exiter.statusCode)
}

func TestCafLintLZFolderWithInvalidTFFiles(t *testing.T) {
	//arrange
	exiter := NewExiter(func(i int) {})
	lzPath := "../testharness/tf/invalid"
	configPath := "../testharness/config/valid"
	expected := `[31mparse variable file: [../testharness/tf/invalid/variables.one.tf:11,1-1: Invalid expression; Expected the start of an expression, but found an invalid expression token.]
[0m
`
	buf := new(bytes.Buffer)
	logger := log.New(buf, "", 0)

	//act
	CafLint(logger, exiter, lzPath, configPath)

	//assert
	assert.Equal(t, expected, buf.String())
	assert.Equal(t, INVALID_VARIABLE_FILE_SYNTAX, exiter.statusCode)
}
