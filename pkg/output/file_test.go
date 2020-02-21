package output

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"rate-calculator/pkg/domain"
	"testing"
)

func TestFileOutput(t *testing.T) {
	outputFares := []*domain.OutputFare{
		{ID: 1, Fare: 1},
		{ID: 2, Fare: 2},
		{ID: 3, Fare: 3},
	}
	expectedOutput := `1, 1.0000
2, 2.0000
3, 3.0000
4, 4.0000
`
	filePath := "./test.txt"

	output := NewFileOutput(filePath)
	err := output.Output(outputFares)
	assert.NoError(t, err)

	outputFares = append(outputFares, &domain.OutputFare{ID: 4, Fare: 4})
	err = output.Output(outputFares)
	assert.NoError(t, err)

	content := readFileContent(t, filePath)
	assert.Equal(t, expectedOutput, content)
}

func TestOpenFileError_Error(t *testing.T) {
	outputFares := []*domain.OutputFare{
		{ID: 1, Fare: 1},
		{ID: 2, Fare: 2},
		{ID: 3, Fare: 3},
	}

	filePath := "./random-wrong/test.txt"

	output := NewFileOutput(filePath)
	err := output.Output(outputFares)
	assert.EqualError(t, err, "open ./random-wrong/test.txt: no such file or directory")
}

func readFileContent(t *testing.T, path string) string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("unexpected error while reading: %s", err)
	}

	return string(data)
}
