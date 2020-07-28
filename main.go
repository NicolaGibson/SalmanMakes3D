package main

import (
	"database/sql"
	//"encoding/json"
	//_ "encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
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

 localhost:4000/employees/search/?name=Joe
 */
//returning data from sqlite database with a get request to an  endpoint where some of the values have a null value

//var db *sqlite3.db

type Employee struct {
	ID int `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	DateOfBirth string `json:"date_of_birth"`
	AddressLineOne string `json:"address_line_one"`
	AddressLineTwo sql.NullString `json:"address_line_two"`
	City string `json:"city"`
	Postcode string `json:"postcode"`
	StartDate string `json:"start_date"`
	NextOfKin string `json:"next_of_kin"`
	Position string `json:"position"`
	EndDate sql.NullString `json:"end_date"`
	RecordCreatedDate string `json:"record_created_date"`
}

func main() {
	http.HandleFunc("/employees", employeeHandler)
	//TO DO HOW TO MAKE /123 DYNAMIC
	http.HandleFunc("/employees/show", employeeByIDHandler)
	http.HandleFunc("/employees/search/text", employeeByTextSearchHandler)
	http.ListenAndServe(":4000", nil)
}

func employeeHandler(w http.ResponseWriter, r *http.Request) {

	switch method := r.Method; method {
	case "GET":
	database, err := sql.Open("sqlite3", "employee.db")
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
			err := rows.Scan(&employee.ID, &employee.FirstName, &employee.LastName, &employee.DateOfBirth, &employee.AddressLineOne, &employee.AddressLineTwo, &employee.City, &employee.Postcode, &employee.StartDate, &employee.NextOfKin, &employee.Position, &employee.EndDate, &employee.RecordCreatedDate)
			if err != nil {
				log.Fatal(err)
			}
			employees = append(employees, employee)

		}
		if err = rows.Err(); err != nil {
			log.Fatal(err)
		}
		for _, employee := range employees {
			//format, _ := json.Marshal(employees)
			//fmt.Fprintf(w,string(format))
			fmt.Fprintf(w, "%b, %s %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s\n", employee.ID, employee.FirstName, employee.LastName, employee.DateOfBirth, employee.AddressLineOne, employee.AddressLineTwo, employee.City, employee.Postcode, employee.StartDate, employee.NextOfKin, employee.Position, employee.EndDate, employee.RecordCreatedDate)
		}
		_, _ = fmt.Fprintln(w, "Endpoint hit: Return all employees records")
	case "POST":
		database, err := sql.Open("sqlite3", "employee.db")
		if err != nil {
			log.Fatal(err)
		}

		firstName := r.FormValue("firstName")
		lastName := r.FormValue ("lastName")
		if firstName == "" || lastName == "" {
			http.Error(w, http.StatusText(400), 400)
			return
		}
		result, err := database.Exec("INSERT INTO employees VALUES ($1, $2)", firstName, lastName)
		if err != nil{
			http.Error(w, http.StatusText(500), 500)
			return
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			http.Error(w, http.StatusText(500),500)
			return
		}

		fmt.Fprintf(w, "Employee %s created successfully (%d row affected)\n", firstName, rowsAffected)

		fmt.Fprintln(w, "Endpoint hit: Create new employee record")
	default:
		fmt.Fprintln(w, "Endpoint hit: This endpoint only supports GET or POST requests")

	}
}

func employeeByIDHandler(w http.ResponseWriter, r *http.Request) {
	// http://localhost:4000/employees/show?ID=2
	switch method := r.Method; method {
		case "GET":
			database, err := sql.Open("sqlite3", "employee.db")
			if err != nil {
				log.Fatal(err)
			}

			ID := r.FormValue("ID")
			if ID == "" {
				http.Error(w, http.StatusText(400), 400)
				return
			}
			row := database.QueryRow("SELECT ID, firstName, lastName, dateOfBirth, addressLineOne, addressLineTwo, city, postcode, startDate, nextOfKin, position, endDate, recordCreatedDate FROM employees WHERE ID =$1",ID)
			employee := new(Employee)
			//scan input text, reads from there and stores space seperated values in successive arguements
			err = row.Scan(&employee.ID, &employee.FirstName, &employee.LastName, &employee.DateOfBirth, &employee.AddressLineOne, &employee.AddressLineTwo, &employee.City, &employee.Postcode, &employee.StartDate, &employee.NextOfKin, &employee.Position, &employee.EndDate, &employee.RecordCreatedDate)

			if err == sql.ErrNoRows{
				http.NotFound(w, r)
				return
			} else if err != nil{
				fmt.Fprint(w, err)
				http.Error(w, http.StatusText(500), 500)
				return
			}
			fmt.Fprintf(w, "%d, %s %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s\n", employee.ID, employee.FirstName, employee.LastName, employee.DateOfBirth, employee.AddressLineOne, employee.AddressLineTwo, employee.City, employee.Postcode, employee.StartDate, employee.NextOfKin, employee.Position, employee.EndDate, employee.RecordCreatedDate)
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