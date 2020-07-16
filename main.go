package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"fmt"
	_ "encoding/json"
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
//returning data from sqlite database with a get request to an  endpoint where some of the values have a null value

type Employee struct {
	ID int `json:"id"`
	firstName string `json:"first_name"`
	lastName string `json:"last_name"`
	dateOfBirth string `json:"date_of_birth"`
	addressLineOne string `json:"address_line_one"`
	addressLineTwo sql.NullString `json:"address_line_two"`
	city string `json:"city"`
	postcode string `json:"postcode"`
	startDate string `json:"start_date"`
	nextOfKin string `json:"next_of_kin"`
	position string `json:"position"`
	endDate sql.NullString `json:"end_date"`
	recordCreatedDate string `json:"record_created_date"`
}

func main() {
	http.HandleFunc("/employees", employeeHandler)
	//TO DO HOW TO MAKE /123 DYNAMIC
	http.HandleFunc("/employees/1", employeeByIDHandler)
	http.HandleFunc("/employees/search/text", employeeByTextSearchHandler)
	http.ListenAndServe(":4000", nil)
}

/*var database *sql.DB

func init() {
	var err error
	database, err := sql.Open("sqlite3", "employee-database.db")
if err != nil {
		log.Fatal(err)
	}

	if err = database.Ping(); err != nil {
		log.Fatal(err)
	}
} */

func employeeHandler(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case "GET":
	database, err := sql.Open("sqlite3", "employee-database.db")
		if err != nil {
			log.Fatal(err)
		}
		rows, err := database.Query("SELECT ID, firstName, lastName, dateOfBirth, addressLineOne, addressLineTwo, city, postcode, startDate, nextOfKin, position, endDate, recordCreatedDate FROM employees")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		employees := make([]*Employee, 0)
		for rows.Next() {
			employee := new(Employee)
			err := rows.Scan(&employee.ID, &employee.firstName, &employee.lastName, &employee.dateOfBirth, &employee.addressLineOne, &employee.addressLineTwo, &employee.city, &employee.postcode, &employee.startDate, &employee.nextOfKin, &employee.position, &employee.endDate, &employee.recordCreatedDate)
			if err != nil {
				log.Fatal(err)
			}
			employees = append(employees, employee)
		}
		if err = rows.Err(); err != nil {
			log.Fatal(err)
		}
		for _, employee := range employees {

			fmt.Fprintf(w, "%b, %s %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s\n", employee.ID, employee.firstName, employee.lastName, employee.dateOfBirth, employee.addressLineOne, employee.addressLineTwo, employee.city, employee.postcode, employee.startDate, employee.nextOfKin, employee.position, employee.endDate, employee.recordCreatedDate)
		}
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
			database, err := sql.Open("sqlite3", "employee-database.db")
			if err != nil {
				log.Fatal(err)
			}
			// TODO: pick up id from url path
			id := 1
			row := database.QueryRow("SELECT ID, firstName, lastName, dateOfBirth, addressLineOne, addressLineTwo, city, postcode, startDate, nextOfKin, position, endDate, recordCreatedDate FROM employees WHERE ID =$1", id)
			employee := new(Employee)
			err = row.Scan(&employee.ID, &employee.firstName, &employee.lastName, &employee.dateOfBirth, &employee.addressLineOne, &employee.addressLineTwo, &employee.city, &employee.postcode, &employee.startDate, &employee.nextOfKin, &employee.position, &employee.endDate, &employee.recordCreatedDate)
			if err == sql.ErrNoRows{
				http.NotFound(w, r)
				return
			} else if err != nil{
				fmt.Fprint(w, err)
				http.Error(w, http.StatusText(500), 500)
				return
			}
			fmt.Fprintf(w, "%d, %s %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s\n", employee.ID, employee.firstName, employee.lastName, employee.dateOfBirth, employee.addressLineOne, employee.addressLineTwo, employee.city, employee.postcode, employee.startDate, employee.nextOfKin, employee.position, employee.endDate, employee.recordCreatedDate)
			fmt.Fprintln(w, "Endpoint hit: Return an employees by ID")

		case "PATCH":
			fmt.Fprintln(w, "Endpoint hit: Update employees record by ID")
		case "DELETE":
			fmt.Fprintln(w, "Endpoint hit: Delete employees record by ID")
		default:
			fmt.Fprintln(w, "Endpoint hit: This endpoint only supports GET, PATCH and DELETE requests by ID")
		}
	}

