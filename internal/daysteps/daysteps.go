package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	stepLength = 0.65 // Длина одного шага в метрах
	mInKm      = 1000 // Количество метров в километре
)

func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("неверный формат данных: ожидается 'шаги,время', получено: '%s'", data)
	}

	// Обработка шагов
	stepStr := strings.TrimSpace(parts[0])
	if stepStr == "" {
		return 0, 0, fmt.Errorf("количество шагов не может быть пустым")
	}

	// Удаляем + в начале если есть
	if strings.HasPrefix(stepStr, "+") {
		stepStr = stepStr[1:]
	}

	steps, err := strconv.Atoi(stepStr)
	if err != nil {
		return 0, 0, fmt.Errorf("неверный формат количества шагов '%s': %v", parts[0], err)
	}
	if steps <= 0 {
		return 0, 0, fmt.Errorf("количество шагов должно быть положительным, получено: %d", steps)
	}

	// Обработка времени
	durationStr := strings.TrimSpace(parts[1])
	if durationStr == "" {
		return 0, 0, fmt.Errorf("продолжительность не может быть пустой")
	}

	// Заменяем . на h для дробных часов (например, 1.5h -> 1h30m)
	if strings.Contains(durationStr, ".") && strings.Contains(durationStr, "h") {
		parts := strings.Split(durationStr, "h")
		if len(parts) != 2 {
			return 0, 0, fmt.Errorf("неверный формат продолжительности '%s'", durationStr)
		}
		hours, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return 0, 0, fmt.Errorf("неверный формат продолжительности '%s': %v", durationStr, err)
		}
		wholeHours := int(hours)
		minutes := int((hours - float64(wholeHours)) * 60)
		durationStr = fmt.Sprintf("%dh%dm", wholeHours, minutes)
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, 0, fmt.Errorf("неверный формат продолжительности '%s': %v", parts[1], err)
	}
	if duration <= 0 {
		return 0, 0, fmt.Errorf("продолжительность должна быть положительной, получено: %v", duration)
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println("Ошибка:", err)
		return ""
	}

	distanceKm := float64(steps) * stepLength / mInKm
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println("Ошибка расчета калорий:", err)
		return ""
	}

	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps,
		distanceKm,
		calories,
	)
}
