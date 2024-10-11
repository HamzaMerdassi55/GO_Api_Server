package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	err := a.Initialise(DbUser, DBPassword, "test")
	if err != nil {
		log.Fatal("Error occurred while initializing the database")
	}
	createTable()
	m.Run()
	
}

func createTable() {
	createTableQuery := `CREATE TABLE IF NOT EXISTS products (
        id int NOT NULL AUTO_INCREMENT,
        name varchar(255) NOT NULL,
        quantity int,
        price DECIMAL(10,1),  -- Change here
		PRIMARY KEY (id)
    );`
	_, err := a.DB.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}


func clearTable() {
	a.DB.Exec("DELETE FROM products")
	a.DB.Exec("ALTER table products AUTO_INCREMENT = 1")
	log.Println("clearTable")
}

func addProduct(name string, quantity int, price float64) {
	query := fmt.Sprintf("INSERT INTO products(name, quantity, price) VALUES('%v', %v, %v)", name, quantity, price)
	_, err := a.DB.Exec(query)
	if err != nil {
		log.Println(err)
	}
}

func TestGetProduct(t *testing.T) {
	clearTable()
	addProduct("Keyboard", 100, 5000)
	request, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(request)  
	checkStatusCode(t, http.StatusOK, response.Code) // Vérifiez le statut attendu
}

func checkStatusCode(t *testing.T, expectedStatusCode int, actualStatusCode int) {
	if expectedStatusCode != actualStatusCode {
		t.Errorf("Expected status: %v, Received: %v", expectedStatusCode, actualStatusCode)
	}
}

func sendRequest(request *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	a.Router.ServeHTTP(recorder, request)
	return recorder
}

func TestCreateProduct(t *testing.T) {
	clearTable()
	var product = []byte(`{"name":"chair", "quantity":1, "price":100}`)
	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(product))
	req.Header.Set("Content-Type", "application/json")

	response := sendRequest(req)
	checkStatusCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "chair" {
		t.Errorf("Expected name: %v, Got: %v", "chair", m["name"])
	}
	log.Printf("%T", m["quantity"])
	if m["quantity"] != 1.0 {
		t.Errorf("Expected quantity: %v, Got: %v", 1.0, m["quantity"])
	}
}

func TestDeleteProduct(t *testing.T){

	clearTable()
	addProduct("connector", 10, 10)

	req , _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(req)
	checkStatusCode(t, http.StatusOK, response.Code)

	req , _ = http.NewRequest("DELETE", "/product/1", nil)
	response = sendRequest(req)
	checkStatusCode(t, http.StatusOK, response.Code)

	req , _ = http.NewRequest("GET", "/product/1", nil)
	response = sendRequest(req)
	checkStatusCode(t, http.StatusNotFound, response.Code)
}

func TestUpdateProduct(t *testing.T){
	clearTable()
	addProduct("connector", 10, 10)
    req, _ := http.NewRequest("GET", "/product/1", nil)
    response := sendRequest(req)

    var oldValue map[string]interface{}
    json.Unmarshal(response.Body.Bytes(), &oldValue)

    var product = []byte(`{"name":"connector", "quantity":1, "price":10}`)
	req, _ = http.NewRequest("PUT", "/product/1", bytes.NewBuffer(product))
	req.Header.Set("Content-Type", "application/json")

	response = sendRequest(req)
	var newValue map[string]interface{}
    json.Unmarshal(response.Body.Bytes(), &newValue)

	if oldValue["id"] != newValue["id"] {
        t.Errorf("Expected id: %v, Got: %v", newValue["id"], oldValue["id"])
    }
}