package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type Task struct {
	ID           string   `json:"id"`          // ID задачи
	Description  string   `json:"description"` // Заголовок
	Note         string   `json:"note"`        // Описание задачи
	Applications []string `json:"application"` // Приложения, которыми будете пользоваться
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Обработчик для получения всех задач
func getTasks(w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(tasks)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)

	if err != nil {
		fmt.Println("Ошибка записи тела ответа: ", err)
	}
}

// Обработчик для отправки задачи на сервер
func createTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, ok := tasks[task.ID]

	if ok {
		http.Error(w, "Задача c указанным id уже существует", http.StatusBadRequest)
		return
	}

	if len(task.Applications) == 0 {
		app := strings.Split(r.Header.Get("User-Agent"), "/")[0]
		task.Applications = append(task.Applications, app)
	}

	tasks[task.ID] = task
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// Обработчик для получения задачи по ID
func getTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	task, ok := tasks[id]

	if !ok {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(task)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)

	if err != nil {
		fmt.Println("Ошибка записи тела ответа: ", err)
	}
}

// Обработчик удаления задачи по ID
func deleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, ok := tasks[id]

	if !ok {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}

	delete(tasks, id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	r.Get("/tasks", getTasks)
	r.Post("/tasks", createTask)
	r.Get("/tasks/{id}", getTask)
	r.Delete("/tasks/{id}", deleteTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
