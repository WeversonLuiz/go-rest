package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"github.com/gorilla/mux"
)

type conta struct {
	ID               string		`json:"id"`
	Numero      	 int        `json:"numero"`
	Saldo       	 float64    `json:"saldo"`
	DataAbertura 	 time.Time  `json:"dataAbertura"`
	Status           bool       `json:"bloqueada"`
}

type allContas []conta


var contas = allContas{
	{
		ID:					"1",
		Numero:       		1,
		Saldo:        		1000.00,
		DataAbertura:       time.Now(),
		Status:             false,
	},
}

func criarConta(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var newConta conta
	reqBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Fprintf(response, "Informe dados para a criação da conta.")
	}
	json.Unmarshal(reqBody, &newConta)
	contas = append(contas, newConta)

	response.WriteHeader(http.StatusCreated)

	json.NewEncoder(response).Encode(newConta)
}

func AllContas(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	json.NewEncoder(response).Encode(contas)
}

func bloquearConta(response http.ResponseWriter, request *http.Request)    {
	response.Header().Set("content-type", "application/json")
	contaID := mux.Vars(request)["id"]
	var bloquear conta
	reqBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Fprintf(response, "Informe valor do saque.")
	}
	json.Unmarshal(reqBody, &bloquear)
	for i, singleConta := range contas {
		if singleConta.ID == contaID {
			singleConta.Status = true
			contas = append(contas[:i], singleConta)
			json.NewEncoder(response).Encode(singleConta)
		}
	}
}

func desbloquearConta(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	contaID := mux.Vars(request)["id"]
	var desbloquear conta
	reqBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Fprintf(response, "Informe valor do saque.")
	}
	json.Unmarshal(reqBody, &desbloquear)
	for i, singleConta := range contas {
		if singleConta.ID == contaID {
			singleConta.Status = false
			contas = append(contas[:i], singleConta)
			json.NewEncoder(response).Encode(singleConta)
		}
	}
}

func creditarEmConta(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	contaID := mux.Vars(request)["id"]
	var creditar conta
	reqBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Fprintf(response, "Informe valor do deposito.")
	}
	json.Unmarshal(reqBody, &creditar)
	for i, singleConta := range contas {
		if singleConta.Status == true{
	
		} else if singleConta.ID == contaID {
			singleConta.Saldo = singleConta.Saldo + creditar.Saldo
			contas = append(contas[:i], singleConta)
			json.NewEncoder(response).Encode(singleConta)
		}
	}
}

func debitarEmConta(response http.ResponseWriter, request *http.Request)     {
	response.Header().Set("content-type", "application/json")
	contaID := mux.Vars(request)["id"]
	var debitar conta
	reqBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Fprintf(response, "Informe valor do saque.")
	}
	json.Unmarshal(reqBody, &debitar)
	for i, singleConta := range contas {
		if singleConta.Status == true {
		
		} else if singleConta.ID == contaID {
			singleConta.Saldo = singleConta.Saldo - debitar.Saldo
			contas = append(contas[:i], singleConta)
			json.NewEncoder(response).Encode(singleConta)
		}
	}
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Seja bem vindo ao Nosso Banco")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/conta", criarConta).Methods("POST")
	router.HandleFunc("/contas/creditar/{id}", creditarEmConta).Methods("PATCH")
	router.HandleFunc("/contas/debitar/{id}", debitarEmConta).Methods("PATCH")
	router.HandleFunc("/contas/desbloquear/{id}", desbloquearConta).Methods("PATCH")
	router.HandleFunc("/contas/bloquear/{id}", bloquearConta).Methods("PATCH")
	router.HandleFunc("/contas", AllContas).Methods("GET")
	
	fmt.Println("Servidor disponível na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

