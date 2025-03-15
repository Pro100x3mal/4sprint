package daysteps

import (
	"errors"
	"fmt"
	"github.com/Yandex-Practicum/go1fl-4-sprint-final/internal/spentcalories"
	"strconv"
	"strings"
	"time"
)

const (
	StepLength = 0.65 // длина шага в метрах
	mInKm      = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	slice := strings.Split(data, ",")
	if len(slice) != 2 {
		return 0, 0, errors.New("некорректная строка - требуется 2 параметра разделенных запятыми")
	}
	steps, err := strconv.Atoi(slice[0])
	if err != nil {
		return 0, 0, errors.New("неверный тип данных - требуется целое число")
	}
	duration, err := time.ParseDuration(slice[1])
	if err != nil {
		return 0, 0, errors.New("неверный тип данных - требуется время в формате 3h50m")
	}
	return steps, duration, nil
}

// DayActionInfo обрабатывает входящий пакет, который передаётся в
// виде строки в параметре data. Параметр storage содержит пакеты за текущий день.
// Если время пакета относится к новым суткам, storage предварительно
// очищается.
// Если пакет валидный, он добавляется в слайс storage, который возвращает
// функция. Если пакет невалидный, storage возвращается без изменений.
func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if steps <= 0 {
		return ""
	}
	distance := float64(steps) * StepLength / mInKm
	calories := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	return fmt.Sprintf("Количество шагов: %v.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.", steps, distance, calories)
}
