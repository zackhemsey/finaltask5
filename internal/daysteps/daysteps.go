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
	Personal personaldata.Personal
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
	if stepsStr != strings.TrimSpace(stepsStr) {
		return fmt.Errorf("недопустимые символы в количестве шагов")
	}
	if stepsStr == "" || stepsStr == "+" || stepsStr == "-" {
		return fmt.Errorf("недопустимое значение количества шагов")
	}
	hasSign := false
	for i, r := range stepsStr {
		if i == 0 && (r == '+' || r == '-') {
			hasSign = true
			continue
		}
		if r < '0' || r > '9' {
			return fmt.Errorf("недопустимые символы в количестве шагов")
		}
	}
	if hasSign && len(stepsStr) == 1 {
		return fmt.Errorf("количество шагов содержит только знак")
	}
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return fmt.Errorf("не удалось преобразовать количество шагов: %v", err)
	}
	if steps <= 0 {
		return fmt.Errorf("количество шагов должно быть положительным числом, получено: %d", steps)
	}
	durationStr := strings.TrimSpace(parts[1])
	if durationStr == "" {
		return fmt.Errorf("продолжительность не может быть пустой")
	}
	duration, err := parseDuration(durationStr)
	if err != nil {
		return fmt.Errorf("ошибка парсинга продолжительности: %v", err)
	}
	if duration <= 0 {
		return fmt.Errorf("продолжительность должна быть положительной, получено: %v", duration)
	}
	ds.Steps = steps
	ds.Duration = duration
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

func (ds DaySteps) ActionInfo() (string, error) {

	distance := spentenergy.Distance(ds.Steps, ds.Personal.Height)
	calories, err := spentenergy.WalkingSpentCalories(ds.Steps, ds.Personal.Weight, ds.Personal.Height, ds.Duration)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		ds.Steps, distance, calories), nil
}

func (ds DaySteps) Print() {
	ds.Personal.Print()
}

func (ds DaySteps) Weight() float64 {
	return ds.Personal.Weight
}

func (ds DaySteps) Height() float64 {
	return ds.Personal.Height
}
