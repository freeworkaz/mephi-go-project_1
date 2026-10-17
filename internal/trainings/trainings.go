package trainings

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

// Типы тренировок.
const (
	walkingTraining = "Ходьба"
	runningTraining = "Бег"
)

// Training хранит данные о тренировке.
type Training struct {
	Steps        int           // количество шагов, проделанных за тренировку.
	TrainingType string        // тип тренировки (бег или ходьба).
	Duration     time.Duration // длительность тренировки.
	personaldata.Personal
}

// Parse парсит строку с данными формата "3456,Ходьба,3h00m" и записывает
// данные в соответствующие поля структуры Training.
func (t *Training) Parse(datastring string) (err error) {
	data := strings.Split(datastring, ",")
	if len(data) != 3 {
		return errors.New("trainings: некорректный формат данных тренировки")
	}

	steps, err := strconv.Atoi(strings.TrimSpace(data[0]))
	if err != nil {
		return fmt.Errorf("trainings: ошибка разбора количества шагов: %w", err)
	}
	if steps <= 0 {
		return errors.New("trainings: количество шагов должно быть положительным")
	}
	t.Steps = steps

	t.TrainingType = data[1]

	duration, err := time.ParseDuration(data[2])
	if err != nil {
		return fmt.Errorf("trainings: ошибка разбора продолжительности тренировки: %w", err)
	}
	if duration <= 0 {
		return errors.New("trainings: продолжительность тренировки должна быть положительной")
	}
	t.Duration = duration

	return nil
}

// ActionInfo формирует и возвращает строку с информацией о тренировке.
func (t Training) ActionInfo() (string, error) {
	distance := spentenergy.Distance(t.Steps, t.Height)
	meanSpeed := spentenergy.MeanSpeed(t.Steps, t.Height, t.Duration)

	var calories float64
	var err error

	switch t.TrainingType {
	case walkingTraining:
		calories, err = spentenergy.WalkingSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
	case runningTraining:
		calories, err = spentenergy.RunningSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}
	if err != nil {
		return "", err
	}

	info := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		t.TrainingType, t.Duration.Hours(), distance, meanSpeed, calories)

	return info, nil
}
