package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Person struct {
	ID        string   `json:"id,omitempty"`
	FirstName string   `json:"firstname,omitempty"`
	LastName  string   `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
}

type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

var people []Person

func GetPeopleEndpoint(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(people)
}

func GetPersonEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for _, item := range people {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Person{})
}

func CreatePersonEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	var person Person
	_ = json.NewDecoder(req.Body).Decode(&person) //Aqui recojo el json de persona y lo decodifico
	person.ID = params["id"]                      //Este paramatro se lo paso por la url en vez de por el json y lo añado a person
	people = append(people, person)               //Añado el nuevo objeto persona al array de people
	json.NewEncoder(w).Encode(people)             //Devuelvo el array entero de people con el nuevo person incluido
}

func DeletePersonEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for index, item := range people {
		if item.ID == params["id"] {
			people = append(people[:index], people[index+1:]...) //Cuando lo encuentre en el index que sea lo borra con los ... y vuelve a conectar el array a la siguiente persona
			break
		}
	}
	json.NewEncoder(w).Encode(people)
}

func main() {
	router := mux.NewRouter()

	// Creo unos valores porque no tengo la bbdd y devuelvo desde el array dos personas
	people = append(people, Person{ID: "1", FirstName: "Ryan", LastName: "Ray", Address: &Address{City: "Dubling", State: "California"}})

	people = append(people, Person{ID: "2", FirstName: "Noel", LastName: "Cerro"})

	// endpoints
	router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
	router.HandleFunc("/people/{id}", GetPersonEndpoint).Methods("GET")
	router.HandleFunc("/people/{id}", CreatePersonEndpoint).Methods("POST")
	router.HandleFunc("/people/{id}", DeletePersonEndpoint).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":3000", router))
}
