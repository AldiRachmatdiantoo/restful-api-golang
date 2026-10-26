package user

import (
	"belajar-golang/database"
	"encoding/json"
	"fmt"
	"net/http"
)

type user struct {
	Id    int64    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}
type ApiResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func SuccessResponse(w http.ResponseWriter, r *http.Request, status int, data any, msg string) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ApiResponse{
		Status:  "success",
		Message: msg,
		Data:    data,
	})
}
func ErrorResponse(w http.ResponseWriter, r *http.Request, status int, data any, msg string) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ApiResponse{
		Status:  "error",
		Message: msg,
		Data:    nil,
	})
}
func Index(w http.ResponseWriter, r *http.Request) {
	db, err := database.Connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()
	var result []user
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer rows.Close()
	for rows.Next() {
		each := user{}
		err := rows.Scan(&each.Id, &each.Name, &each.Email, &each.Phone)
		if err != nil {
			ErrorResponse(w, r, http.StatusInternalServerError, nil, "Gagal Mengambil Data!")
			return
		}
		result = append(result, each)
	}
	if err = rows.Err(); err != nil {
		fmt.Println(err.Error())
		return
	}
	SuccessResponse(w, r, http.StatusOK, result, "Berhasil Mengambil Data!")

}
func Create(w http.ResponseWriter, r *http.Request) {

	var data user
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		ErrorResponse(w, r, http.StatusBadRequest, nil, "Data salah!")
	}
	db, err := database.Connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO users (name, email, phone) VALUES (?, ?, ?)")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	result, err := stmt.Exec(data.Name, data.Email, data.Phone)
	id, err := result.LastInsertId()
	if err != nil {
		return
	}
	data.Id = id
	SuccessResponse(w, r, http.StatusCreated, data, "Berhasil menambah Data!")
}
