package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
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

type Employee struct {
	ID                int            `json:"id"`
	FirstName         string         `json:"first_name"`
	LastName          string         `json:"last_name"`
	DateOfBirth       string         `json:"date_of_birth"`
	AddressLineOne    string         `json:"address_line_one"`
	AddressLineTwo    sql.NullString `json:"address_line_two,omitempty"`
	City              string         `json:"city"`
	Postcode          string         `json:"postcode"`
	StartDate         string         `json:"start_date"`
	NextOfKin         string         `json:"next_of_kin"`
	Position          string         `json:"position"`
	EndDate           sql.NullString `json:"end_date,omitempty"`
	RecordCreatedDate string         `json:"record_created_date"`
	employeeStatus    string         `json:"employeeStatus"`
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
	r := mux.NewRouter()
	//r.HandleFunc("/employees", employeeHandler)
	//http.HandleFunc("/employees/", employeeByIDHandler)
	//r.HandleFunc("/employees/{id}", getEmployeeByIDHandler).Methods("GET")
	//r.HandleFunc("/employees/{id}", employeeByIDHandler).Methods("DELETE")
	r.HandleFunc("/employees", employeeSearchHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":4000", r))

}

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
			err := rows.Scan(&employee.ID, &employee.LastName, &employee.FirstName, &employee.DateOfBirth, &employee.AddressLineOne, &employee.AddressLineTwo, &employee.City, &employee.Postcode, &employee.StartDate, &employee.NextOfKin, &employee.Position, &employee.EndDate, &employee.RecordCreatedDate)
			if err != nil {
				log.Fatal(err)
			}
			employees = append(employees, employee)

		}
		if err = rows.Err(); err != nil {
			log.Fatal(err)
		}
		for _, employee := range employees {
			json, err := json.MarshalIndent(employee, "", "")
			if err != nil {
				log.Println(err)
			}
			fmt.Fprint(w, string(json))
			//fmt.Fprintf(w, "%b, %s %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s\n", employee.ID, employee.LastName, employee.FirstName, employee.DateOfBirth, employee.AddressLineOne, employee.AddressLineTwo, employee.City, employee.Postcode, employee.StartDate, employee.NextOfKin, employee.Position, employee.EndDate, employee.RecordCreatedDate)
		}

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
			fmt.Println(w, err)
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

func getEmployeeByIDHandler(w http.ResponseWriter, r *http.Request) {
	muxvars := mux.Vars(r)
	ID := muxvars["id"]
	row := db.QueryRow("SELECT ID, firstName, lastName, dateOfBirth, addressLineOne, addressLineTwo, city, postcode, startDate, nextOfKin, position, endDate, recordCreatedDate, employeeStatus FROM employees WHERE ID =$1", ID)

	employee := new(Employee)
	err := row.Scan(&employee.ID, &employee.FirstName, &employee.LastName, &employee.DateOfBirth, &employee.AddressLineOne, &employee.AddressLineTwo, &employee.City, &employee.Postcode, &employee.StartDate, &employee.NextOfKin, &employee.Position, &employee.EndDate, &employee.RecordCreatedDate, &employee.employeeStatus)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		fmt.Println("sql no rows error", err)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}
	json, err := json.MarshalIndent(employee, "", "")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(w, string(json))

}

