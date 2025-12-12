package trainings

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

type Training struct {
	Steps        int
	TrainingType string
	Duration     time.Duration
	Personal     personaldata.Personal
}

func (t Training) Print() {
	t.Personal.Print()
}

func (t *Training) Parse(datastring string) error {
	if datastring == "" {
		return fmt.Errorf("пустая строка")
	}
	parts := strings.Split(datastring, ",")
	if len(parts) != 3 {
		return fmt.Errorf("неверный формат: ожидалось 3 части, получено %d", len(parts))
	}
	stepsStr := parts[0]
	if stepsStr != strings.TrimSpace(stepsStr) {
		return fmt.Errorf("недопустимые символы в количестве шагов")
	}
	stepsStr = strings.TrimPrefix(stepsStr, "+")
	if stepsStr == "" || stepsStr == "-" {
		return fmt.Errorf("недопустимое значение количества шагов")
	}
	for _, r := range stepsStr {
		if r < '0' || r > '9' {
			return fmt.Errorf("недопустимые символы в количестве шагов")
		}
	}
	steps, err := strconv.Atoi(stepsStr)
	if err != nil || steps <= 0 {
		return fmt.Errorf("количество шагов должно быть положительным числом")
	}
	trainingType := strings.TrimSpace(parts[1])
	// Accept any training type in Parse, validate in ActionInfo
	durationStr := strings.TrimSpace(parts[2])
	if durationStr == "" {
		return fmt.Errorf("продолжительность не может быть пустой")
	}
	duration, err := parseDuration(durationStr)
	if err != nil || duration <= 0 {
		return fmt.Errorf("ошибка парсинга продолжительности: %v", err)
	}
	t.Steps = steps
	t.TrainingType = trainingType
	t.Duration = duration
	return nil
}

func parseDuration(s string) (time.Duration, error) {
	var totalDuration time.Duration
	hourPart := ""
	minutePart := ""

	if hIdx := strings.Index(s, "h"); hIdx != -1 {
		if hIdx == 0 {
			return 0, fmt.Errorf("не указано значение для часов")
		}
		hourPart = s[:hIdx]
		s = s[hIdx+1:]

		hours, err := strconv.ParseFloat(hourPart, 64)
		if err != nil {
			return 0, fmt.Errorf("неверный формат часов: %v", err)
		}
		if hours < 0 {
			return 0, fmt.Errorf("часы не могут быть отрицательными")
		}
		totalDuration += time.Duration(hours * float64(time.Hour))
	}

	if mIdx := strings.Index(s, "m"); mIdx != -1 {
		if mIdx == 0 {
			return 0, fmt.Errorf("не указано значение для минут")
		}
		minutePart = s[:mIdx]

		minutes, err := strconv.ParseFloat(minutePart, 64)
		if err != nil {
			return 0, fmt.Errorf("неверный формат минут: %v", err)
		}
		if minutes < 0 {
			return 0, fmt.Errorf("минуты не могут быть отрицательными")
		}
		totalDuration += time.Duration(minutes * float64(time.Minute))
	}

	if hourPart == "" && minutePart == "" {
		return 0, fmt.Errorf("не указана единица измерения времени")
	}
	if totalDuration <= 0 {
		return 0, fmt.Errorf("продолжительность должна быть положительной")
	}
	return totalDuration, nil
}

func (t Training) ActionInfo() (string, error) {
	if t.Steps <= 0 {
		return "", fmt.Errorf("количество шагов должно быть положительным числом")
	}
	if t.Duration <= 0 {
		return "", fmt.Errorf("продолжительность должна быть положительной")
	}
	if t.Personal.Weight <= 0 {
		return "", fmt.Errorf("вес должен быть положительным числом")
	}
	if t.Personal.Height <= 0 {
		return "", fmt.Errorf("рост должен быть положительным числом")
	}
	if t.TrainingType != "Ходьба" && t.TrainingType != "Бег" {
		return "", fmt.Errorf("неизвестный тип тренировки")
	}
	distance := spentenergy.Distance(t.Steps, t.Personal.Height)
	speed := spentenergy.MeanSpeed(t.Steps, t.Personal.Height, t.Duration)
	var calories float64
	var err error
	if t.TrainingType == "Ходьба" {
		calories, err = spentenergy.WalkingSpentCalories(t.Steps, t.Personal.Weight, t.Personal.Height, t.Duration)
	} else {
		calories, err = spentenergy.RunningSpentCalories(t.Steps, t.Personal.Weight, t.Personal.Height, t.Duration)
	}
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		t.TrainingType, t.Duration.Hours(), distance, speed, calories), nil
}
