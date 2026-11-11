package api

import (
	"net/http"

	"github.com/aleksandradidenko/go_final_project/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeJson(w, http.StatusMethodNotAllowed, map[string]string{"error": "Метод не поддерживается"})
		return
	}

	search := r.FormValue("search")

	tasks, err := db.Tasks(50, search)
	if err != nil {
		writeJson(w, http.StatusOK, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJson(w, http.StatusOK, TasksResp{
		Tasks: tasks,
	})
}
