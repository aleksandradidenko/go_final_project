package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
)

func Init(password string) {
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/task", auth(taskHandler, password))
	http.HandleFunc("/api/tasks", auth(tasksHandler, password))
	http.HandleFunc("/api/task/done", auth(doneTaskHandler, password))
	http.HandleFunc("/api/signin", func(w http.ResponseWriter, r *http.Request) {
		signInHandler(w, r, password)
	})
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodPut:
		editTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		w.Header().Set("Allow", "GET, POST, PUT, DELETE")
		writeJson(w, http.StatusMethodNotAllowed, map[string]string{"error": "Метод не поддерживается"})
	}
}

func validateTaskID(id string) error {
	if id == "" {
		return errors.New("Не указан идентификатор задачи")
	}
	number, err := strconv.ParseInt(id, 10, 64)
	if err != nil || number <= 0 {
		return errors.New("Неверный идентификатор задачи")
	}
	return nil
}

func writeDatabaseError(w http.ResponseWriter, err error) {
	if errors.Is(err, sql.ErrNoRows) {
		writeJson(w, http.StatusNotFound, map[string]string{"error": "Задача не найдена"})
		return
	}

	log.Printf("Ошибка работы с базой данных: %v", err)
	writeJson(w, http.StatusInternalServerError, map[string]string{"error": "Ошибка работы с базой данных"})
}
