package main

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sfn"
)

type request events.APIGatewayProxyRequest

type response events.APIGatewayProxyResponse

type requestParams struct {
	ResponseID  string `json:"responseId"`
	QueryResult struct {
		QueryText  string `json:"queryText"`
		Action     string `json:"action"`
		Parameters struct {
			Time time.Time `json:"time"`
		} `json:"parameters"`
	} `json:"queryResult"`
}

func Handler(ctx context.Context, req request) (response, error) {
	var params requestParams
	json.Unmarshal([]byte(req.Body), &params)

	// new client
	sess := session.Must(session.NewSession())
	svc := sfn.New(sess)

	// parameters for StepFunctions
	layout := "2006-01-02T15-04-05"
	name := "TurnOnAirconStepFunc-" + time.Now().Format(layout)
	arn := os.Getenv("TurnOnAirconStepFuncARN")

	layout = "2006-01-02T15:04:05Z07:00"
	timeToStart := params.QueryResult.Parameters.Time
	input := struct {
		StartDate string `json:"start_date"`
	}{
		timeToStart.Format(layout),
	}

	inputJson, err := json.Marshal(input)
	if err != nil {
		return response{StatusCode: 500}, err
	}

	// kick StepFunctions
	_, err = svc.StartExecution(&sfn.StartExecutionInput{
		Name:            aws.String(name),
		StateMachineArn: aws.String(arn),
		Input:           aws.String(string(inputJson)),
	})
	if err != nil {
		return response{StatusCode: 500}, err
	}

	resp := response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            timeToStart.Format(layout),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
