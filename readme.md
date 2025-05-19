# 🚀 Setting Up a Go Backend

Follow these steps to initialize and run a Go backend with necessary dependencies.

## 0️⃣ Create .env File as in .env.example

## 1️⃣ Initialize Go Module  
Run the following command to create a `go.mod` file:

```sh
go mod init backend
```

## 2️⃣ Fetch Dependencies  
To create a `go.sum` file and add packages to `go.mod`, run:

```sh
go get .
```

## 3️⃣ Install Required Packages  
Run the following commands to install essential Go packages:

```sh
go get -u github.com/gin-gonic/gin
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
go get -u github.com/joho/godotenv
```

## 4️⃣ Run the Application  
Execute the Go application using:

```sh
go run main.go
```

## 🎯 Notes:
- Ensure **Go** is installed on your system.
- The `go.mod` and `go.sum` files should be committed to version control.
- If you face dependency issues, try:

  ```sh
  go mod tidy
  ```

Happy coding! 🚀