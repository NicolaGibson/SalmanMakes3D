package main

import (
	"fmt"
	"net/http"
)

func main(){
	http.HandleFunc("/employee/create", createEmployee)
	http.HandleFunc("/employee/update", updateEmployeeDetails)
	http.HandleFunc ("/employee/delete", deleteEmployee)
	http.HandleFunc("/employee/search/all", getAllEmployees)
	http.HandleFunc("/employee/search/ID", getDetailsByID)
	http.HandleFunc("/employee/search/forename", getDetailsByForename)
	http.HandleFunc("/employee/search/surname", getDetailsBySurname)
	http.HandleFunc("/employee/search/position", getDetailsByPosition)
	http.HandleFunc("/employee/search/start-date", getDetailsByStartDate)
	http.HandleFunc("/employee/search/end-date", getDetailsByEndDate)
	http.ListenAndServe(":4000", nil)
}

func createEmployee(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w,"Endpoint hit: Create new employee record")
}

func updateEmployeeDetails(w http.ResponseWriter, r *http.Request ){
	fmt.Fprintln(w,"Endpoint hit: Update employee details")
}

func deleteEmployee(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "Endpoint hit: Delete employee")
}

func getAllEmployees(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "Endpoint hit: Return all employees")
}
func getDetailsByID(w http.ResponseWriter, r *http.Request){
 fmt.Fprintln(w, "Endpoint hit: Search for employee by ID")
}

func getDetailsByForename(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "Endpoint hit: Search for employee by forename")
}

func getDetailsBySurname(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "Endpoint hit: Search for employee by surname")
}

func getDetailsByPosition(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "Endpoint hit: Search for employee by position")
}

func getDetailsByStartDate(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "Endpoint hit: Search for employee by start date")
}

func getDetailsByEndDate(w http.ResponseWriter, r *http.Request){
	fmt.Fprint(w, "Endpoint hit: Search for employee by end date")
}
