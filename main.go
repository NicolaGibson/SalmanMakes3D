package main

import (
	"fmt"
	"net/http"
)
/*


	/employee/ POST create
	/employee/123 PUT / PATCH update
	/employee/ GET return all
	/employee/123213 GET return specific user
	/employee/1232 DELETE
	/employee/?startDate=13213&endDate=123213&name=nikki GET return filtered
	/employee/?search=nikki GET return filtered

 localhost:4000/employees/search/?name=Rob
 */

func main(){
	http.HandleFunc("/employees", employeeHandler)
	//TO DO HOW TO MAKE /123 DYNAMIC
	http.HandleFunc("/employees/123", employeeByIDHandler)
	http.HandleFunc("/employees/search/text", employeeByTextSearchHandler)
	http.HandleFunc("/employees/search/date", employeeByDateHandler)
	http.ListenAndServe(":4000", nil)
}

func employeeHandler(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case "GET":
			fmt.Fprintln(w, "Endpoint hit: Return all employees records")
	case "POST":
		fmt.Fprintln(w, "Endpoint hit: Create new employee record")
	default:
		fmt.Fprintln(w, "Endpoint hit: This endpoint only supports GET or POST requests")

	}
}

func employeeByIDHandler(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case "GET":
		fmt.Fprintln(w, "Endpoint hit: Return all employees by ID")
	case "PATCH":
		fmt.Fprintln(w, "Endpoint hit: Update employees record by ID")
	case "DELETE":
		fmt.Fprintln(w, "Endpoint hit: Delete employees record by ID")
	default:
		fmt.Fprintln(w, "Endpoint hit: This endpoint only supports GET, PATCH and DELETE requests by ID")
	}
}

func employeeByTextSearchHandler(w http.ResponseWriter, r*http.Request){
	switch method := r.Method; method {
	case "GET":
		if r.URL.Query().Get("name") != "" || r.URL.Query().Get ("position") != "" {
			fmt.Fprintln(w, "Endpoint hit: Search for employees by name or position")
		} else {
			fmt.Fprintln(w, "Endpoint hit: Please enter name or position to search for employees")
		}
	default:
	fmt.Fprintln(w, "Endpoint hit: This endpoint only supports GET requests")
	}
}

func employeeByDateHandler(w http.ResponseWriter, r*http.Request) {
	switch method := r.Method; method {
	case "GET":
		if r.URL.Query().Get("startDate") != "" || r.URL.Query().Get("endDate") != "" {
			fmt.Fprintln(w, "Endpoint hit: Search for employees by end date or start date")
		} else {
			fmt.Fprintln(w, "Endpoint hit: Please enter end date or start date to search for employees")
		}
	default:
		fmt.Fprintln(w, "Endpoint hit: This endpoint only supports GET requests")
	}
}


/*func createEmployeesTable(){
	//CREATE TABLE "employees" ( "ID" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE, "firstName" TEXT NOT NULL, "lastName" TEXT NOT NULL, "dateOfBirth" INTEGER NOT NULL, "addressLineOne" TEXT NOT NULL, "addressLineTwo" TEXT, "city" TEXT NOT NULL, "postcode" TEXT NOT NULL, "startDate" INTEGER NOT NULL, "nextOfKin" TEXT NOT NULL, "position" TEXT NOT NULL, "endDate" INTEGER )
}
func addEmployeeToTable(){
	/*INSERT INTO employees (ID, firstName, lastName, dateOfBirth, addressLineOne, city, postcode, startDate, nextOfKin, position)
	VALUES ('1','Salman','Ahmed','1999-01-01 00:00:00:000','33 Holborn', 'London', 'EC1N 2HT','2020-01-09 00:00:00.000', 'Steve Stotter','CEO' );
} */