//Insert, delete and update do not return rows.
//func employeeByIDHandler(w http.ResponseWriter, r *http.Request) {
//	switch method := r.Method; method {
//	case "GET":
//		getID(r.URL.Path)
//		if ID == "" {
//			http.Error(w, http.StatusText(400), 400)
//			return
//		}
//		row := db.QueryRow("SELECT ID, firstName, lastName, dateOfBirth, addressLineOne, addressLineTwo, city, postcode, startDate, nextOfKin, position, endDate, recordCreatedDate, employeeStatus FROM employees WHERE ID =$1", ID)
//
//		employee := new(Employee)
//		err := row.Scan(&employee.ID, &employee.FirstName, &employee.LastName, &employee.DateOfBirth, &employee.AddressLineOne, &employee.AddressLineTwo, &employee.City, &employee.Postcode, &employee.StartDate, &employee.NextOfKin, &employee.Position, &employee.EndDate, &employee.RecordCreatedDate, &employee.employeeStatus)
//		if err == sql.ErrNoRows {
//			http.NotFound(w, r)
//			fmt.Println(w, "sql error", err)
//			return
//		} else if err != nil {
//			fmt.Println(w, err)
//			http.Error(w, http.StatusText(500), 500)
//			return
//		}
//		json, err := json.MarshalIndent(employee, "", "")
//		if err != nil {
//			log.Fatal(err)
//		}
//		fmt.Fprint(w, string(json))
//			//fmt.Fprintf(w, "%d, %s %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s\n", employee.ID, employee.FirstName, employee.LastName, employee.DateOfBirth, employee.AddressLineOne, employee.AddressLineTwo, employee.City, employee.Postcode, employee.StartDate, employee.NextOfKin, employee.Position, employee.EndDate, employee.RecordCreatedDate, employee.employeeStatus)
//		/*case "PATCH":
//			v:r.URL.Query()
//
//
//
//			result, err := db.Exec("UPDATE employees SET firstName = ?", firstName)
//			//result, err := db.Exec("UPDATE employees SET firstName = $1, lastName = $2, dateOfBirth =$3, addressLineOne =$4, addressLineTwo =$5, city =$6, postcode =$7, nextOfKin =$8, position =$9, endDate =$10 WHERE ID =?",ID)
//
//			if err != nil {
//				log.Println(err.Error())
//				return
//			}
//			rowsAffected, err := result.RowsAffected()
//			if err != nil {
//				http.Error(w, http.StatusText(500),500)
//				return
//			}
//
//			fmt.Fprintf(w, "Employee %s updated successfully (%d row affected)\n", ID, rowsAffected)
//*/
//		case "DELETE":
//			getID(r.URL.Path)
//			if ID == ""{
//				http.Error(w, http.StatusText(400), 400)
//				return
//			}
//			result, err := db.Exec("UPDATE employees SET employeeStatus ='Disabled' WHERE ID =$1", ID)
//
//			//employee := new(Employee)
//			//err := row.Scan(&employee.ID, &employee.FirstName, &employee.LastName, &employee.DateOfBirth, &employee.AddressLineOne, &employee.AddressLineTwo, &employee.City, &employee.Postcode, &employee.StartDate, &employee.NextOfKin, &employee.Position, &employee.EndDate, &employee.RecordCreatedDate, &employee.employeeStatus)
//
//			if err != nil{
//				http.Error(w, http.StatusText(500), 500)
//				fmt.Println("DB exec error", err)
//				return
//			}
//			rowsAffected, err := result.RowsAffected()
//			if err != nil {
//				http.Error(w, http.StatusText(500),500)
//				fmt.Println("Results row effected error", err)
//				return
//			}
//
//			fmt.Fprintf(w, "Employee %s deleted successfully (%d row affected)\n", ID, rowsAffected)
//		default:
//			fmt.Fprint(w, "Endpoint hit: This endpoint only supports GET, PATCH and DELETE requests by ID")
//		}
//	}

func employeeSearchHandler(w http.ResponseWriter, r*http.Request) {
	// 1. Get the filter criteria
	//1b. Get criteria to work with different fields.
	filterValues := r.URL.Query()
	fmt.Printf("%+v\n", filterValues)

	if filterValues["first_name"][0] != ""{
		fmt.Printf("Value: %s *****\n", filterValues["first_name"][0])

		//var firstName string = filterValues["first_name"][0]
		var firstName string = filterValues["first_name"][0]
		rows, err := db.Query("SELECT ID, firstName, lastName, dateOfBirth, addressLineOne, addressLineTwo, city, postcode, startDate, nextOfKin, position, endDate, recordCreatedDate FROM employees WHERE firstName LIKE $1", firstName)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
			return

		}
		defer rows.Close()

		employees := make([]*Employee, 0)
		for rows.Next() {
			fmt.Println("Employee found")
			employee := new(Employee)
			err := rows.Scan(&employee.ID, &employee.FirstName, &employee.LastName, &employee.DateOfBirth, &employee.AddressLineOne, &employee.AddressLineTwo, &employee.City, &employee.Postcode, &employee.StartDate, &employee.NextOfKin, &employee.Position, &employee.EndDate, &employee.RecordCreatedDate)
			if err != nil {
				log.Fatal(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, err)
				return
			}
			employees = append(employees, employee)

		}
		fmt.Printf("%+v\n", employees)
		if err = rows.Err(); err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, err)
			return
		}
		for _, employee := range employees {
			json, err := json.MarshalIndent(employee, "", "")
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, err)
				return
			}
			fmt.Fprint(w, string(json))
		}
	}

}

/*func createTable() {
	CREATE TABLE "employees" ( "ID" INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE, "firstName" TEXT NOT NULL, "lastName" TEXT NOT NULL, "dateOfBirth" TEXT NOT NULL, "addressLineOne" TEXT NOT NULL, "addressLineTwo" TEXT, "city" TEXT NOT NULL, "postcode" TEXT NOT NULL, "startDate" TEXT NOT NULL, "nextOfKin" TEXT NOT NULL, "position" TEXT NOT NULL, "endDate" INTEGER, "recordCreatedDate" TEXT DEFAULT CURRENT_TIMESTAMP)
	INSERT INTO employees (firstName, lastName, dateOfBirth, addressLineOne, city, postcode, startDate, nextOfKin, position)
	VALUES ('Salman','Ahmed','1999-01-01 00:00:00:000','33 Holborn', 'London', 'EC1N 2HT','2020-01-09 00:00:00.000', 'Steve Stotter','CEO');

	INSERT INTO employees (firstName, lastName, dateOfBirth, addressLineOne, city, postcode, startDate, nextOfKin, position)
	  VALUES ('Joe','Jenkins','1992-03-31 00:00:00:000','12 Little Tree Lane', 'London', 'E7 0TR','2020-03-10 00:00:00.000', 'Laura Jenkins','Head of Design');

}*/