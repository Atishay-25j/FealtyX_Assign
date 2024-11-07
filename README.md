This project provides a simple REST API built with Go to manage a list of students. The API supports basic CRUD (Create, Read, Update, Delete) operations and integrates with the Ollama API to generate an AI-based summary for each student.

Table of Contents
Requirements
Installation
Project Structure
API Documentation
Create a New Student
Get All Students
Get a Student by ID
Update a Student by ID
Delete a Student by ID
Generate a Summary for a Student by ID
Running the Project
Testing with Postman
Troubleshooting
Requirements
Go (version 1.16 or later)
Ollama installed locally
A REST client like Postman or curl
Installation
Clone the repository:

bash
Copy code
git clone https://github.com/Atishay-25j/FealtyX_Assign.git
cd students.go
Initialize the Go module (if not already done):

bash
Copy code
go mod init FealtyX_Assign
Install dependencies:

For routing, we use gorilla/mux:
bash
Copy code
go get -u github.com/gorilla/mux
Any other dependencies should be installed automatically as they’re listed in the go.mod file.
Set up and verify Ollama:

Download and install Ollama on your local machine.
Make sure Ollama is running and accessible at http://localhost:11434. Start Ollama if it's not already running:
bash
Copy code
ollama start
Project Structure
main.go: Main file containing API routes and handlers.
student.go: Contains handlers for CRUD operations and the summary generation.
go.mod and go.sum: Go module files that manage dependencies.
API Documentation
Base URL
All requests assume that the server is running locally:

bash
Copy code
http://localhost:8080/students
Endpoints
1. Create a New Student
Method: POST
URL: /students
Headers:
Content-Type: application/json
Body:
json
Copy code
{
  "name": "John Doe",
  "age": 20,
  "email": "johndoe@example.com"
}
Response:
201 Created on success
400 Bad Request if the request body is invalid
2. Get All Students
Method: GET
URL: /students
Response:
200 OK with a JSON array of students
Example:
json
Copy code
[
  {
    "id": 1,
    "name": "John Doe",
    "age": 20,
    "email": "johndoe@example.com"
  }
]
3. Get a Student by ID
Method: GET
URL: /students/{id}
Response:
200 OK with a JSON object of the student
404 Not Found if the student does not exist
4. Update a Student by ID
Method: PUT
URL: /students/{id}
Headers:
Content-Type: application/json
Body:
json
Copy code
{
  "name": "Jane Doe",
  "age": 22,
  "email": "janedoe@example.com"
}
Response:
200 OK on successful update
400 Bad Request if the request body is invalid
404 Not Found if the student does not exist
5. Delete a Student by ID
Method: DELETE
URL: /students/{id}
Response:
204 No Content on successful deletion
404 Not Found if the student does not exist
6. Generate a Summary for a Student by ID
Method: GET
URL: /students/{id}/summary
Response:
200 OK with a summary of the student
404 Not Found if the student does not exist
500 Internal Server Error if there is an error communicating with the Ollama API
Example Response:
json
Copy code
{
  "summary": "John Doe is a 20-year-old student with the email johndoe@example.com."
}
Running the Project
Run the Go server:

bash
Copy code
go run main.go
This starts the server on http://localhost:8080.

Test the Endpoints: Use a REST client like Postman to interact with each endpoint (details in the next section).

Testing with Postman
Importing Requests
Create a collection in Postman for each API request, using the following examples:

Create Student (POST): http://localhost:8080/students
Get All Students (GET): http://localhost:8080/students
Get Student by ID (GET): http://localhost:8080/students/{id}
Update Student by ID (PUT): http://localhost:8080/students/{id}
Delete Student by ID (DELETE): http://localhost:8080/students/{id}
Generate Summary by ID (GET): http://localhost:8080/students/{id}/summary
Make sure to replace {id} with the actual student ID for the requests that require it.

Example Usage
To create a new student in Postman:

Select POST as the method.
Set the URL to http://localhost:8080/students.
In Headers, add Content-Type: application/json.
In Body, select raw and enter the JSON object:
json
Copy code
{
  "name": "John Doe",
  "age": 20,
  "email": "johndoe@example.com"
}
Send the request, and you should receive a 201 Created response.
Troubleshooting
Error Generating Summary
If you encounter an error when generating a summary for a student:

Make sure Ollama is running by checking the port 11434.
Verify the ollamaURL in your code is set to http://localhost:11434/api/v1/generate.
If the error persists, try restarting Ollama by running:
bash
Copy code
ollama start
Check Port Availability
If the server fails to start, make sure port 8080 is not in use by another process. If needed, you can change the port in main.go:

go
Copy code
http.ListenAndServe(":8080", r)
