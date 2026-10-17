package spentenergy

import (
	"errors"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе.
)

// Ошибки валидации входных параметров.
var (
	ErrNegativeSteps = errors.New("spentenergy: количество шагов не может быть отрицательным")
	ErrZeroSteps     = errors.New("spentenergy: количество шагов должно быть больше нуля")
	ErrZeroWeight    = errors.New("spentenergy: вес должен быть больше нуля")
	ErrZeroHeight    = errors.New("spentenergy: рост должен быть больше нуля")
	ErrZeroDuration  = errors.New("spentenergy: продолжительность должна быть больше нуля")
)

// Distance возвращает дистанцию в километрах, пройденную пользователем
// за steps шагов при росте height (в метрах).
func Distance(steps int, height float64) float64 {
	// Рассчитываем длину шага.
	stepLength := height * stepLengthCoefficient
	// Умножаем количество шагов на длину шага и переводим метры в километры.
	return float64(steps) * stepLength / float64(mInKm)
}

// MeanSpeed возвращает среднюю скорость в км/ч.
func MeanSpeed(steps int, height float64, duration time.Duration) float64 {
	// Если шагов или продолжительности нет — скорость равна нулю.
	if steps <= 0 || duration <= 0 {
		return 0
	}
	distance := Distance(steps, height)
	// Переводим продолжительность в часы и считаем среднюю скорость.
	return distance / duration.Hours()
}

// meanSpeed — внутренняя версия с валидацией шагов для расчёта калорий.
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	return MeanSpeed(steps, height, duration)
}

// RunningSpentCalories возвращает количество калорий, потраченных при беге.
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if err := validate(steps, weight, height, duration); err != nil {
		return 0, err
	}

	speed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()

	return (weight * speed * durationInMinutes) / minInH, nil
}

// WalkingSpentCalories возвращает количество калорий, потраченных при ходьбе.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if err := validate(steps, weight, height, duration); err != nil {
		return 0, err
	}

	speed := meanSpeed(steps, height, duration)
	durationInMinutes := duration.Minutes()

	calories := (weight * speed * durationInMinutes) / minInH

	return calories * walkingCaloriesCoefficient, nil
}

// validate проверяет входные параметры на корректность.
func validate(steps int, weight, height float64, duration time.Duration) error {
	if steps < 0 {
		return ErrNegativeSteps
	}
	if steps == 0 {
		return ErrZeroSteps
	}
	if weight <= 0 {
		return ErrZeroWeight
	}
	if height <= 0 {
		return ErrZeroHeight
	}
	if duration <= 0 {
		return ErrZeroDuration
	}
	return nil
}
