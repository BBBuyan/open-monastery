Database Name: om

To run docker container:
```
docker run --name om-db -e MYSQL_ROOT_PASSWORD=Password123! -p 3306:3306 -d mysql 
```

To create migration
```
task mig-create -- <mig_name>
```

To run migration:
```
task mig-run
```

Swagger:
```
go install github.com/swaggo/swag/cmd/swag@latest
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
```

To generate swagger docs, run from `server/`:
```
task swag-gen
```
