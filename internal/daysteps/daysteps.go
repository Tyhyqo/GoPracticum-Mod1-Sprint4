package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

var (
	errInvalidFormat   = errors.New("invalid data format")
	errInvalidSteps    = errors.New("invalid steps count")
	errInvalidDuration = errors.New("invalid duration")
)

// Основные константы, необходимые для расчётов.
const (
	stepLength = 0.65 // длина одного шага в метрах
	mInKm      = 1000 // количество метров в километре
)

// parsePackage разбирает строку вида "678,0h50m"
// на количество шагов и продолжительность прогулки.
func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("%w: expected 2 comma-separated values, got %d",
			errInvalidFormat, len(parts))
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("%w: %w", errInvalidSteps, err)
	}
	if steps <= 0 {
		return 0, 0, fmt.Errorf("%w: %d, must be positive", errInvalidSteps, steps)
	}

	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("%w: %w", errInvalidDuration, err)
	}
	if duration <= 0 {
		return 0, 0, fmt.Errorf("%w: %v, must be positive", errInvalidDuration, duration)
	}

	return steps, duration, nil
}

// DayActionInfo возвращает отчёт о прогулке: количество шагов,
// дистанцию в километрах и сожжённые калории.
// При ошибке выводит её на экран и возвращает пустую строку.
func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println("daysteps:", err)
		return ""
	}

	distanceMeters := float64(steps) * stepLength
	distanceKm := distanceMeters / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		fmt.Println("daysteps:", err)
		return ""
	}

	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps,
		distanceKm,
		calories,
	)
}
