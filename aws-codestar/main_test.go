package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestGetRandomPerson(t *testing.T) {
	t.Run("GetRandomPerson", func(t *testing.T) {
		gotPerson, err := GetRandomPerson()
		if err != nil || gotPerson == nil {
			t.Errorf("GetRandomPerson() error = %s", err)
			return
		}
	})
}

func TestHandler(t *testing.T) {
	type args struct {
		request events.APIGatewayProxyRequest
	}
	tests := []struct {
		name    string
		args    args
		want    events.APIGatewayProxyResponse
		wantErr bool
	}{{
		name: "Handler",
		args: args{events.APIGatewayProxyRequest{}},
		want: events.APIGatewayProxyResponse{StatusCode: 200},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Handler(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Handler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.StatusCode != tt.want.StatusCode {
				t.Errorf("Handler() = %v, want %v", got, tt.want)
			}
			t.Logf("Got %+v", got)
		})
	}
}
