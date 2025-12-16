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
	personaldata.Personal
}

func (t *Training) Parse(datastring string) error {
	if datastring == "" {
		return fmt.Errorf("пустая строка")
	}

	parts := strings.Split(datastring, ",")
	if len(parts) != 3 {
		return fmt.Errorf("неверный формат: ожидалось 3 части, получено %d", len(parts))
	}

	stepsStr := strings.TrimSpace(parts[0])
	trainingType := strings.TrimSpace(parts[1])
	durationStr := strings.TrimSpace(parts[2])

	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return fmt.Errorf("не удалось преобразовать количество шагов: %w", err)
	}

	if steps <= 0 {
		return fmt.Errorf("количество шагов должно быть положительным числом, получено: %d", steps)
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return fmt.Errorf("ошибка парсинга продолжительности: %w", err)
	}

	if duration <= 0 {
		return fmt.Errorf("продолжительность должна быть положительной, получено: %v", duration)
	}

	t.Steps = steps
	t.TrainingType = trainingType
	t.Duration = duration
	return nil
}

func (t Training) ActionInfo() (string, error) {
	if t.Height <= 0 {
		return "", fmt.Errorf("рост должен быть положительным числом")
	}

	var calories float64
	var err error

	switch t.TrainingType {
	case "Ходьба":
		calories, err = spentenergy.WalkingSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
	case "Бег":
		calories, err = spentenergy.RunningSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}

	if err != nil {
		return "", err
	}

	distance := spentenergy.Distance(t.Steps, t.Height)
	speed := spentenergy.MeanSpeed(t.Steps, t.Height, t.Duration)

	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		t.TrainingType, t.Duration.Hours(), distance, speed, calories), nil
}
