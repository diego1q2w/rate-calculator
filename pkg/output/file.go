package output

import (
	"bufio"
	"fmt"
	"os"
	"rate-calculator/pkg/domain"
)

type FileOutput struct {
	filePath string
}

func NewFileOutput(filePath string) *FileOutput {
	return &FileOutput{filePath: filePath}
}

func (o *FileOutput) Output(output []*domain.OutputFare) error {
	file, err := os.Create(o.filePath)
	if err != nil {
		return fmt.Errorf("unable to open file: %w", err)
	}
	defer file.Close()
	w := bufio.NewWriter(file)
	defer w.Flush()

	for _, o := range output {
		if _, err := w.WriteString(fmt.Sprintf("%d, %.4f\n", o.ID, o.Fare)); err != nil {
			return fmt.Errorf("unable to write into the file: %w", err)
		}
	}

	return nil
}
