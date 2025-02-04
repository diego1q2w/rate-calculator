// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package api

import (
	"rate-calculator/pkg/estimator/domain"
	"sync"
)

var (
	locksegmenterMockSegment sync.RWMutex
)

// Ensure, that segmenterMock does implement segmenter.
// If this is not the case, regenerate this file with moq.
var _ segmenter = &segmenterMock{}

// segmenterMock is a mock implementation of segmenter.
//
//     func TestSomethingThatUsessegmenter(t *testing.T) {
//
//         // make and configure a mocked segmenter
//         mockedsegmenter := &segmenterMock{
//             SegmentFunc: func(position *domain.Position) error {
// 	               panic("mock out the Segment method")
//             },
//         }
//
//         // use mockedsegmenter in code that requires segmenter
//         // and then make assertions.
//
//     }
type segmenterMock struct {
	// SegmentFunc mocks the Segment method.
	SegmentFunc func(position *domain.Position) error

	// calls tracks calls to the methods.
	calls struct {
		// Segment holds details about calls to the Segment method.
		Segment []struct {
			// Position is the position argument value.
			Position *domain.Position
		}
	}
}

// Segment calls SegmentFunc.
func (mock *segmenterMock) Segment(position *domain.Position) error {
	if mock.SegmentFunc == nil {
		panic("segmenterMock.SegmentFunc: method is nil but segmenter.Segment was just called")
	}
	callInfo := struct {
		Position *domain.Position
	}{
		Position: position,
	}
	locksegmenterMockSegment.Lock()
	mock.calls.Segment = append(mock.calls.Segment, callInfo)
	locksegmenterMockSegment.Unlock()
	return mock.SegmentFunc(position)
}

// SegmentCalls gets all the calls that were made to Segment.
// Check the length with:
//     len(mockedsegmenter.SegmentCalls())
func (mock *segmenterMock) SegmentCalls() []struct {
	Position *domain.Position
} {
	var calls []struct {
		Position *domain.Position
	}
	locksegmenterMockSegment.RLock()
	calls = mock.calls.Segment
	locksegmenterMockSegment.RUnlock()
	return calls
}
