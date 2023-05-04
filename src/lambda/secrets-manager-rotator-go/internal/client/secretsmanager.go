package client

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"log"
	"time"
)

type SecretsManagerClient struct {
	Client *secretsmanager.Client
}

type SecretsManager interface {
	ListAll() ([]*secretsmanager.DescribeSecretOutput, error)
}

func (s *SecretsManagerClient) ListAll() ([]*secretsmanager.DescribeSecretOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var allOutput []types.SecretListEntry
	var nextToken *string
externalLoop:
	for {
		listOutput, err := s.Client.ListSecrets(
			ctx,
			&secretsmanager.ListSecretsInput{
				//Filters: []types.Filter{
				//	{
				//		Key: "name",
				//		Values: []string{
				//			fmt.Sprintf("service/%s", omdServiceName),
				//		},
				//	},
				//},
				NextToken: nextToken,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("error listing secret values: %w", err)
		}
		log.Printf("[DEBUG] secrets found: %d", len(listOutput.SecretList))
		log.Printf("[DEBUG] pagination: %v", listOutput.NextToken)
		allOutput = append(allOutput, listOutput.SecretList...)
		nextToken = listOutput.NextToken
		if nextToken == nil {
			break externalLoop
		}
	}

	var allSecrets []*secretsmanager.DescribeSecretOutput
	for _, secret := range allOutput {
		secretOutput, err := s.Client.DescribeSecret(
			ctx,
			&secretsmanager.DescribeSecretInput{
				SecretId: secret.ARN,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("error describing secret: %w", err)
		}
		allSecrets = append(allSecrets, secretOutput)
	}

	return allSecrets, nil
}

func NewSecretsManager(cfg aws.Config) *SecretsManagerClient {
	return &SecretsManagerClient{
		Client: secretsmanager.NewFromConfig(cfg),
	}
}
