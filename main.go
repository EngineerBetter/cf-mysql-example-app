package main

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/EngineerBetter/cf-mysql-example-app/mysql"
	"github.com/cloudfoundry-community/go-cfenv"
)

func main() {
	appEnv, err := cfenv.Current()
	if err != nil {
		log.Fatal("Could not get CF env:", err)
	}

	services, err := appEnv.Services.WithTag("mysql")
	if err != nil || len(services) == 0 {
		log.Fatal("Could not find service with mysql tag:", err)
	}

	username, _ := services[0].CredentialString("username")
	password, _ := services[0].CredentialString("password")
	hostname, _ := services[0].CredentialString("hostname")
	dbName, _ := services[0].CredentialString("name")

	dbUrl := username + ":" + password + "@tcp(" + hostname + ":3306)/" + dbName

	repo, err := mysql.NewMySQLRepository(dbUrl)
	if err != nil {
		log.Fatal("Could not start MySQLRepository:", err)
	}

	handler := NewPutGetHandler(repo)
	err = http.ListenAndServe(":"+os.Getenv("PORT"), handler)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

type Repository interface {
	Write(key, value string) error
	Read(key string) (string, error)
}

type PutGetHandler struct {
	repo Repository
}

func NewPutGetHandler(repository Repository) http.Handler {
	var handler PutGetHandler
	handler.repo = repository
	return &handler
}

func (handler *PutGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	re := regexp.MustCompile("^/")
	key := re.ReplaceAllString(r.URL.Path, "")

	responseBody := "bad request"
	statusCode := http.StatusBadRequest

	if r.Method == http.MethodPut {
		buff := new(bytes.Buffer)
		buff.ReadFrom(r.Body)
		value := buff.String()

		err := handler.repo.Write(key, value)

		if err != nil {
			statusCode = http.StatusInternalServerError
			responseBody = err.Error()
		} else {
			statusCode = http.StatusCreated
			responseBody = "created"
		}
	} else if r.Method == http.MethodGet {
		value, err := handler.repo.Read(key)

		if err != nil {
			statusCode = http.StatusInternalServerError
			responseBody = err.Error()
		} else {
			if value == "" {
				statusCode = http.StatusNotFound
				responseBody = "key not found"
			} else {
				statusCode = http.StatusOK
				responseBody = value
			}
		}
	}

	w.WriteHeader(statusCode)
	w.Write([]byte(responseBody + "\n"))
}
