package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
)

type Person struct {
	Gender, Email, DOB, Nat string

	Name struct {
		Title, First, Last string
	}
	Location struct {
		Street, City, State string
	}
	Picture struct {
		Large string
	}
}

type RandomUserMeJSON struct {
	Results []*Person
	Info    struct {
		Version string
	}
}

func GetRandomPerson() (person *Person, err error) {
	resp, err := http.Get("https://randomuser.me/api/?results=1")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var responseJSON RandomUserMeJSON
	err = json.NewDecoder(resp.Body).Decode(&responseJSON)
	if err != nil {
		return nil, err
	} else if len(responseJSON.Results) < 1 {
		return nil, errors.New("Nothing returned!")
	}

	return responseJSON.Results[0], nil
}

func Handler() error {
	person, err := GetRandomPerson()
	if err != nil {
		return err
	}

	log.Printf("Logging random person %+v", person)

	return nil

}

func main() {
	lambda.Start(Handler)
}
