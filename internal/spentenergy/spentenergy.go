package spentenergy

import (
	"fmt"
	"time"
)

const (
	mInKm                      = 1000.0
	stepLengthCoefficient      = 0.45
	walkingCaloriesCoefficient = 0.5
	minInH                     = 60.0
)

func Distance(steps int, height float64) float64 {
	if steps <= 0 || height <= 0 {
		return 0
	}

	stepLength := height * stepLengthCoefficient
	return (float64(steps) * stepLength) / mInKm
}

func MeanSpeed(steps int, height float64, duration time.Duration) float64 {
	if steps <= 0 || height <= 0 || duration <= 0 {
		return 0
	}

	distance := Distance(steps, height)
	hours := duration.Hours()

	if hours == 0 {
		return 0
	}

	return distance / hours
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("некорректные входные параметры")
	}

	meanSpeed := MeanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()

	calories := (weight * meanSpeed * durationInMinutes) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("некорректные входные параметры")
	}

	meanSpeed := MeanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()

	calories := (weight * meanSpeed * durationInMinutes) / minInH

	return calories * walkingCaloriesCoefficient, nil
}
