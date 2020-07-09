package main

import (
	"fmt"
	"net/http"
)




func main(){
	http.HandleFunc("/employee/create", createEmployee)
	http.HandleFunc("/employee/update", updateEmployeeDetails)
	http.HandleFunc ("/employee/delete", deleteEmployee)
	http.HandleFunc("/search-ID", getDetailsByID)
	http.ListenAndServe(":4000", nil)
}

func createEmployee(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w,"Endpoint hit: Create new employee record")
}

func updateEmployeeDetails(w http.ResponseWriter, r *http.Request ){
	fmt.Fprintln(w,"Endpoint hit: Update employee details.")
}

func deleteEmployee(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "Endpoint hit: Delete Employee ")
}

func getDetailsByID(w http.ResponseWriter, r *http.Request){
 fmt.Fprintln(w, "Endpoint hit: Search Details By ID.")
}