package api

import (
	"net/http"

	"github.com/aleksandradidenko/go_final_project/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	writeJson(w, http.StatusOK, task)
}
