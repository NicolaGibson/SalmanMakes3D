package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestemployeeHandler(t *testing.T){
// It should return all employees
// It should create a new employee
	req := httptest.NewRequest(http.MethodGet, "/employees", nil)
	res := httptest.NewRecorder()

	employeeHandler(res, req)

	result := res.Result()

	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		t.Fatal(err)
	}
	result.Body.Close()

	if http.StatusOK != result.StatusCode  {
		t.Error("expected", http.StatusOK, "got", result.StatusCode)
	}
	if "employee" != string(body) {
		t.Error("expected employee got", string(body))
	}
}


func TestgetEmployeeByIDHandler(t *testing.T){
	// It should return one specific employee when given an ID
	t.Error()
}

func TestdeleteEmployeeByIDHandler(t *testing.T){
	// It should delete one specific employee when given an ID
	t.Error()
}

func TestupdateEmployeeByIDHandler(t *testing.T){
	// It should update one specific employee's details when given an ID
	t.Error()
}

func TestemployeeSearch(t *testing.T){
	// It should search for all employees that meet search criteria.
	t.Error()
}