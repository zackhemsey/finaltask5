package spentenergy

import (
	"fmt"
	"time"
)

const (
	mInKm                      = 1000
	stepLengthCoefficient      = 0.45
	walkingCaloriesCoefficient = 0.5
)

func Distance(steps int, height float64) float64 {
	if steps <= 0 || height <= 0 {
		return 0
	}
	return float64(steps) * height * stepLengthCoefficient / mInKm
}

func MeanSpeed(steps int, height float64, duration time.Duration) float64 {
	if steps <= 0 || height <= 0 || duration <= 0 {
		return 0
	}
	distance := Distance(steps, height)
	hours := duration.Hours()
	if hours <= 0 {
		return 0
	}
	return distance / hours
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть положительным числом")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть положительным числом")
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть положительным числом")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность должна быть положительной")
	}
	return Distance(steps, height) * weight, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть положительным числом")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть положительным числом")
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть положительным числом")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность должна быть положительной")
	}
	return Distance(steps, height) * weight * walkingCaloriesCoefficient, nil
}
