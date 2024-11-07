package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	// other imports you need
)

type Student struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

var (
	students  = make(map[int]Student)
	idCounter = 1
	mu        sync.Mutex // Protects the students map and idCounter
)

func createStudent(w http.ResponseWriter, r *http.Request) {
	var student Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	mu.Lock()
	student.ID = idCounter
	students[idCounter] = student
	idCounter++
	mu.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(student)
}

func getStudents(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	var allStudents []Student
	for _, student := range students {
		allStudents = append(allStudents, student)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(allStudents)
}

func getStudentByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	mu.Lock()
	student, exists := students[id]
	mu.Unlock()

	if !exists {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(student)
}

func updateStudent(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	var updatedStudent Student
	if err := json.NewDecoder(r.Body).Decode(&updatedStudent); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	mu.Lock()
	_, exists := students[id]
	if exists {
		updatedStudent.ID = id
		students[id] = updatedStudent
	}
	mu.Unlock()

	if !exists {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(updatedStudent)
}

func deleteStudent(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	mu.Lock()
	_, exists := students[id]
	if exists {
		delete(students, id)
	}
	mu.Unlock()

	if !exists {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getStudentSummary(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	mu.Lock()
	student, exists := students[id]
	mu.Unlock()

	if !exists {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	// Call Ollama API with the student data for a summary
	summary, err := generateOllamaSummary(student)
	if err != nil {
		http.Error(w, "Error generating summary ", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(summary))
}

func generateOllamaSummary(student Student) (string, error) {
	requestBody, _ := json.Marshal(map[string]interface{}{
		"prompt": "Generate a profile summary for a student",
		"context": map[string]string{
			"Name":  student.Name,
			"Age":   strconv.Itoa(student.Age),
			"Email": student.Email,
		},
	})

	resp, err := http.Post("http://localhost:11434/ollama/api", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "Error", err
	}
	defer resp.Body.Close()

	var response map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	return response["summary"], nil
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/students", createStudent).Methods("POST")
	r.HandleFunc("/students", getStudents).Methods("GET")
	r.HandleFunc("/students/{id}", getStudentByID).Methods("GET")
	r.HandleFunc("/students/{id}", updateStudent).Methods("PUT")
	r.HandleFunc("/students/{id}", deleteStudent).Methods("DELETE")
	r.HandleFunc("/students/{id}/summary", getStudentSummary).Methods("GET")

	http.ListenAndServe(":8080", r)
}
