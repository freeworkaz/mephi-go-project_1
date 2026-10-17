package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

// DaySteps хранит данные о дневной прогулке.
type DaySteps struct {
	Steps    int           // количество шагов.
	Duration time.Duration // длительность прогулки.
	personaldata.Personal
}

// Parse парсит строку с данными формата "678,0h50m" и записывает данные
// в соответствующие поля структуры DaySteps.
func (ds *DaySteps) Parse(datastring string) (err error) {
	data := strings.Split(datastring, ",")
	if len(data) != 2 {
		return errors.New("daysteps: некорректный формат данных прогулки")
	}

	steps, err := strconv.Atoi(data[0])
	if err != nil {
		return fmt.Errorf("daysteps: ошибка разбора количества шагов: %w", err)
	}
	if steps <= 0 {
		return errors.New("daysteps: количество шагов должно быть положительным")
	}
	ds.Steps = steps

	duration, err := time.ParseDuration(data[1])
	if err != nil {
		return fmt.Errorf("daysteps: ошибка разбора продолжительности прогулки: %w", err)
	}
	if duration <= 0 {
		return errors.New("daysteps: продолжительность прогулки должна быть положительной")
	}
	ds.Duration = duration

	return nil
}

// ActionInfo формирует и возвращает строку с информацией о прогулке.
func (ds DaySteps) ActionInfo() (string, error) {
	distance := spentenergy.Distance(ds.Steps, ds.Height)

	calories, err := spentenergy.WalkingSpentCalories(ds.Steps, ds.Weight, ds.Height, ds.Duration)
	if err != nil {
		return "", err
	}

	info := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		ds.Steps, distance, calories)

	return info, nil
}
