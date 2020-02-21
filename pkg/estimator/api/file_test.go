package api

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"rate-calculator/pkg/estimator/domain"
	"testing"
)

func TestFileReader(t *testing.T) {
	testCases := map[string]struct {
		numberOfFields int
		filePath       string
		segmenterErr   error
		expectedPoints []*domain.Position
		expectedErr    error
	}{
		"if wrong path error expected": {
			filePath:       "./random-paht",
			numberOfFields: 4,
			expectedErr:    errors.New("unable to open file: open ./random-paht: no such file or directory"),
			expectedPoints: []*domain.Position{},
		},
		"if wrong number of expected field error expected": {
			filePath:       "./test.csv",
			numberOfFields: 5,
			expectedErr:    errors.New("fields count expected to be at least 5, got 4"),
			expectedPoints: []*domain.Position{},
		},
		"if error while segmenting error expected": {
			filePath:       "./test.csv",
			numberOfFields: 4,
			segmenterErr:   errors.New("test"),
			expectedErr:    errors.New("unable to segment: test"),
			expectedPoints: []*domain.Position{},
		},
		"it should get the correct pounts": {
			filePath:       "./test.csv",
			numberOfFields: 4,
			expectedPoints: []*domain.Position{
				{RideID: 8, Lat: 38.012436, Long: 23.821972, Timestamp: 1405589186},
				{RideID: 8, Lat: 38.013108, Long: 23.821800, Timestamp: 1405589186},
				{RideID: 9, Lat: 37.953066, Long: 23.735606, Timestamp: 1405587697},
				{RideID: 9, Lat: 37.953009, Long: 23.735593, Timestamp: 1405587707},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			positions := make([]*domain.Position, 0)
			segmenter := &segmenterMock{SegmentFunc: func(position *domain.Position) error {
				if tc.segmenterErr == nil {
					positions = append(positions, position)
				}

				return tc.segmenterErr
			}}
			fileReader := NewFileReader(segmenter, tc.filePath, tc.numberOfFields)
			err := fileReader.Process()
			if fmt.Sprintf("%s", err) != fmt.Sprintf("%s", tc.expectedErr) {
				t.Errorf("expected error: %s, got: %s", tc.expectedErr, err)
			}

			assert.Equal(t, tc.expectedPoints, positions)
		})
	}
}
