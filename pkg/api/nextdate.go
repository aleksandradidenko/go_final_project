package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DateLayout = "20060102"

func afterNow(date, now time.Time) bool {
	return date.Format(DateLayout) > now.Format(DateLayout)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("не указано правило повторения")
	}

	date, err := time.Parse(DateLayout, dstart)
	if err != nil {
		return "", err
	}

	parts := strings.Split(repeat, " ")

	switch parts[0] {

	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("неверный формат правила d")
		}

		days, err := strconv.Atoi(parts[1])
		if err != nil || days < 1 || days > 400 {
			return "", fmt.Errorf("недопустимый интервал повторения")
		}

		for {
			date = date.AddDate(0, 0, days)

			if afterNow(date, now) {
				break
			}
		}

	case "y":
		if len(parts) != 1 {
			return "", fmt.Errorf("неверный формат правила y")
		}

		for {
			date = date.AddDate(1, 0, 0)

			if afterNow(date, now) {
				break
			}
		}

	case "w":
		if len(parts) != 2 {
			return "", fmt.Errorf("неверный формат правила w")
		}

		var weekdays [8]bool

		for _, value := range strings.Split(parts[1], ",") {
			day, err := strconv.Atoi(value)
			if err != nil || day < 1 || day > 7 {
				return "", fmt.Errorf("недопустимый день недели")
			}

			weekdays[day] = true
		}

		for {
			date = date.AddDate(0, 0, 1)

			weekday := int(date.Weekday())
			if weekday == 0 {
				weekday = 7
			}

			if weekdays[weekday] && afterNow(date, now) {
				break
			}
		}

	case "m":
		if len(parts) < 2 || len(parts) > 3 {
			return "", fmt.Errorf("неверный формат правила m")
		}

		var days [32]bool
		var months [13]bool

		var lastDay bool
		var beforeLastDay bool

		// Проверяем дни месяца
		for _, value := range strings.Split(parts[1], ",") {
			day, err := strconv.Atoi(value)
			if err != nil {
				return "", fmt.Errorf("неверный день месяца")
			}

			switch {
			case day >= 1 && day <= 31:
				days[day] = true

			case day == -1:
				lastDay = true

			case day == -2:
				beforeLastDay = true

			default:
				return "", fmt.Errorf("недопустимый день месяца")
			}
		}

		// Проверяем месяцы
		if len(parts) == 3 {
			for _, value := range strings.Split(parts[2], ",") {
				month, err := strconv.Atoi(value)
				if err != nil || month < 1 || month > 12 {
					return "", fmt.Errorf("недопустимый месяц")
				}

				months[month] = true
			}
		} else {
			for month := 1; month <= 12; month++ {
				months[month] = true
			}
		}

		// Проверяем, существует ли хотя бы одна подходящая дата
		possible := lastDay || beforeLastDay

		for month := 1; month <= 12 && !possible; month++ {
			if !months[month] {
				continue
			}

			daysInMonth := time.Date(
				2000, time.Month(month)+1, 0,
				0, 0, 0, 0, time.UTC,
			).Day()

			for day := 1; day <= daysInMonth; day++ {
				if days[day] {
					possible = true
					break
				}
			}
		}

		if !possible {
			return "", fmt.Errorf("нет подходящих дат для повторения")
		}

		// Ищем следующую подходящую дату
		for {
			date = date.AddDate(0, 0, 1)

			if !afterNow(date, now) {
				continue
			}

			if !months[int(date.Month())] {
				continue
			}

			day := date.Day()

			isLast := date.AddDate(0, 0, 1).Month() != date.Month()

			isBeforeLast := !isLast &&
				date.AddDate(0, 0, 2).Month() != date.Month()

			if days[day] ||
				(lastDay && isLast) ||
				(beforeLastDay && isBeforeLast) {
				break
			}
		}

	default:
		return "", fmt.Errorf("неподдерживаемое правило повторения")
	}

	return date.Format(DateLayout), nil
}

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	now := time.Now()

	if value := r.FormValue("now"); value != "" {
		var err error
		now, err = time.Parse(DateLayout, value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	next, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, next)
}
