package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	//"github.com/lann/builder"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"strconv"
	sq "github.com/Masterminds/squirrel"


)
/*
 localhost:4000/employees/search/?name=Joe
 */

var ID = ""

//type UpdateBuilder builder.Builder

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
	r.Path("/employees").HandlerFunc(employeeHandler)
	r.HandleFunc("/employees", employeeHandler)
	r.HandleFunc("/employees/{id:[0-9]+}", getEmployeeByIDHandler).Methods("GET")
	r.HandleFunc("/employees/{id:[0-9]+}", deleteEmployeeByIDHandler).Methods("DELETE")
	r.HandleFunc ("/employees/{id:[0-9]+}", updateEmployeeByIDHandler).Methods("PATCH")
	r.HandleFunc("/employees", employeeSearchHandler).Methods("GET").Queries("key, {[0-9A-Za-z_]}")
	log.Fatal(http.ListenAndServe(":4000", r))
}

func employeeHandler(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case "GET":
		//mysqlite := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
		users := sq.Select("ID, firstName, lastName, dateOfBirth, addressLineOne, addressLineTwo, city, postcode, startDate, nextOfKin, position, endDate, recordCreatedDate, employeeStatus").From("employees")
		sql, args, err := users.ToSql()
		fmt.Println(sql, args, err)


		rows, err := db.Query(sql)
		if err != nil {
			fmt.Println("error: ", err)
			log.Fatal(err)
		}
		defer rows.Close()

		employees := make([]*Employee, 0)
		for rows.Next() {
			employee := new(Employee)
			err := rows.Scan(&employee.ID, &employee.LastName, &employee.FirstName, &employee.DateOfBirth, &employee.AddressLineOne, &employee.AddressLineTwo, &employee.City, &employee.Postcode, &employee.StartDate, &employee.NextOfKin, &employee.Position, &employee.EndDate, &employee.RecordCreatedDate, &employee.employeeStatus)
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
func deleteEmployeeByIDHandler(w http.ResponseWriter, r *http.Request){
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

	/*if employeeReq.AddressLineTwo != "" {
		b = b.Set("addressLineTwo", employeeReq.AddressLineOne)
	}*/

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

	/*if employeeReq.EndDate == "inactive" {
		b = b.Set("endDate", employeeReq.EndDate)
	}*/

	mysql, args, err := b.ToSql()
	if err != nil {
		log.Fatal("err toSQL: ", err)
	}
	fmt.Println("My final SQL query with args>>>>>", mysql,args)

	_, err = b.Exec()
	if err != nil {
		log.Fatal("error executing query: ", err)
	}

	return

//stmtcacher squirrel? instead of transaction?
		//tx, _ := db.Begin()
		//stmt, _ := tx.Prepare("UPDATE employees SET firstName = ?, lastName = ? WHERE ID = ?")
		//
		//result, err := stmt.Exec(firstName, lastName, ID)
		//if err != nil {
		//	w.WriteHeader(http.StatusInternalServerError)
		//	fmt.Fprint(w, err)
		//	return
		//}
		//tx.Commit()
		//if err == sql.ErrNoRows {
		//	http.NotFound(w, r)
		//	fmt.Println("sql no rows error", err)
		//	return
		//} else if err != nil {
		//	w.WriteHeader(http.StatusInternalServerError)
		//	fmt.Fprint(w, err)
		//	return
		//}
		//rowsAffected, err := result.RowsAffected()
		//if err != nil {
		//	log.Fatal(err)
		//	return
		//}
		//fmt.Fprintf(w, "Employee %d updated successfully (%d row affected)\n", ID, rowsAffected)

	//}
}

func employeeSearchHandler(w http.ResponseWriter, r*http.Request) {
	filterValues := r.URL.Query()
	fmt.Println("%+v\n", filterValues)
	var employeeReq Employee
	fmt.Printf("empReq: %+v******\n", employeeReq)
	employees := sq.Select("ID, firstName, lastName, dateOfBirth, addressLineOne, addressLineTwo, city, postcode, startDate, nextOfKin, position, endDate, recordCreatedDate, employeeStatus").From("employees")
	d := employees.Where(sq.Eq{"filterValues": filterValues})
	sqlStr, args, err := d.ToSql()
	return
	sqlStr == "SELECT ID, firstName, lastName, dateOfBirth, addressLineOne, addressLineTwo, city, postcode, startDate, nextOfKin, position, endDate, recordCreatedDate, employeeStatus FROM employees WHERE filterValues = ?, filterValues"
	fmt.Println(sqlStr, args, err)

	if err != nil {
		fmt.Println("error: ", err)
		log.Fatal(err)
	}
	sql := &bytes.Buffer{}

	if len(d.Prefixes) > 0 {
		args, err = appendToSql(d.Prefixes, sql, " ", args)
		if err != nil {
			return
		}

		sql.WriteString(" ")
	}
	if employeeReq.FirstName != "" {
		d = d.Select("firstName", employeeReq.FirstName)
	}

	if employeeReq.LastName != "" {
		d = d.Set("lastName", employeeReq.LastName)
	}

	if employeeReq.DateOfBirth != "" {
		d = d.Set("dateOfBirth", employeeReq.DateOfBirth)
	}

	if employeeReq.AddressLineOne != "" {
		d = d.Set("addressLineOne", employeeReq.AddressLineOne)
	}

	if employeeReq.AddressLineTwo != "" {
		d = d.Set("addressLineTwo", employeeReq.AddressLineOne)
	}

	if employeeReq.City != "" {
		d = d.Set("city", employeeReq.City)
	}

	if employeeReq.Postcode != "" {
		d = d.Set("postcode", employeeReq.Postcode)
	}

	if employeeReq.StartDate != "" {
		d = d.Set("startDate", employeeReq.StartDate)
	}

	if employeeReq.NextOfKin != "" {
		d = d.Set("nextOfKin", employeeReq.NextOfKin)
	}
	if employeeReq.Position != "" {
		d = d.Set("position", employeeReq.Position)
	}

	if employeeReq.EndDate == "inactive" {
		d = d.Set("endDate", employeeReq.EndDate)
	}

	mysql, args, err := d.ToSql()
	if err != nil {
		log.Fatal("err toSQL: ", err)
	}
	fmt.Println("My final SQL query with args>>>>>", mysql, args)

	_, err = d.Exec()
	if err != nil {
		log.Fatal("error executing query: ", err)
	}

	return
}
	/*
		users := sq.Select("ID, firstName, lastName, dateOfBirth, addressLineOne, addressLineTwo, city, postcode, startDate, nextOfKin, position, endDate, recordCreatedDate, employeeStatus").From("employees")
		sql, args, err := users.ToSql()
		fmt.Println(sql, args, err)

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
			} */

