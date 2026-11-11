package api

import (
	"log"
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
	if err := validateTaskID(id); err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeDatabaseError(w, err)
		return
	}

	if task.Repeat == "" {
		err = db.DeleteTask(id)
	} else {
		next, nextErr := NextDate(time.Now(), task.Date, task.Repeat)
		if nextErr != nil {
			log.Printf("Ошибка расчёта следующей даты задачи %s: %v", id, nextErr)
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": "Не удалось рассчитать следующую дату"})
			return
		}

		err = db.UpdateDate(id, next)
	}

	if err != nil {
		writeDatabaseError(w, err)
		return
	}

	writeJson(w, http.StatusOK, map[string]string{})
}
