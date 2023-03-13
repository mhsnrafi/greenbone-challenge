## Greenbone Challenge
### Requirements
1. The system administrator wants to be able to add a new computer to an employee
2. The system administrator wants to be informed when an employee is assigned 3 or
more computers
3. The system administrator wants to be able to get all computers
4. The system administrator wants to be able to get all assigned computers for an
employee
5. The system administrator wants to be able to get the data of a single computer
6. The system administrator wants to be able to remove a computer from an employee
7. The system administrator wants to be able to assign a computer to another employee

## Tech/framework used
1. Go programming language
2. GORM library for database operations
3. Redis (for cache)
4. Postgres DB
5. Docker

These technologies were chosen to provide a scalable, performant, and maintainable solution for computer tracking system

## Problem Statement
The problem is to create an application that tracks computers issued by a company and stores their details in a database. The required details for each computer are MAC address, computer name, IP address, employee abbreviation, and description. The system administrator wants to be able to perform CRUD operations and retrieve information about all computers via a REST interface. Additionally, if 3 or more devices are assigned to a single user, the administrator wants to be notified using a messaging service running in a Docker container. The messaging service listens to requests on port 8080, and the expected body of the REST endpoint is defined in the problem statement.


### Tool and techniques used
There are couple of API's Designed to solve this challenge and to secure the solution
1. Take JWT authentication approach to make API's secure, How we generate access token and refresh that tokens how implement below is the approach
2. The GenerateAccessTokens function creates two types of tokens for a given email.
3. These tokens are an access token and a refresh token. 
4. The access token has a set time to expire after a certain number of minutes and the refresh token has a set time to expire after a certain number of days. 
5. The function calls a CreateToken function twice to create both the access and refresh tokens and returns them. 
6. If there is an error during the creation of either token, the function returns an error. 
7. API Collection Json already included in the project
8. Proper logging is added used zap logger
9. To secure api, used JWT authentication mechanism


## How to run the solution, follow these steps:

Clone the repository: Run the following command to clone the repository to your local machine: 
```bash
git clone https://github.com/mhsnrafi/greenbone-challenge.git

```

### Install dependencies: 
Change into the project directory and run go mod download to install the required dependencies.

### Start the Project: 
Use the command:
```bash
docker-compose up
```

### Configure credentials: 
The credentials required to connect to the database and run the API are described in the .env.local file.

###  Generate access token: 
Call the "Generate access token" endpoint to obtain an access token, which is required to authorize the API calls. Add the header "Bearer-Token" to each API request, using the access token obtained in this step.

### Use the API: 
The Postman collection is attached for easy use of the API.
```json
Greenbone.postman_collection.json
```


### Endpoints
- Generate access token endpoint: `http://localhost:8000/v1/auth/generate_access_token`
- Refresh Token endpoint: `http://localhost:8000/v1/auth/refresh`
- Create Employee: `http://localhost:8000/v1/api/employees/`
- Create Computer: `http://localhost:8000/v1/computers`
- Get All Computer: `http://localhost:8000/v1/computers`
- Get Computer By Id: `http://localhost:8050/v1/computers/3`
- Delete Computer: `http://localhost:8000/v1/api/employees/computers/3/JDE`
- Get All Assigned Computer of Employee: `http://localhost:8050/v1/api/employees/computers/JDE`
- Get Assigned Computer of Another Employee: `http://localhost:8000/v1/computers/3/JAD`
- Swagger Endpoint: `http://localhost:8080/swagger/index.html#/`


### POST /auth/generate_access_token
This endpoint used to authenticate and validate the used is verified and generate access token details.

### POST /auth/refresh

This endpoint used refresh the access token

### Request Payload
```json
{
  "Token": "refresh_token",
  "Email": "user@example.com"
}
```

### Request Payload
#### Create Employee #1
```json
{
   "first_name": "David",
   "last_name": "Lee",
   "email": "david.lee@example.com",
   "abbreviation": "DLL",
   "computers": []
}
```

#### Create Employee #2
```json
{
   "first_name": "Alice",
   "last_name": "Johnson",
   "email": "alice.johnson@example.com",
   "abbreviation": "AJK",
   "computers": []
}
```

#### Create Employee #3
```json
 {
   "first_name": "John",
   "last_name": "Doe",
   "email": "john.doe@example.com",
   "abbreviation": "JDE",
   "computers": []
}
```

#### Create Employee #4
```json
{
   "first_name": "Bob",
   "last_name": "Smith",
   "email": "bob.smith@example.com",
   "abbreviation": "BSS",
   "computers": []
}
```






#### Create Computer #1
```json
{
  "mac_address": "12:34:56:78:90:ab",
  "computer_name": "John's Laptop",
  "ip_address": "192.168.1.103",
  "employee_abbrev": "JDE",
  "description": "MacBook Air"
}
```

#### Create Computer #2
```json
{
  "mac_address": "11:22:33:44:55:66",
  "computer_name": "John's  Desktop",
  "ip_address": "192.168.1.104",
  "employee_abbrev": "JDE",
  "description": "HP EliteDesk"
}
```

#### Create Computer #3
```json
{
  "mac_address": "ff:ee:dd:cc:bb:aa",
  "computer_name": "John's  Test Computer",
  "ip_address": "192.168.1.105",
  "employee_abbrev": "JDE",
  "description": "Virtual machine"
}
```

#### Create Computer #4
```json
{
  "mac_address": "aa:bb:cc:dd:ee:ff",
  "computer_name": "Alice's Desktop",
  "ip_address": "192.168.1.102",
  "employee_abbrev": "AJK",
  "description": "Custom-built PC"
}
```

#### Create Computer #5
```json
{
  "mac_address": "aa:bb:cc:dd:ee:ff",
  "computer_name": "Alice's Desktop",
  "ip_address": "192.168.1.102",
  "employee_abbrev": "AJK",
  "description": "Custom-built PC"
}
```


#### Create Computer #6
```json
{
  "mac_address": "00:11:22:33:44:55",
  "computer_name": "David's Laptop",
  "ip_address": "192.168.1.106",
  "employee_abbrev": "DLL",
  "description": "Lenovo ThinkPad"
}
```

#### Create Computer #7
```json
{
  "mac_address": "55:44:33:22:11:00",
  "computer_name": "David's Desktop",
  "ip_address": "192.168.1.107",
  "employee_abbrev": "DLL",
  "description": "Custom-built PC"
}
```

## Tests
The API includes a set of unit tests to ensure proper functionality. To run the tests, use the following command.
```bash
go test -v ./...
```

## API Documentation
To test the API endpoints directly from the documentation, making it easier to ensure that the API is working as expected build swagger api documentationa  user-friendly interface to quickly understand the API’s capabilities and functions
```bash
http://localhost:8080/swagger/index.html#/
```


## Improvement Area
Instead of send warning notification to system admin on a docker service, we need to be utilize messaging service like RabbitMQ can provide better reliability and scalability for sending notifications, as it allows for asynchronous message passing and can handle a large volume of messages. However, it also adds complexity to the system, as you need to set up and manage a RabbitMQ server and potentially write additional code to handle messaging

On the other hand, sending notifications directly to the Docker service may be simpler and more straightforward, as it doesn't require any additional infrastructure or code. However, it may be less scalable and reliable, as the Docker service may not be able to handle a large volume of requests or may be more prone to failure

This the area where we need some improvement to make a bettle reliable and scalable system. if we expect a high volume of notifications or need a high level of reliability, using a messaging service like RabbitMQ may be the better option. If we expect a low to moderate volume of notifications and simplicity is a priority, sending notifications directly to the Docker service may be sufficient.

