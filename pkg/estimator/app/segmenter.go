package app

import (
	"errors"
	"fmt"
	"github.com/umahmood/haversine"
	"rate-calculator/pkg/estimator/domain"
	"time"
)

type Segmenter struct {
	workerCh     chan segment
	filter       segmentFilter
	rideSegment  rideSegment
	distanceCalc distanceCalc
}

func NewSegmenter(filter segmentFilter, distanceCalc distanceCalc, workers int) *Segmenter {
	s := &Segmenter{
		filter:       filter,
		distanceCalc: distanceCalc,
		rideSegment:  make(map[domain.RideID]segment),
	}
	if workers > 0 {
		s.spinWorkers(workers)
	}

	return s
}

type distanceCalc = func(p, q haversine.Coord) (mi, km float64)

//go:generate moq -out segment_filter_mock_test.go . segmentFilter
type segmentFilter interface {
	Filter(delta *domain.SegmentDelta) error
}

type rideSegment map[domain.RideID]segment

type segment struct {
	id domain.RideID
	p1 *domain.Position
	p2 *domain.Position
}

func (s *segment) isReady() bool {
	return s.p1 != nil && s.p2 != nil
}

func (s *segment) pushElement(position *domain.Position) error {
	if position == nil {
		return errors.New("can not have a nil position")
	}
	if s.p1 == nil {
		s.p1 = position
	} else if s.p2 == nil {
		s.p2 = position
	} else {
		s.p1 = s.p2
		s.p2 = position
	}

	return nil
}

func (s *segment) calculate(distanceCalc distanceCalc) (*domain.SegmentDelta, error) {
	sDelta := &domain.SegmentDelta{
		RideID: s.id,
		Dirty:  false,
		Date:   time.Unix(s.p1.Timestamp, 0).UTC(),
	}

	p1 := haversine.Coord{Lat: s.p1.Lat, Lon: s.p1.Long}
	p2 := haversine.Coord{Lat: s.p2.Lat, Lon: s.p2.Long}
	_, km := distanceCalc(p1, p2)
	sDelta.Distance = float32(km)

	sDelta.Duration = float32(s.p2.Timestamp-s.p1.Timestamp) / float32(3600)

	sDelta.Velocity = sDelta.Distance / sDelta.Duration
	return sDelta, nil
}

func (s *Segmenter) Segment(position *domain.Position) error {
	rSegment, ok := s.rideSegment[position.RideID]
	if !ok {
		rSegment = segment{id: position.RideID}
	}
	if err := rSegment.pushElement(position); err != nil {
		return err
	}
	s.rideSegment[position.RideID] = rSegment

	if rSegment.isReady() {
		if s.workerCh != nil {
			s.workerCh <- rSegment
		} else {
			s.calculate(rSegment)
		}
	}
	return nil
}

func (s *Segmenter) spinWorkers(numberOfWorkers int) {
	ch := make(chan segment)
	s.workerCh = ch
	for i := 0; i < numberOfWorkers; i++ {
		go s.worker()
	}
}

func (s *Segmenter) Close() {
	close(s.workerCh)
}

func (s *Segmenter) worker() {
	for segment := range s.workerCh {
		s.calculate(segment)
	}
}

func (s *Segmenter) calculate(segment segment) {
	sDelta, err := segment.calculate(s.distanceCalc)
	if err != nil {
		fmt.Printf("error calculating distance: %s", err)
		return
	}

	if err := s.filter.Filter(sDelta); err != nil {
		fmt.Printf("error applying filer: %s", err)
	}
}
