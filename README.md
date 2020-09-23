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

## Creating a User

Using a POST request in Postman enter the address http://localhost:4000/employees and select params.  The keys for an employee record are as follows, values will be the employee's individual personal details:

![Alt text](/create employee.png?raw=true "Create New Employee")

All listed fields are required except addressLineTwo which is optional.

If an employee record has been successfully created you will see a message indicated that this has happened.