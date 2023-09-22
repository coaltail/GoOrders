# GoOrders

Backend for a social media managing app. The backend is created entirely in GO, using Fiber and GORM as its main frameworks.


## How to run
The backend uses docker to locally host the database and the server. It is recommended to install Docker Desktop and to run it with the Visual Studio Code extension "DevContainers".
If another code editor is being used, it can be run as follows:
1. Compose and build the project: 
```docker
docker compose build
dokcer compose up
```
2. Run the file "main.go" to start the server:
```go
go run main.go
```
Or, alternatively, if you're using air:
```go
air run main.go
```


