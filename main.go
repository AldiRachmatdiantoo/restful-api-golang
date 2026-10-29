package main

import (
	"belajar-golang/database"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type User struct {
	Id    int64  `json:"id"`
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required"`
	Phone string `json:"phone" validate:"required"`
	Role  string `json:"role" validate:"required"`
}
type ApiResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func SuccessResponse(w http.ResponseWriter, r *http.Request, status int, message string, data any) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ApiResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}
func ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message string) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ApiResponse{
		Status:  "error",
		Message: message,
		Data:    "",
	})
}

func main() {
	db := database.Connect()
	defer db.Close()

	router := http.NewServeMux()

	var validate = validator.New()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM users")
		if err != nil {
			ErrorResponse(w, r, http.StatusInternalServerError, "Gagal Melakukan Query Database")
			return
		}
		defer rows.Close()
		users := []User{}
		for rows.Next() {
			catchUser := User{}
			err := rows.Scan(&catchUser.Id, &catchUser.Name, &catchUser.Email, &catchUser.Phone, &catchUser.Role)
			if err != nil {
				ErrorResponse(w, r, http.StatusInternalServerError, "Error: "+err.Error())
				return
			}
			users = append(users, catchUser)
		}
		if err = rows.Err(); err != nil {
			ErrorResponse(w, r, http.StatusInternalServerError, "Gagal Mengambil data!")
			return
		}
		SuccessResponse(w, r, http.StatusOK, "Berhasil mengambil data!", users)
	})
	router.HandleFunc("GET /{id}", func(w http.ResponseWriter, r *http.Request) {
		userId, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			ErrorResponse(w, r, http.StatusBadRequest, "ID Salah!")
			return
		}
		getUser := User{}
		stmt, err := db.Prepare("SELECT * FROM users WHERE id = ?")
		if err != nil {
			ErrorResponse(w, r, http.StatusBadRequest, "Gagal melakukan Query Database!"+err.Error())
			return
		}
		defer stmt.Close()
		err = stmt.QueryRow(userId).Scan(&getUser.Id, &getUser.Name, &getUser.Email, &getUser.Phone, &getUser.Role)
		if err == sql.ErrNoRows {
			ErrorResponse(w, r, http.StatusBadRequest, "User Tidak Ditemukan!")
			return
		} else if err != nil {
			ErrorResponse(w, r, http.StatusInternalServerError, "Error: "+err.Error())
			return
		}
		SuccessResponse(w, r, http.StatusOK, "Berhasil Mengambil Data!", getUser)

	})
	router.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		getUser := User{}
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&getUser)
		if err != nil {
			ErrorResponse(w, r, http.StatusBadRequest, "error: "+err.Error())
			return
		}
		err = validate.Struct(getUser)
		if err != nil {
			ErrorResponse(w, r, http.StatusBadRequest, "Input salah!")
			return
		}
		stmt, err := db.Prepare("INSERT INTO users (name, email, phone, role) VALUES (?, ?, ?, ?)")
		if err != nil {
			ErrorResponse(w, r, http.StatusBadRequest, "Gagal melakukan Query Database! "+err.Error())
			return
		}
		defer stmt.Close()
		result, err := stmt.Exec(getUser.Name, getUser.Email, getUser.Phone, getUser.Role)
		if err != nil {
			ErrorResponse(w, r, http.StatusBadRequest, "Gagal melakukan Query Database! "+err.Error())
			return
		}
		id, err := result.LastInsertId()
		if err != nil {
			ErrorResponse(w, r, http.StatusNotFound, "Tidak dapat menyimpan data!")
			return
		}
		getUser.Id = id
		SuccessResponse(w, r, http.StatusCreated, "Berhasil menambahkan user!", getUser)
	})
	router.HandleFunc("PUT /{id}", func(w http.ResponseWriter, r *http.Request) {
		users := User{}
		data := json.NewDecoder(r.Body)
		data.DisallowUnknownFields()
		err := data.Decode(&users)
		if err != nil {
			ErrorResponse(w, r, http.StatusBadRequest, "error: "+err.Error())
			return
		}
		err = validate.Struct(users)
		if err != nil {
			ErrorResponse(w, r, http.StatusBadRequest, "Input salah!")
			return
		}
		userId, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			ErrorResponse(w, r, http.StatusBadRequest, "ID Salah!")
			return
		}
		stmt, err := db.Prepare("UPDATE users SET name = ?, email = ?, phone = ?, role = ? WHERE id = ?")
		if err != nil {
			ErrorResponse(w, r, http.StatusBadRequest, "Gagal melakukan Query Database! "+err.Error())
			return
		}
		defer stmt.Close()
		_, err = stmt.Exec(users.Name, users.Email, users.Phone, users.Role, userId)
		if err != nil {
			ErrorResponse(w, r, http.StatusBadRequest, "Gagal melakukan Query Database! "+err.Error())
			return
		}
		SuccessResponse(w, r, http.StatusCreated, "Berhasil mengupdate user!", []User{})
	})
	router.HandleFunc("DELETE /{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			ErrorResponse(w, r, http.StatusBadRequest, "ID Salah!")
			return
		}
		stmt, err := db.Prepare("DELETE FROM users WHERE id = ?")
		if err != nil {
			ErrorResponse(w, r, http.StatusBadRequest, "Gagal melakukan Query Database! "+err.Error())
			return
		}
		defer stmt.Close()
		_, err = stmt.Exec(id)
		if err != nil {
			ErrorResponse(w, r, http.StatusBadRequest, "Gagal melakukan Query Database! "+err.Error())
			return
		}
		SuccessResponse(w, r, http.StatusOK, "Berhasil Menghapus user!", []User{})
	})
	err := http.ListenAndServe("localhost:8080", router)
	if err != nil {
		log.Fatal(err)
	}

}
