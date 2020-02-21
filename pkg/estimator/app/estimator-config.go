package app

import (
	"fmt"
	"rate-calculator/pkg/estimator/domain"
	"time"
)

type dateRule struct {
	start  time.Time
	finish time.Time
	check  time.Time
}
type TimeRule struct {
	Start  string
	Finish string
	Fare   domain.Fare
}

type SpeedRule struct {
	Limit float32
	Fare  domain.Fare
}

func GetEstimatorConfig(tRuleDay, tRuleNight TimeRule, speedRule SpeedRule) ([]RateConfig, error) {
	timeLayout := "15:04"
	parseErrs := make([]error, 0)
	dStart, err := time.Parse(timeLayout, tRuleDay.Start)
	parseErrs = append(parseErrs, err)
	dFinish, err := time.Parse(timeLayout, tRuleDay.Finish)
	parseErrs = append(parseErrs, err)
	nStart, err := time.Parse(timeLayout, tRuleNight.Start)
	parseErrs = append(parseErrs, err)
	nFinish, err := time.Parse(timeLayout, tRuleNight.Finish)
	parseErrs = append(parseErrs, err)
	if err := checkErrors(parseErrs); err != nil {
		return nil, fmt.Errorf("unable to parse dates: %w", err)
	}

	return []RateConfig{
		{
			Rule: func(delta *domain.SegmentDelta) (b bool, m multiplier) { // > 10 Km/H (05:00 - 00:00]
				dRule := &dateRule{
					start:  dStart,
					finish: dFinish,
					check:  time.Date(0, 1, 1, delta.Date.Hour(), delta.Date.Minute(), delta.Date.Second(), delta.Date.Nanosecond(), delta.Date.Location()),
				}
				return delta.Velocity > speedRule.Limit && inTimeSpan(dRule), delta.Distance
			}, Fare: tRuleDay.Fare},
		{
			Rule: func(delta *domain.SegmentDelta) (b bool, m multiplier) { // > 10 Km/H (00:00 - 05:00]
				dRule := &dateRule{
					start:  nStart,
					finish: nFinish,
					check:  time.Date(0, 1, 1, delta.Date.Hour(), delta.Date.Minute(), delta.Date.Second(), delta.Date.Nanosecond(), delta.Date.Location()),
				}
				return delta.Velocity > speedRule.Limit && inTimeSpan(dRule), delta.Distance
			}, Fare: tRuleNight.Fare},
		{
			Rule: func(delta *domain.SegmentDelta) (b bool, m multiplier) { // <= 10 Km/H
				return delta.Velocity <= speedRule.Limit, delta.Duration
			}, Fare: speedRule.Fare},
	}, nil
}

func checkErrors(errs []error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}

// Non inclusive start
func inTimeSpan(timeRule *dateRule) bool {
	if timeRule.start.Before(timeRule.finish) {
		return timeRule.check.After(timeRule.start) && !timeRule.check.After(timeRule.finish)
	}
	if timeRule.start.Equal(timeRule.finish) {
		return timeRule.check.Equal(timeRule.start)
	}
	return timeRule.start.Before(timeRule.check) || !timeRule.finish.Before(timeRule.check)
}
