package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/excoriate/aws-secrets-rotation-lambda/internal/adapter"
	"github.com/excoriate/aws-secrets-rotation-lambda/internal/client"
	sm "github.com/excoriate/aws-secrets-rotation-lambda/internal/rotation"
)

//var client = lambda.New(session.New())
//
//func callLambda() (string, error) {
//	input := &lambda.GetAccountSettingsInput{}
//	req, resp := client.GetAccountSettingsRequest(input)
//	err := req.Send()
//	output, _ := json.Marshal(resp.AccountUsage)
//	return string(output), err
//}
//
func handleRequest(ctx context.Context, event sm.RotationEvent) (string, error) {
	// event
	eventJson, _ := json.MarshalIndent(event, "", "  ")
	fmt.Print(eventJson)

	awsCfg, err := adapter.NewAWS("us-east-1")
	if err != nil {
		return "", fmt.Errorf("failed to create adapter client: %v", err)
	}

	sm := client.NewSecretsManager(awsCfg)

	// List all secrets
	result, err := sm.ListAll()
	if err != nil {
		return "", fmt.Errorf("failed to list secrets: %v", err)
	}

	for _, secret := range result {
		fmt.Print(secret)
	}

	return "Hello, World!", nil
}

func main() {
	lambda.Start(handleRequest)
}
