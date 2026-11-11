package api

import (
	"net/http"
	"time"

	"github.com/aleksandradidenko/go_final_project/pkg/db"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeJson(w, http.StatusMethodNotAllowed, map[string]string{"error": "Метод не поддерживается"})
		return
	}

	id := r.FormValue("id")
	if id == "" {
		writeJson(w, http.StatusOK, map[string]string{
			"error": "Не указан идентификатор задачи",
		})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, http.StatusOK, map[string]string{"error": err.Error()})
		return
	}

	if task.Repeat == "" {
		err = db.DeleteTask(id)
	} else {
		next, nextErr := NextDate(time.Now(), task.Date, task.Repeat)
		if nextErr != nil {
			writeJson(w, http.StatusOK, map[string]string{"error": nextErr.Error()})
			return
		}

		err = db.UpdateDate(id, next)
	}

	if err != nil {
		writeJson(w, http.StatusOK, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, http.StatusOK, map[string]string{})
}
