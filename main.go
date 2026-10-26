package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ApiResponse struct {
	Status string `json:"status"`
	Data   any    `json:"data"`
}

type User struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}
type Handler struct {
	Users []User
}

func SuccessResponseHandler(w http.ResponseWriter, status int, data any, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ApiResponse{
		Status: "success",
		Data: data,
	})
}
func newHandler() *Handler {
	return &Handler{
		Users: []User{
			{
				Id:    1,
				Name:  "Budi",
				Email: "budi@example.com",
				Phone: "08123456789",
			},
			{
				Id:    2,
				Name:  "Andi",
				Email: "andi@example.com",
				Phone: "08987654321",
			},
		},
	}
}
func (h *Handler) getUsers(w http.ResponseWriter, r *http.Request) {
	SuccessResponseHandler(w, http.StatusOK, h.Users, r)
}
func main() {
	userHandler := newHandler()
	userRoutes := http.NewServeMux()
	userRoutes.HandleFunc("GET /users", userHandler.getUsers)

	err := http.ListenAndServe(":8000", userRoutes)
	if err != nil {
		fmt.Println("Gagal tersambung ke port 8000")
	}

}
