package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep = 0.65 // средняя длина шага.
	mInKm   = 1000 // количество метров в километре.
	minInH  = 60   // количество минут в часе.
	//kmhInMsec = 0.278 // коэффициент для преобразования км/ч в м/с.
	//cmInM     = 100   // количество сантиметров в метре.
)

func parseTraining(data string) (int, string, time.Duration, error) {
	slice := strings.Split(data, ",")
	if len(slice) != 3 {
		return 0, "", 0, errors.New("некорректная строка - требуется 3 параметра разделенных запятыми")
	}
	steps, err := strconv.Atoi(slice[0])
	if err != nil {
		return 0, "", 0, errors.New("неверный тип данных - требуется целое число")
	}
	duration, err := time.ParseDuration(slice[2])
	if err != nil {
		return 0, "", 0, errors.New("неверный тип данных - требуется время в формате 3h50m")
	}
	return steps, slice[1], duration, nil
}

// distance возвращает дистанцию(в километрах), которую преодолел пользователь за время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий (число шагов при ходьбе и беге).
func distance(steps int) float64 {
	return float64(steps) * lenStep / mInKm
}

// meanSpeed возвращает значение средней скорости движения во время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий(число шагов при ходьбе и беге).
// duration time.Duration — длительность тренировки.
func meanSpeed(steps int, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	return distance(steps) / duration.Hours()
}

// ShowTrainingInfo возвращает строку с информацией о тренировке.
//
// Параметры:
//
// data string - строка с данными.
// weight, height float64 — вес и рост пользователя.
func TrainingInfo(data string, weight, height float64) string {
	steps, trainingType, duration, err := parseTraining(data)
	if err != nil {
		fmt.Println(err)
		return "неизвестный тип тренировки"
	}
	dist := distance(steps)
	speed := meanSpeed(steps, duration)
	var calories float64
	switch trainingType {
	case "Ходьба":
		calories = WalkingSpentCalories(steps, weight, height, duration)
	case "Бег":
		calories = RunningSpentCalories(steps, weight, duration)
	}
	return fmt.Sprintf("Тип тренировки: %v\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f", trainingType, duration.Hours(), dist, speed, calories)

}

// Константы для расчета калорий, расходуемых при беге.
const (
	runningCaloriesMeanSpeedMultiplier = 18.0 // множитель средней скорости.
	runningCaloriesMeanSpeedShift      = 20.0 // среднее количество сжигаемых калорий при беге.
)

// RunningSpentCalories возвращает количество потраченных колорий при беге.
//
// Параметры:
//
// steps int - количество шагов.
// weight float64 — вес пользователя.
// duration time.Duration — длительность тренировки.
func RunningSpentCalories(steps int, weight float64, duration time.Duration) float64 {
	return ((runningCaloriesMeanSpeedMultiplier * meanSpeed(steps, duration)) - runningCaloriesMeanSpeedShift) * weight
}

// Константы для расчета калорий, расходуемых при ходьбе.
const (
	walkingCaloriesWeightMultiplier = 0.035 // множитель массы тела.
	walkingSpeedHeightMultiplier    = 0.029 // множитель роста.
)

// WalkingSpentCalories возвращает количество потраченных калорий при ходьбе.
//
// Параметры:
//
// steps int - количество шагов.
// duration time.Duration — длительность тренировки.
// weight float64 — вес пользователя.
// height float64 — рост пользователя.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) float64 {
	speed := meanSpeed(steps, duration)
	return ((walkingCaloriesWeightMultiplier * weight) + (speed*speed/height)*walkingSpeedHeightMultiplier) * duration.Hours() * minInH
}
