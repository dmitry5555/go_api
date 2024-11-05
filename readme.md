# about
api @ go using net/http go package and sqlite db

- /users				get all users
- /users/{id}			get user by id
- /users/add 			[POST] add user(name, email)
						test: curl -X POST http://localhost:8080/users/add -H "Content-Type: application/json" -d '{"name": "John Doe", "email": "john@example1.com"}'

# init project
go mod init myproject

# run
go run main.go

go mod tidy: clean go.mod

# etc
go build: compile to bin
go run: compile and run
go fmt: fmt the code

