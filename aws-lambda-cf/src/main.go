package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
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

func GetRandomPeople(numResults int) (person []*Person, err error) {
	resp, err := http.Get(os.Getenv("RANDOMUSER_API_URL") + "?results=" + string(numResults))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	log.Printf("Response from RANDOMUSER API\n%d\t%s", resp.StatusCode, resp.Status)

	rawResponse, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("Response body from RANDOMUSER API\n%s", string(rawResponse))

	var responseJSON RandomUserMeJSON
	err = json.Unmarshal(rawResponse, &responseJSON)
	if err != nil {
		return nil, err
	} else if len(responseJSON.Results) < 1 {
		return nil, errors.New("Nothing returned from api!")
	}

	return responseJSON.Results, nil
}

func APIRequestHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	numPeople, err := strconv.Atoi(request.QueryStringParameters["amount"])
	if err != nil || numPeople < 1 {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"message": "The 'amount' query parameter is required, and must be a number greater than 0!"}`}, nil
	}

	people, err := GetRandomPeople(numPeople)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusServiceUnavailable}, err
	}

	peopleJSON, err := json.Marshal(people)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(peopleJSON),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil

}

func main() {
	lambda.Start(APIRequestHandler)
}
