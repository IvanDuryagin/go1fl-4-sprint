package spentcalories

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	lenStep                    = 0.65
	mInKm                      = 1000
	minInH                     = 60
	stepLengthCoefficient      = 0.45
	walkingCaloriesCoefficient = 0.5
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("неверный формат данных: ожидается 'шаги,активность,время', получено: '%s'", data)
	}

	// Parse steps
	stepStr := strings.TrimSpace(parts[0])
	if stepStr == "" {
		return 0, "", 0, fmt.Errorf("количество шагов не может быть пустым")
	}

	// Handle positive sign
	if strings.HasPrefix(stepStr, "+") {
		stepStr = stepStr[1:]
	}

	steps, err := strconv.Atoi(stepStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("неверный формат количества шагов '%s': %v", parts[0], err)
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("количество шагов должно быть положительным, получено: %d", steps)
	}

	// Parse activity
	activity := strings.TrimSpace(parts[1])
	if activity == "" {
		return 0, "", 0, fmt.Errorf("вид активности не может быть пустым")
	}

	// Parse duration
	durationStr := strings.TrimSpace(parts[2])
	if durationStr == "" {
		return 0, "", 0, fmt.Errorf("продолжительность не может быть пустой")
	}

	// Handle fractional hours (e.g., 1.5h -> 1h30m)
	if strings.Contains(durationStr, ".") && strings.Contains(durationStr, "h") {
		timeParts := strings.Split(durationStr, "h")
		if len(timeParts) != 2 {
			return 0, "", 0, fmt.Errorf("неверный формат продолжительности '%s'", durationStr)
		}
		hours, err := strconv.ParseFloat(timeParts[0], 64)
		if err != nil {
			return 0, "", 0, fmt.Errorf("неверный формат продолжительности '%s': %v", durationStr, err)
		}
		wholeHours := int(hours)
		minutes := int((hours - float64(wholeHours)) * 60)
		durationStr = fmt.Sprintf("%dh%dm", wholeHours, minutes)
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("неверный формат продолжительности '%s': %v", parts[2], err)
	}
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("продолжительность должна быть положительной, получено: %v", duration)
	}

	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	if steps <= 0 || height <= 0 {
		return 0
	}
	return float64(steps) * height * stepLengthCoefficient / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if steps <= 0 || height <= 0 || duration <= 0 {
		return 0
	}
	dist := distance(steps, height)
	hours := duration.Hours()
	if hours <= 0 {
		return 0
	}
	return dist / hours
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть положительным, получено: %d", steps)
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть положительным, получено: %.2f", weight)
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть положительным, получено: %.2f", height)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность должна быть положительной, получено: %v", duration)
	}

	speed := meanSpeed(steps, height, duration)
	if speed <= 0 {
		return 0, fmt.Errorf("не удалось рассчитать среднюю скорость")
	}

	return (weight * speed * duration.Minutes()) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	calories, err := RunningSpentCalories(steps, weight, height, duration)
	if err != nil {
		return 0, err
	}
	return calories * walkingCaloriesCoefficient, nil
}

func TrainingInfo(data string, weight, height float64) (string, error) {
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

	var calories float64
	var calcErr error

	switch activity {
	case "Ходьба":
		calories, calcErr = WalkingSpentCalories(steps, weight, height, duration)
	case "Бег":
		calories, calcErr = RunningSpentCalories(steps, weight, height, duration)
	default:
		return "", fmt.Errorf("неизвестный тип тренировки: %s", activity)
	}

	if calcErr != nil {
		return "", fmt.Errorf("ошибка расчета показателей: %v", calcErr)
	}

	return fmt.Sprintf(
		"Тип тренировки: %s\n"+
			"Длительность: %.2f ч.\n"+
			"Дистанция: %.2f км.\n"+
			"Скорость: %.2f км/ч\n"+
			"Сожгли калорий: %.2f\n",
		activity,
		duration.Hours(),
		distance(steps, height),
		meanSpeed(steps, height, duration),
		calories,
	), nil
}
