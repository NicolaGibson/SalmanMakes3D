package main

import (
	"database/sql"
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

var ID = ""
var param = ""

type Employee struct {
	ID                int            `json:"id"`
	FirstName         string         `json:"first_name"`
	LastName          string         `json:"last_name"`
	DateOfBirth       string         `json:"date_of_birth"`
	AddressLineOne    string         `json:"address_line_one"`
	AddressLineTwo    sql.NullString `json:"address_line_two"`
	City              string         `json:"city"`
	Postcode          string         `json:"postcode"`
	StartDate         string         `json:"start_date"`
	NextOfKin         string         `json:"next_of_kin"`
	Position          string         `json:"position"`
	EndDate           sql.NullString `json:"end_date"`
	RecordCreatedDate string         `json:"record_created_date"`
	employeeStatus    string
}

var db *sql.DB

func init(){
	var err error
	db, err = sql.Open("sqlite3", "employee.db")
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}


func main() {
	http.HandleFunc("/employees", employeeHandler)
	http.HandleFunc("/employees/", employeeByIDHandler)
	//http.HandleFunc("/employees/?", employeeSearchHandler)
	http.ListenAndServe(":4000", nil)


}

//code for next page
func employeeHandler(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case "GET":
		rows, err := db.Query("SELECT ID, firstName, lastName, dateOfBirth, addressLineOne, addressLineTwo, city, postcode, startDate, nextOfKin, position, endDate, recordCreatedDate FROM employees ORDER BY lastName ASC LIMIT 0, 50")
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
			fmt.Fprintf(w, "%b, %s %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s\n", employee.ID, employee.LastName, employee.FirstName, employee.DateOfBirth, employee.AddressLineOne, employee.AddressLineTwo, employee.City, employee.Postcode, employee.StartDate, employee.NextOfKin, employee.Position, employee.EndDate, employee.RecordCreatedDate)
		}
		fmt.Fprintln(w, "Endpoint hit: Return all employees records")
	case "POST":
		firstName := r.FormValue("firstName")
		lastName := r.FormValue("lastName")
		dateOfBirth := r.FormValue("dateOfBirth")
		addressLineOne := r.FormValue("addressLineOne")
		addressLineTwo := r.FormValue("addressLineTwo")
		city := r.FormValue("city")
		postcode := r.FormValue("postcode")
		startDate := r.FormValue("startDate")
		nextOfKin := r.FormValue("nextOfKin")
		position := r.FormValue("position")
		employeeStatus := r.FormValue("employeeStatus")

		if firstName == "" ||lastName == "" ||addressLineOne == "" ||city == ""||postcode == "" ||startDate == "" ||nextOfKin == "" ||position == "" ||employeeStatus == "" {
			http.Error(w, http.StatusText(400), 400)
			return
		}
		result, err := db.Exec("INSERT INTO employees (firstName, lastName, dateOfBirth, addressLineOne, addressLineTwo, city, postcode, startDate, nextOfKin, position, employeeStatus) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)", firstName, lastName, dateOfBirth, addressLineOne, addressLineTwo, city, postcode, startDate, nextOfKin, position, employeeStatus)
		if err != nil{
			http.Error(w, http.StatusText(500), 500)
			return
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			http.Error(w, http.StatusText(500),500)
			return
		}

		fmt.Fprintf(w, "Employee %s created successfully (%d row affected)\n", firstName + " " + lastName, rowsAffected)

	default:
		fmt.Fprintln(w, "Endpoint hit: This endpoint only supports GET or POST requests")

	}
}

func getIDParam(path string) (ps string){
	//ignore first / when hitting second / return everything after as a parameter
	for i := 1; i<len(path); i++ {
		if path[i] == '/'{
			ps = path [i +1:]
		}
	}
	ID = ps
	return
}


func employeeByIDHandler(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
		case "GET":
			getIDParam(r.URL.Path)

			if ID == ""{
				http.Error(w, http.StatusText(400), 400)
				return
			}
			row := db.QueryRow("SELECT ID, firstName, lastName, dateOfBirth, addressLineOne, addressLineTwo, city, postcode, startDate, nextOfKin, position, endDate, recordCreatedDate, employeeStatus FROM employees WHERE ID =$1", ID)
			employee := new(Employee)
			//scan input text, reads from there and stores space seperated values in successive arguements
			err := row.Scan(&employee.ID, &employee.FirstName, &employee.LastName, &employee.DateOfBirth, &employee.AddressLineOne, &employee.AddressLineTwo, &employee.City, &employee.Postcode, &employee.StartDate, &employee.NextOfKin, &employee.Position, &employee.EndDate, &employee.RecordCreatedDate, &employee.employeeStatus)
			if err == sql.ErrNoRows{
				http.NotFound(w, r)
				return
			} else if err != nil{
				fmt.Fprint(w, err)
				http.Error(w, http.StatusText(500), 500)
				return
			}
			fmt.Fprintf(w, "%d, %s %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s\n", employee.ID, employee.FirstName, employee.LastName, employee.DateOfBirth, employee.AddressLineOne, employee.AddressLineTwo, employee.City, employee.Postcode, employee.StartDate, employee.NextOfKin, employee.Position, employee.EndDate, employee.RecordCreatedDate, employee.employeeStatus)
			fmt.Fprint(w, "Endpoint hit: Return an employee by ID")

		case "PATCH":

			//UPDATE employees SET position = "Design Manager" WHERE ID = 2
			fmt.Fprintln(w, "Endpoint hit: Update employee record by ID")
		case "DELETE":
			getIDParam(r.URL.Path)

			if ID == ""{
				http.Error(w, http.StatusText(400), 400)
				return
			}
			row := db.QueryRow("UPDATE employees SET employeeStatus ='Disabled' WHERE ID =$1", ID)

			employee := new(Employee)
			err := row.Scan(&employee.ID, &employee.FirstName, &employee.LastName, &employee.DateOfBirth, &employee.AddressLineOne, &employee.AddressLineTwo, &employee.City, &employee.Postcode, &employee.StartDate, &employee.NextOfKin, &employee.Position, &employee.EndDate, &employee.RecordCreatedDate, &employee.employeeStatus)
			if err == sql.ErrNoRows {
				http.NotFound(w, r)
				return
			}else if err != nil {
					fmt.Print(w, err)
					http.Error(w, http.StatusText(500), 500)
					return
				}

				fmt.Fprintf(w, "Employee ID %d, %s, %s deleted successfully \n", employee.ID, employee.FirstName+" "+employee.LastName)
				fmt.Fprintln(w, "Endpoint hit: Delete employee record by ID")
		default:
			fmt.Fprintln(w, "Endpoint hit: This endpoint only supports GET, PATCH and DELETE requests by ID")
		}
	}


