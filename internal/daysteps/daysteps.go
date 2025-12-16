package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

type DaySteps struct {
	Steps    int
	Duration time.Duration
	personaldata.Personal
}

func (ds *DaySteps) Parse(datastring string) error {
	if datastring == "" {
		return fmt.Errorf("пустая строка")
	}

	parts := strings.Split(datastring, ",")
	if len(parts) != 2 {
		return fmt.Errorf("неверный формат: ожидалось 2 части, получено %d", len(parts))
	}

	stepsStr := parts[0]
	durationStr := strings.TrimSpace(parts[1])

	if stepsStr != strings.TrimSpace(stepsStr) {
		return fmt.Errorf("недопустимые символы в количестве шагов")
	}

	stepsStr = strings.TrimSpace(stepsStr)

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

	ds.Steps = steps
	ds.Duration = duration
	return nil
}

func (ds DaySteps) ActionInfo() (string, error) {
	distance := spentenergy.Distance(ds.Steps, ds.Height)
	calories, err := spentenergy.WalkingSpentCalories(ds.Steps, ds.Weight, ds.Height, ds.Duration)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		ds.Steps, distance, calories), nil
}

func (ds DaySteps) Print() {
	fmt.Printf("Шаги: %d, Продолжительность: %v\n", ds.Steps, ds.Duration)
}
