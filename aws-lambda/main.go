package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Person struct {
	Gender string
	Email  string
	DOB    string
	Name   *struct {
		Title string
		First string
		Last  string
	}
	Picture *struct {
		Large     string
		Medium    string
		Thumbnail string
	}
	Location *struct {
		Street   string
		City     string
		State    string
		Postcode int
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

// Handler is executed by AWS Lambda in the main function. Once the request
// is processed, it returns an Amazon API Gateway response object to AWS Lambda
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	person, err := GetRandomPerson()
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusServiceUnavailable}, err
	}

	personJSON, err := json.Marshal(person)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(personJSON),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil

}

func main() {
	lambda.Start(Handler)
}