/*func getfirstParam(path string) (ps string){
	//ignore first / when hitting second / return everything after as a parameter
	for i := 1; i<len(path); i++ {
		if path[i] == '?'{
			for j := 0; j <len(path); j++{
				ps = j
			}

		}

	}
	param = ps

	fmt.Println(param)
	return
} */
func employeeSearchHandler(w http.ResponseWriter, r*http.Request){
	/* To Do: return employee record by name using current implementation err no such column name. Rewrite implementation using map for search by lastName, position, startDate, endDate, use map to pull out
	different values. change url path, needs to be part of employee. Search on the employees endpoint
	SELECT* FROM employees WHERE firstName like 'A%' OR lastName like '%AAA%' OR position like '%AAA% OR startDate like '%1111%' OR endDate like '%1111%' OR recordCreatedDate like '%1111%'
	firstName := r.FormValue("firstName")  - OR firstName like '$1%'"
	*/

	switch method := r.Method; method {
	case "GET":
		//getfirstParam(r.URL.Path)
		employees := make([]*Employee,0)
		rows, err := db.Query("SELECT ID, firstName, lastName, dateOfBirth, addressLineOne, addressLineTwo, city, postcode, startDate, nextOfKin, position, endDate, recordCreatedDate FROM employees WHERE firstName LIKE param")
		if err != nil{
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			employee := new(Employee)
			//scan input text, reads from there and stores space seperated values in successive arguements
			if err = rows.Scan(&employee.ID, &employee.FirstName, &employee.LastName, &employee.DateOfBirth, &employee.AddressLineOne, &employee.AddressLineTwo, &employee.City, &employee.Postcode, &employee.StartDate, &employee.NextOfKin, &employee.Position, &employee.EndDate, &employee.RecordCreatedDate); err != nil{
			panic(err)
				}
				employees = append(employees, employee)
				fmt.Fprintf(w, "%d, %s %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s\n", employee.ID, employee.FirstName, employee.LastName, employee.DateOfBirth, employee.AddressLineOne, employee.AddressLineTwo, employee.City, employee.Postcode, employee.StartDate, employee.NextOfKin, employee.Position, employee.EndDate, employee.RecordCreatedDate)
			}
			/*if r.URL.Query().Get("firstName") != "" {
			fmt.Fprintln(w, "Endpoint hit: Search for employees by name, position or date")
			}*/
			if err := rows.Err(); err != nil{
				panic(err)
			}

	default:
		fmt.Fprintln(w, "Endpoint hit: This endpoint only supports GET requests")

	}
}
//search should always return an array of results.

//|| r.URL.Query().Get ("position") != "" || r.URL.Query().Get("startDate") != "" || r.URL.Query().Get("endDate") != ""

/* TO DO
Pagination
Patch for updating employee attributes address, phone number

*/


func createTable() {
	//CREATE TABLE "employees" ( "ID" INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE, "firstName" TEXT NOT NULL, "lastName" TEXT NOT NULL, "dateOfBirth" TEXT NOT NULL, "addressLineOne" TEXT NOT NULL, "addressLineTwo" TEXT, "city" TEXT NOT NULL, "postcode" TEXT NOT NULL, "startDate" TEXT NOT NULL, "nextOfKin" TEXT NOT NULL, "position" TEXT NOT NULL, "endDate" INTEGER, "recordCreatedDate" TEXT DEFAULT CURRENT_TIMESTAMP)
	/*INSERT INTO employees (firstName, lastName, dateOfBirth, addressLineOne, city, postcode, startDate, nextOfKin, position)
	  VALUES ('Salman','Ahmed','1999-01-01 00:00:00:000','33 Holborn', 'London', 'EC1N 2HT','2020-01-09 00:00:00.000', 'Steve Stotter','CEO');
	*/
	/*INSERT INTO employees (firstName, lastName, dateOfBirth, addressLineOne, city, postcode, startDate, nextOfKin, position)
	  VALUES ('Joe','Jenkins','1992-03-31 00:00:00:000','12 Little Tree Lane', 'London', 'E7 0TR','2020-03-10 00:00:00.000', 'Laura Jenkins','Head of Design');
	*/
}