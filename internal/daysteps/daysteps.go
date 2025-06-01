package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	stepsDuration := strings.Split(data, ",")
	if len(stepsDuration) != 2 {
		return 0, 0, fmt.Errorf("неверный формат данных: ожидается 'шаги,время', получено: '%s'", data)
	}

	stepsStr := stepsDuration[0]
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, 0, fmt.Errorf("неверный формат количества шагов '%s':'%v'", stepsStr, err)
	}
	if steps <= 0 {
		return 0, 0, fmt.Errorf("количество шагов < 0 '%d'", steps)
	}

	durationStr := stepsDuration[1]
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, 0, fmt.Errorf("неверный формат продолжительности '%s':'%v'", durationStr, err)
	}
	if duration <= 0 {
		return 0, 0, fmt.Errorf("ошибка: Время < 0 '%s'", duration)
	}
	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return ""
	}

	if steps <= 0 {
		return ""
	}

	distanceM := float64(steps) * stepLength
	distanceKm := distanceM / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		fmt.Println("Ошибка расчета калорий:", err)
		return ""
	}

	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.",
		steps,
		distanceKm,
		calories,
	)
}
