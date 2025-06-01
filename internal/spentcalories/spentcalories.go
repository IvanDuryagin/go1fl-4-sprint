package spentcalories

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	// TODO: реализовать функцию
	stepsActivityDuration := strings.Split(data, ",")
	if len(stepsActivityDuration) != 3 {
		return 0, "", 0, fmt.Errorf("неверный формат данных: ожидается 'шаги,активность,время', получено: '%s'", data)
	}

	stepsStr := stepsActivityDuration[0]
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("неверный формат количества шагов '%s':'%v'", stepsStr, err)
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("количество шагов < 0 '%d'", steps)
	}

	activity := stepsActivityDuration[2]
	if activity == "" {
		return 0, "", 0, fmt.Errorf("вид активности не может быть пустым")
	}

	durationStr := stepsActivityDuration[2]
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("неверный формат продолжительности '%s':'%v'", durationStr, err)
	}
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("ошибка: Время < 0 '%s'", duration)
	}
	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	// TODO: реализовать функцию
	stepLength := height * stepLengthCoefficient

	distanceM := float64(steps) * stepLength

	distanceKm := distanceM / mInKm

	return distanceKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// TODO: реализовать функцию
	if duration <= 0 {
		return 0
	}

	distanceKm := distance(steps, height)
	hours := duration.Hours()

	if hours == 0 {
		return 0
	}

	return distanceKm / hours
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	// TODO: реализовать функцию
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return "", fmt.Errorf("ошибка парсинга данных тренировки: %v", err)
	}

	if weight <= 0 {
		return "", fmt.Errorf("некорректный вес: %.2f кг", weight)
	}
	if height <= 0 {
		return "", fmt.Errorf("некорректный рост: %.2f м", height)
	}

	var dist, speed, calories float64
	var calcErr error

	switch activity {
	case "Ходьба":
		dist = distance(steps, height)
		speed = meanSpeed(steps, height, duration)
		calories, calcErr = WalkingSpentCalories(steps, weight, height, duration)
	case "Бег":
		dist = distance(steps, height)
		speed = meanSpeed(steps, height, duration)
		calories, calcErr = RunningSpentCalories(steps, weight, height, duration)
	default:
		return "", fmt.Errorf("неизвестный тип тренировки: %s", activity)
	}

	if calcErr != nil {
		return "", fmt.Errorf("ошибка расчета показателей: %v", calcErr)
	}

	durationHours := duration.Hours()
	result := fmt.Sprintf(
		"Тип тренировки: %s\n"+
			"Длительность: %.2f ч.\n"+
			"Дистанция: %.2f км.\n"+
			"Скорость: %.2f км/ч\n"+
			"Сожгли калорий: %.2f",
		activity,
		durationHours,
		dist,
		speed,
		calories,
	)

	return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// TODO: реализовать функцию
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов < 0, получено: %d", steps)
	}

	if weight <= 0 {
		return 0, fmt.Errorf("вес < 0, получено: %.2f", weight)
	}

	if height <= 0 {
		return 0, fmt.Errorf("рост < 0, получено: %.2f", height)
	}

	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность < 0, получено: %v", duration)
	}

	speed := meanSpeed(steps, height, duration)
	if speed <= 0 {
		return 0, fmt.Errorf("не удалось рассчитать среднюю скорость")
	}

	durationInMinutes := duration.Minutes()

	calories := (weight * speed * durationInMinutes) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// TODO: реализовать функцию
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов < 0, получено: %d", steps)
	}

	if weight <= 0 {
		return 0, fmt.Errorf("вес < 0, получено: %.2f", weight)
	}

	if height <= 0 {
		return 0, fmt.Errorf("рост < 0, получено: %.2f", height)
	}

	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность < 0, получено: %v", duration)
	}

	speed := meanSpeed(steps, height, duration)
	if speed <= 0 {
		return 0, fmt.Errorf("не удалось рассчитать среднюю скорость")
	}

	durationInMinutes := duration.Minutes()

	calories := (weight * speed * durationInMinutes) / minInH

	return calories * walkingCaloriesCoefficient, nil
}
