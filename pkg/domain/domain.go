package domain

import "time"

type RideID uint64

type Position struct {
	RideID    RideID
	Lat       float64
	Long      float64
	Timestamp int64
}

type SegmentDelta struct {
	RideID   RideID
	Dirty    bool
	Distance float32 //Km
	Duration float32 // Hours
	Date     time.Time
	Velocity float32
}