func employeeByTextSearchHandler(w http.ResponseWriter, r*http.Request){
	//SELECT* FROM employees WHERE firstName like '%AAA%' OR lastName like '%AAA%' OR position like '%AAA% OR startDate like '%1111%' OR endDate like '%1111%' OR recordCreatedDate like '%1111%'
	switch method := r.Method; method {
	case "GET":
		if r.URL.Query().Get("name") != "" || r.URL.Query().Get ("position") != "" || r.URL.Query().Get("startDate") != "" || r.URL.Query().Get("endDate") != ""  {
			fmt.Fprintln(w, "Endpoint hit: Search for employees by name, position or date")
		} else {
			fmt.Fprintln(w, "Endpoint hit: Please enter name, position, start or end date to search for employees")
		}
	default:
	fmt.Fprintln(w, "Endpoint hit: This endpoint only supports GET requests")

	}
}





/*var returnAllEmployees = func (){
	database, err := sql.Open("sqlite3", "employee-database.db")
	if err != nil {
		log.Fatal(err)
	}
	rows, err := database.Query("SELECT ID, firstName, lastName, dateOfBirth, addressLineOne, addressLineTwo, city, postcode, startDate, nextOfKin, position, endDate, recordCreatedDate FROM employees")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	employees := make([]*Employee, 0)
	for rows.Next() {
		employee := new(Employee)
		err := rows.Scan(&employee.ID, &employee.firstName, &employee.lastName, &employee.dateOfBirth, &employee.addressLineOne, &employee.addressLineTwo, &employee.city, &employee.postcode, &employee.startDate, &employee.nextOfKin, &employee.position, &employee.endDate, &employee.recordCreatedDate)
		if err != nil {
			log.Fatal(err)
		}
		employees = append(employees, employee)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	for _, employee := range employees {
		//fmt.Fprintf(w, "%b, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s\n", employee.ID, &employee.firstName, &employee.lastName, &employee.dateOfBirth, &employee.addressLineOne, &employee.city, &employee.postcode, &employee.startDate, &employee.nextOfKin, &employee.position, &employee.recordCreatedDate)
		fmt.Println(employee.ID,":" + " " + employee.firstName, employee.lastName, employee.dateOfBirth, employee.addressLineOne, employee.addressLineTwo, employee.city, employee.postcode, employee.startDate, employee.nextOfKin, employee.position, employee.endDate, employee.recordCreatedDate)
	}
} */




/* TO DO
Test by a) returning all employees (employeeHandler - get request) b) creating a new employee (employeeHandler - post request)
*/



//CREATE TABLE "employees" ( "ID" INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE, "firstName" TEXT NOT NULL, "lastName" TEXT NOT NULL, "dateOfBirth" TEXT NOT NULL, "addressLineOne" TEXT NOT NULL, "addressLineTwo" TEXT, "city" TEXT NOT NULL, "postcode" TEXT NOT NULL, "startDate" TEXT NOT NULL, "nextOfKin" TEXT NOT NULL, "position" TEXT NOT NULL, "endDate" INTEGER, "recordCreatedDate" TEXT DEFAULT CURRENT_TIMESTAMP)

/*INSERT INTO employees (firstName, lastName, dateOfBirth, addressLineOne, city, postcode, startDate, nextOfKin, position)
VALUES ('Salman','Ahmed','1999-01-01 00:00:00:000','33 Holborn', 'London', 'EC1N 2HT','2020-01-09 00:00:00.000', 'Steve Stotter','CEO');
 */

/*INSERT INTO employees (firstName, lastName, dateOfBirth, addressLineOne, city, postcode, startDate, nextOfKin, position)
VALUES ('Joe','Jenkins','1992-03-31 00:00:00:000','12 Little Tree Lane', 'London', 'E7 0TR','2020-03-10 00:00:00.000', 'Laura Jenkins','Head of Design');
 */