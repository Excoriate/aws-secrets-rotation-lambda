package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/lambdacontext"
	sm "github.com/excoriate/aws-secrets-rotation-lambda/internal/rotation"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

//var mockCreateSecretEvent = "../../../mock/events/secretsmanager-event-1.json"
//var mockCreateSecretEvent = "../../../mock/events/secret-to-rotate.json"
var mockCreateSecretEvent = "../../../mock/events/secret-to-rotate-valid.json"

func TestMainTest(t *testing.T) {
	d := time.Now().Add(50 * time.Millisecond)
	_ = os.Setenv("AWS_LAMBDA_FUNCTION_NAME", "rotator-lambda-go")

	ctx, _ := context.WithDeadline(context.Background(), d)
	ctx = lambdacontext.NewContext(ctx, &lambdacontext.LambdaContext{
		AwsRequestID:       "495b12a8-xmpl-4eca-8168-160484189f99",
		InvokedFunctionArn: "arn:adapter:lambda:us-east-2:123456789012:function:blank-go",
	})

	t.Run("ReadEventJSON", func(t *testing.T) {
		inputJson := ReadJSONFromFile(t, mockCreateSecretEvent)
		var event sm.Event
		err := json.Unmarshal(inputJson, &event)

		assert.NoError(t, err, "should not error")
		assert.NotNil(t, event, "event should not be nil")
	})

	t.Run("HandleRequest", func(t *testing.T) {
		inputJson := ReadJSONFromFile(t, mockCreateSecretEvent)
		var event sm.Event
		_ = json.Unmarshal(inputJson, &event)

		result, err := handleRequest(ctx, event)

		assert.NoError(t, err, "should not error")
		assert.NotNil(t, result, "result should not be nil")
	})
}

func ReadJSONFromFile(t *testing.T, inputFile string) []byte {
	inputJSON, err := os.ReadFile(inputFile)
	if err != nil {
		t.Errorf("could not open test file. details: %v", err)
	}

	return inputJSON
}
