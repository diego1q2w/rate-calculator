package api

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"rate-calculator/pkg/domain"
	"strconv"
)

//go:generate moq -out segmenter_mock_test.go . segmenter
type segmenter interface {
	Segment(position *domain.Position) error
}

type FileReader struct {
	segmenter  segmenter
	filePath   string
	lineLength int
}

func NewFileReader(segmenter segmenter, filePath string, lineLength int) *FileReader {
	return &FileReader{segmenter: segmenter, filePath: filePath, lineLength: lineLength}
}

func (f *FileReader) Process() error {
	csvFile, err := os.Open(f.filePath)
	if err != nil {
		return fmt.Errorf("unable to open file: %w", err)
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("unable to read line: %w", err)
		}
		if f.lineLength > len(line) {
			return fmt.Errorf("fields count expected to be at least %d, got %d", f.lineLength, len(line))
		}

		if err := f.process(line); err != nil {
			return err
		}
	}

	return nil
}

func (f *FileReader) process(line []string) error {
	parseErrs := make([]error, 0)
	rideId, err := strconv.ParseUint(line[0], 10, 64)
	parseErrs = append(parseErrs, err)
	lat, err := strconv.ParseFloat(line[1], 16)
	parseErrs = append(parseErrs, err)
	long, err := strconv.ParseFloat(line[2], 16)
	parseErrs = append(parseErrs, err)
	timestamp, err := strconv.ParseInt(line[3], 10, 64)
	parseErrs = append(parseErrs, err)
	if err := parseErrors(parseErrs); err != nil {
		return err
	}

	if err := f.segmenter.Segment(&domain.Position{
		RideID:    domain.RideID(rideId),
		Lat:       lat,
		Long:      long,
		Timestamp: timestamp,
	}); err != nil {
		return fmt.Errorf("unable to segment: %w", err)
	}

	return nil
}

func parseErrors(errs []error) error {
	for _, err := range errs {
		if err != nil {
			return fmt.Errorf("unable to parse field: %w", err)
		}
	}

	return nil
}
