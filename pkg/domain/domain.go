package domain

import "time"

type RideID uint64
type Fare float32

// Position is for each position point
type Position struct {
	RideID    RideID
	Lat       float64
	Long      float64
	Timestamp int64
}

// SegmentDelta is the result of calculating the difference between 2 points
type SegmentDelta struct {
	RideID   RideID
	Dirty    bool
	Distance float32 //Km
	Duration float32 // Hours
	Date     time.Time
	Velocity float32
}

// SegmentFare is the fare of a given segment
type SegmentFare struct {
	ID   RideID
	Fare Fare
}

// OutputFare is the fare of the whole ride with all the SegmentFare aggregated
type OutputFare struct {
	ID   RideID
	Fare Fare
}
