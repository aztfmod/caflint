package lint

import (
	"bytes"
	"fmt"
	"io"
	"os"
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
	message := "hello world"
	expectedOut := `[31mhello world
[0m` // this is idented weird, but needed to validate the output of printError. The non printable characters is how colorize sets it to red.

	//redirect stdout
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	//act
	printError(message)

	//grab stdout and revert redirect
	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()
	// back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC

	//assert
	assert.Equal(t, expectedOut, out)
}
