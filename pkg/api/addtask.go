package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/aleksandradidenko/go_final_project/pkg/db"
)

// Проверка даты и правила повторения задачи
func checkDate(task *db.Task) error {
	if len(task.Repeat) > 128 {
		return fmt.Errorf("правило повторения не должно превышать 128 символов")
	}

	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(DateLayout)
	}

	t, err := time.Parse(DateLayout, task.Date)
	if err != nil {
		return err
	}

	var next string

	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	if afterNow(now, t) {
		if task.Repeat == "" {
			task.Date = now.Format(DateLayout)
		} else {
			task.Date = next
		}
	}

	return nil
}

// Отправка ответа в формате JSON
func writeJson(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Ошибка записи JSON-ответа: %v", err)
	}
}

// Обработчик добавления задачи
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if task.Title == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{
			"error": "Не указан заголовок задачи",
		})
		return
	}

	if err := checkDate(&task); err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeDatabaseError(w, err)
		return
	}

	writeJson(w, http.StatusOK, map[string]string{
		"id": strconv.FormatInt(id, 10),
	})
}
