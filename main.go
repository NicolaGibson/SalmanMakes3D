package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"strconv"
	//"time"
)
/*


	/employee/ POST create
	/employee/123 PUT / PATCH update
	/employee/ GET return all
	/employee/123213 GET return specific user
	/employee/1232 DELETE
	/employee/?startDate=13213&endDate=123213&name=nikki GET return filtered
	/employee/?search=nikki GET return filtered

 localhost:4000/employees/search/?name=Joe
 */


func main() {
	http.HandleFunc("/employees", employeeHandler)
	//TO DO HOW TO MAKE /123 DYNAMIC
	http.HandleFunc("/employees/1", employeeByIDHandler)
	http.HandleFunc("/employees/search/text", employeeByTextSearchHandler)
	http.HandleFunc("/employees/search/date", employeeByDateHandler)
	http.ListenAndServe(":4000", nil)
}

func employeeHandler(w http.ResponseWriter, r *http.Request) {
	//	SELECT ID, firstName, lastName, dateOfBirth, addressLineOne, city, postcode, startDate, nextOfKin, position, recordCreatedDate FROM employees
	switch method := r.Method; method {
	case "GET":
		returnAllEmployees()
			fmt.Fprintln(w, "Endpoint hit: Return all employees records")


	case "POST":
		fmt.Fprintln(w, "Endpoint hit: Create new employee record")
	default:
		fmt.Fprintln(w, "Endpoint hit: This endpoint only supports GET or POST requests")

	}
}

func employeeByIDHandler(w http.ResponseWriter, r *http.Request) {
	//	SELECT ID, firstName, lastName, dateOfBirth, addressLineOne, city, postcode, startDate, nextOfKin, position, recordCreatedDate FROM employees WHERE ID = 1
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
	//SELECT* FROM employees WHERE firstName like '%AAA%' OR lastName like '%AAA%' OR position like '%AAA%'
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
	//SELECT* FROM employees WHERE startDate like '%1111%' OR endDate like '%1111%' OR recordCreatedDate like '%1111%'
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

func returnAllEmployees(){
	database, _ := sql.Open("sqlite3", "employee-database.db")
	rows, _ := database.Query("SELECT ID, firstName, lastName, dateOfBirth, addressLineOne, city, postcode, startDate, nextOfKin, position, recordCreatedDate FROM employees")
	var ID int
	var firstName string
	var lastName string
	var dateOfBirth string
	var addressLineOne string
	//var addressLineTwo string
	var city string
	var postcode string
	var startDate string
	var nextOfKin string
	var position string
	var recordCreatedDate string
	for rows.Next(){
		rows.Scan(&ID, &firstName, &lastName, &dateOfBirth, &addressLineOne,  &city, &postcode, &startDate, &nextOfKin, &position, &recordCreatedDate)
		fmt.Println(strconv.Itoa(ID) + ": " + firstName + " " + lastName + " " + addressLineOne +  " " + city + " " + postcode + " " + startDate + " " + nextOfKin + " " + position + " " + recordCreatedDate)
	}
	rows.Close()
}


/* TO DO
Connect employee database to go project
Create table if does not exist
Test by a) returning all employees (employeeHandler - get request) b) creating a new employee (employeeHandler - post request)
*/

/* TO DO Oge

Create endpoints, add in date created field, figure out SQL queries, write notes on decision to use type Interger for date.
Split search endpoint - search by text and search by date.

*/

//CREATE TABLE "employees" ( "ID" INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE, "firstName" TEXT NOT NULL, "lastName" TEXT NOT NULL, "dateOfBirth" TEXT NOT NULL, "addressLineOne" TEXT NOT NULL, "addressLineTwo" TEXT, "city" TEXT NOT NULL, "postcode" TEXT NOT NULL, "startDate" TEXT NOT NULL, "nextOfKin" TEXT NOT NULL, "position" TEXT NOT NULL, "endDate" INTEGER, "recordCreatedDate" TEXT DEFAULT CURRENT_TIMESTAMP)

/*INSERT INTO employees (firstName, lastName, dateOfBirth, addressLineOne, city, postcode, startDate, nextOfKin, position)
VALUES ('Salman','Ahmed','1999-01-01 00:00:00:000','33 Holborn', 'London', 'EC1N 2HT','2020-01-09 00:00:00.000', 'Steve Stotter','CEO');
 */

/*INSERT INTO employees (firstName, lastName, dateOfBirth, addressLineOne, city, postcode, startDate, nextOfKin, position)
VALUES ('Joe','Jenkins','1992-03-31 00:00:00:000','12 Little Tree Lane', 'London', 'E7 0TR','2020-03-10 00:00:00.000', 'Laura Jenkins','Head of Design');
 */