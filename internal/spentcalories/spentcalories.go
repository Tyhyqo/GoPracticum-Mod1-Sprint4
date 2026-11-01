package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	errInvalidFormat   = errors.New("invalid data format")
	errInvalidSteps    = errors.New("invalid steps count")
	errInvalidDuration = errors.New("invalid duration")
	errInvalidWeight   = errors.New("invalid weight")
	errInvalidHeight   = errors.New("invalid height")
	errUnknownTraining = errors.New("неизвестный тип тренировки") // текст задан ТЗ и проверяется тестом
)

// Основные константы, необходимые для расчётов.
const (
	mInKm                      = 1000 // количество метров в километре
	minInH                     = 60   // количество минут в часе
	stepLengthCoefficient      = 0.45 // коэффициент для расчёта длины шага на основе роста
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчёта калорий при ходьбе
)

// Виды тренировок.
const (
	trainingRunning = "Бег"
	trainingWalking = "Ходьба"
)

// parseTraining разбирает строку вида "3456,Ходьба,3h00m"
// на количество шагов, вид тренировки и продолжительность.
func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("%w: expected 3 comma-separated values, got %d",
			errInvalidFormat, len(parts))
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("%w: %w", errInvalidSteps, err)
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("%w: %d, must be positive", errInvalidSteps, steps)
	}

	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("%w: %w", errInvalidDuration, err)
	}
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("%w: %v, must be positive", errInvalidDuration, duration)
	}

	return steps, parts[1], duration, nil
}

// distance возвращает дистанцию в километрах. Длина шага зависит от роста.
func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	return float64(steps) * stepLength / mInKm
}

// meanSpeed возвращает среднюю скорость в км/ч.
// Для неположительной продолжительности возвращает 0.
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	return distance(steps, height) / duration.Hours()
}

// TrainingInfo возвращает отчёт о тренировке: вид, длительность,
// дистанцию, среднюю скорость и сожжённые калории.
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, trainingType, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var calories float64
	switch trainingType {
	case trainingRunning:
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	case trainingWalking:
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", fmt.Errorf("%w: %q", errUnknownTraining, trainingType)
	}
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		trainingType,
		duration.Hours(),
		distance(steps, height),
		meanSpeed(steps, height, duration),
		calories,
	), nil
}

// validateInput проверяет, что все параметры для расчёта калорий положительные.
func validateInput(steps int, weight, height float64, duration time.Duration) error {
	switch {
	case steps <= 0:
		return fmt.Errorf("%w: %d, must be positive", errInvalidSteps, steps)
	case weight <= 0:
		return fmt.Errorf("%w: %.2f, must be positive", errInvalidWeight, weight)
	case height <= 0:
		return fmt.Errorf("%w: %.2f, must be positive", errInvalidHeight, height)
	case duration <= 0:
		return fmt.Errorf("%w: %v, must be positive", errInvalidDuration, duration)
	}
	return nil
}

// RunningSpentCalories возвращает количество калорий, потраченных при беге.
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if err := validateInput(steps, weight, height, duration); err != nil {
		return 0, err
	}

	speed := meanSpeed(steps, height, duration)
	return weight * speed * duration.Minutes() / minInH, nil
}

// WalkingSpentCalories возвращает количество калорий, потраченных при ходьбе.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if err := validateInput(steps, weight, height, duration); err != nil {
		return 0, err
	}

	speed := meanSpeed(steps, height, duration)
	return weight * speed * duration.Minutes() / minInH * walkingCaloriesCoefficient, nil
}
