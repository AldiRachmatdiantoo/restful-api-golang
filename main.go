package main

import (
	"belajar-golang/controllers/user"
	"fmt"
	"net/http"
)

func main() {
	routes := http.NewServeMux()
	routes.HandleFunc("GET /users", user.Index)
	routes.HandleFunc("POST /users", user.Create)

	err := http.ListenAndServe(":8000", routes)
	if err != nil{
		fmt.Println(err.Error())
		return
	}
}
