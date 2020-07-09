package main

import (
	"fmt"
	"net/http"
)




func main(){
	http.HandleFunc("/create-employee", createEmployee)
	http.HandleFunc("/search-ID", getDetailsByID)
	http.ListenAndServe(":4000", nil)
}

func createEmployee(w http.ResponseWriter, r *http.Request){
	fmt.Println("Endpoint hit: Create new employee record")
}

func getDetailsByID(w http.ResponseWriter, r *http.Request){
 fmt.Println("Endpoint hit: Search Details By ID.")



}