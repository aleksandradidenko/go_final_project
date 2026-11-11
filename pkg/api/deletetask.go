package api

import (
	"net/http"

	"github.com/aleksandradidenko/go_final_project/pkg/db"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	if id == "" {
		writeJson(w, http.StatusOK, map[string]string{
			"error": "Не указан идентификатор задачи",
		})
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		writeJson(w, http.StatusOK, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJson(w, http.StatusOK, map[string]string{})
}
