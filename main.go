package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/TerrenceHo/CalHacks4-Backend/controllers"
	"github.com/TerrenceHo/CalHacks4-Backend/models"
	"github.com/gorilla/mux"
)

const (
	host     = "localhost"
	port     = "5432"
	user     = "kho"
	password = ""
	name     = "calhacks"
	pepper   = "dev-pepper"
)

func main() {
	connection := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", host, port, user, name)
	services, err := models.NewServices(
		models.WithGorm("postgres", connection),
		models.WithLogMode(true),
		models.WithUser(pepper),
	)
	must(err)
	defer services.Close()
	err = services.AutoMigrate()
	must(err)

	usersC := controllers.NewUsers(services.User)

	router := mux.NewRouter()
	router.HandleFunc("/", homePage).Methods("GET")
	router.HandleFunc("/api/v1/user/register", usersC.Create).Methods("POST")

	port := 12345
	fmt.Println("Listening on Port", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), router))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello World!</h1>")
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
