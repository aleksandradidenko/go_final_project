package api

import (
	"net/http"

	"github.com/aleksandradidenko/go_final_project/pkg/db"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	if err := validateTaskID(id); err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		writeDatabaseError(w, err)
		return
	}

	writeJson(w, http.StatusOK, map[string]string{})
}
