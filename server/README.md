To run docker container:
```
docker run --name airb -e MYSQL_ROOT_PASSWORD=Password123! -p 3306:3306 -d mysql 
```
To create migration
```
migrate create -ext sql -dir .\migrations\ -seq <name>
```
Database Name: om

To run migration:
```
migrate  -path migrations  -database "mysql://root:Password123!@tcp(localhost:3306)/om" up
```

Swagger:
```
go install github.com/swaggo/swag/cmd/swag@latest
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
```
To generate swagger docs, run from `server/`:
```
swag init -g main.go -d cmd/api,internal/app/appmodels,internal/domain/models,internal/mapi/handlers,internal/mapi/dtos,internal/mapi/apimodels,internal/domain/enums --parseInternal
```
