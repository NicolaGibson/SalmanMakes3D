###  Salman Makes


Salman Makes 3D is a new start-up in the world of additive manufacturing. It has recently started growing at a remarkable pace. Their HR department has been managing employees on paper but it is no longer suitable due to the rapid growth of the company.
 
You are required to create an Employee Management System for their HR system to replace their paper based system.

##### You will need to design a RESTful API that allow clients to communicate with.
 
The API needs to have the following features;

1. Ability to create an employee record with the following details:
    * First name
    * Last name
    * DOB
    * Address line 1
    * Address line 2
    * City
    * Postcode
    * Start date
    * Next of kin
    * Position
    * End date
    
2. Ability to get employee details by their unique id
3. Ability to update employee details
4. Ability to delete employee details
5. Ability to get all current employees details
6. Ability to search for employees by their;
    * First name
    * Last name
    * Position
    * Start date
    * End date
 
Notes:
 
* Data will need to be stored in a database
* Reasonable error messages should be returned when an action cannot be completed
* This is a prototype and must be easy to change and modify as requirements change (hint: use interfaces, polymorphism, etc)
* The business requirements need to be fully tested in code
* The system needs to be fully documented
* Non-functional requirements should be considered, even if not fully implemented in code.
* The REST API needs to comply with industry standards
* No front-end needs to be developed for this project

-----

## Instructions for Use 

This API runs on localhost:4000 and handles creating an employee, updating an employees details, deleting an employee, searching for employee(s) using their personal details and searching for all employees.

Functionality can be accessed and edited using Postman or in the IDE using curl requests.

## Getting Started

Enter go run main.go in the terminal to start the server.

## Creating an Employee

Using a POST request in Postman enter http://localhost:4000/employees and select params.  The keys for an employee record are as follows, values will be the employee's individual personal details:

![Alt text](/create_employee.png "Create New Employee")

All listed fields are required except addressLineTwo which is optional. Once all key and values have been entered press the Send button.

If an employee record has been successfully created you will see a message indicated that this has happened.

## Deleting an Employee 

Employees are deleted by their ID number, please make sure you know the ID number before deleting an employee.  In the example below we will delete employee 28.

![Alt text](/delete_employee.png "Delete Employee")

Using a DELETE request in Postman enter http://localhost:4000/employees/28, then press the Send button.

If an employee has been successfully deleted you will see the following message: "Employee 28 deleted successfully (1 row affected)".


## Update an Employee's Details

An employee's details can be updated by their ID, details are updated in a JSON object and keys must match those specified here. One or more of the following details can be edited at a time:

first_name         
last_name          
date_of_birth       
address_line_one   
address_line_ywo    
city              
postcode          
start_date         
next_Of_kin         
position          
end_date           
	         
Using a PATCH request in Postman enter http://localhost:4000/employees/28

## Search for an Employee by ID

This search requires an employee's ID, please make sure that you know it before starting the search.  In the example below we will be using employee ID 28.

Using a GET request in Postman enter http://localhost:4000/employees/28, then press the Send button.

## Search for an Employee(s)