package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"strconv"
)

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

func init() {
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
	r.HandleFunc("/employees/{id:[0-9]+}", getEmployeeByIDHandler).Methods("GET")
	r.HandleFunc("/employees/{id:[0-9]+}", deleteEmployeeByIDHandler).Methods("DELETE")
	r.HandleFunc("/employees/{id:[0-9]+}", updateEmployeeByIDHandler).Methods("PATCH")
	r.HandleFunc("/employees", employeeSearchHandler).Methods("GET")
	r.HandleFunc("/employees", createEmployeeHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(":4000", r))
}

func createEmployeeHandler(w http.ResponseWriter, r *http.Request) {
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

	if firstName == "" || lastName == "" || addressLineOne == "" || city == "" || postcode == "" || startDate == "" || nextOfKin == "" || position == "" || employeeStatus == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	result, err := db.Exec("INSERT INTO employees (firstName, lastName, dateOfBirth, addressLineOne, addressLineTwo, city, postcode, startDate, nextOfKin, position, employeeStatus) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)", firstName, lastName, dateOfBirth, addressLineOne, addressLineTwo, city, postcode, startDate, nextOfKin, position, employeeStatus)
	if err != nil {
		fmt.Println(w, err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	fmt.Fprintf(w, "Employee %s created successfully (%d row affected)\n", firstName+" "+lastName, rowsAffected)

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
func deleteEmployeeByIDHandler(w http.ResponseWriter, r *http.Request) {
	muxvars := mux.Vars(r)
	ID := muxvars["id"]
	result, err := db.Exec("UPDATE employees SET employeeStatus ='Disabled' WHERE ID =$1", ID)
	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		fmt.Println("sql no rows error", err)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Fprintf(w, "Employee %s deleted successfully (%d row affected)\n", ID, rowsAffected)
}
func updateEmployeeByIDHandler(w http.ResponseWriter, r *http.Request) {
	muxvars := mux.Vars(r)
	id := muxvars["id"]
	ID, _ := strconv.Atoi(id)
	fmt.Println("ID: ", ID)

	var employeeReq Employee
	err := json.NewDecoder(r.Body).Decode(&employeeReq)
	if err != nil {
		log.Fatal("err decoding req body: ", err)
	}

	fmt.Printf("empReq: %+v******\n", employeeReq)

	b := sq.Update("employees").Where(sq.Eq{"id": ID}).RunWith(db)

	if employeeReq.FirstName != "" {
		b = b.Set("firstName", employeeReq.FirstName)
	}

	if employeeReq.LastName != "" {
		b = b.Set("lastName", employeeReq.LastName)
	}

	if employeeReq.DateOfBirth != "" {
		b = b.Set("dateOfBirth", employeeReq.DateOfBirth)
	}

	if employeeReq.AddressLineOne != "" {
		b = b.Set("addressLineOne", employeeReq.AddressLineOne)
	}

	if employeeReq.AddressLineTwo.String != "" {
		b = b.Set("addressLineTwo", employeeReq.AddressLineTwo)
	}

	if employeeReq.City != "" {
		b = b.Set("city", employeeReq.City)
	}

	if employeeReq.Postcode != "" {
		b = b.Set("postcode", employeeReq.Postcode)
	}

	if employeeReq.StartDate != "" {
		b = b.Set("startDate", employeeReq.StartDate)
	}

	if employeeReq.NextOfKin != "" {
		b = b.Set("nextOfKin", employeeReq.NextOfKin)
	}
	if employeeReq.Position != "" {
		b = b.Set("position", employeeReq.Position)
	}

	if employeeReq.EndDate.String != "inactive" {
		b = b.Set("endDate", employeeReq.EndDate)
	}

	mysql, args, err := b.ToSql()
	if err != nil {
		log.Fatal("err toSQL: ", err)
	}
	fmt.Println("My final SQL query with args>>>>>", mysql, args)

	_, err = b.Exec()
	if err != nil {
		log.Fatal("error executing query: ", err)
	}

	return

}

func employeeSearchHandler(w http.ResponseWriter, r *http.Request) {
	filterValues := r.URL.Query()
	fmt.Printf("filterValues: %+v\n", filterValues)

	employeesSQL := sq.Select("ID, firstName, lastName, dateOfBirth, addressLineOne, addressLineTwo, city, postcode, startDate, nextOfKin, position, endDate, recordCreatedDate, employeeStatus").From("employees").RunWith(db)

	for k, v := range filterValues {
		switch k {
		case "first_name":
			employeesSQL = employeesSQL.Where("firstName = ?", v[0])
		case "last_name":
			employeesSQL = employeesSQL.Where("lastName = ?", v[0])
		case "date_of_birth":
			employeesSQL = employeesSQL.Where("date_of_birth = ?", v[0])
		}
	}

	sql, args, err := employeesSQL.ToSql()
	fmt.Printf("SQL: %v, Args: %+v, Err: %v\n", sql, args, err)

	rows, err := employeesSQL.Query()
	if err != nil {
		fmt.Println("db query error: ", err)
		log.Fatal(err)
	}
	defer rows.Close()

	employees := make([]*Employee, 0)
	for rows.Next() {
		employee := new(Employee)
		err := rows.Scan(&employee.ID, &employee.FirstName, &employee.LastName, &employee.DateOfBirth, &employee.AddressLineOne, &employee.AddressLineTwo, &employee.City, &employee.Postcode, &employee.StartDate, &employee.NextOfKin, &employee.Position, &employee.EndDate, &employee.RecordCreatedDate, &employee.employeeStatus)
		if err != nil {
			fmt.Println("db rows.Scan error: ", err)
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
			fmt.Println("JSON marshall error: ", err)
			log.Println(err)
		}
		fmt.Fprint(w, string(json))
	}
}
