package api

import (
	"encoding/json"
	"net/http"

	"github.com/aleksandradidenko/go_final_project/pkg/db"
)

func editTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJson(w, http.StatusOK, map[string]string{"error": err.Error()})
		return
	}

	if task.ID == "" {
		writeJson(w, http.StatusOK, map[string]string{"error": "Не указан идентификатор задачи"})
		return
	}

	if task.Title == "" {
		writeJson(w, http.StatusOK, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	if err := checkDate(&task); err != nil {
		writeJson(w, http.StatusOK, map[string]string{"error": err.Error()})
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		writeJson(w, http.StatusOK, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, http.StatusOK, map[string]string{})
}
