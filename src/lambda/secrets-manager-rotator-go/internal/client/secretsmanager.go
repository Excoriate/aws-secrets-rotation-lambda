package client

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"go.uber.org/zap"
	"time"
)

type SecretsManagerClient struct {
	Client *secretsmanager.Client
	Logger *zap.Logger
}

type SecretsManager interface {
	ListAll() ([]*secretsmanager.DescribeSecretOutput, error)
	GetSecret(arn string) (*secretsmanager.DescribeSecretOutput, error)
	GetSecretValue(arn, token, stage string) (*secretsmanager.GetSecretValueOutput, error)
	GetSecretValueByStageLabel(arn, token, stageLabel string) (*secretsmanager.
	GetSecretValueOutput,
		error)
	PutSecretValue(arn, token, value, stage string) (*secretsmanager.PutSecretValueOutput, error)
	GenerateRandomPassword(excludeChars string) (string, error)
	UpdateSecretVersion(arn, token, stage, currentVersion string) (*secretsmanager.UpdateSecretVersionStageOutput, error)
}

func (s *SecretsManagerClient) UpdateSecretVersion(arn, token, stage,
	currentVersion string) (*secretsmanager.UpdateSecretVersionStageOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	secretVersionOutput, err := s.Client.UpdateSecretVersionStage(
		ctx,
		&secretsmanager.UpdateSecretVersionStageInput{
			SecretId:            aws.String(arn),
			VersionStage:        aws.String(stage),
			MoveToVersionId:     aws.String(token),
			RemoveFromVersionId: aws.String(currentVersion),
		},
	)

	if err != nil {
		s.Logger.Error("error updating secret version stage", zap.Error(err))
		return nil, fmt.Errorf("error updating secret version stage: %w", err)
	}

	return secretVersionOutput, nil
}

func (s *SecretsManagerClient) GenerateRandomPassword(excludeChars string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	passwordOutput, err := s.Client.GetRandomPassword(
		ctx,
		&secretsmanager.GetRandomPasswordInput{
			ExcludeCharacters: aws.String(excludeChars),
		},
	)

	if err != nil {
		s.Logger.Error("error generating random password", zap.Error(err))
		return "", fmt.Errorf("error generating random password: %w", err)
	}

	return *passwordOutput.RandomPassword, nil
}

func (s *SecretsManagerClient) GetSecretValueByStageLabel(arn,
	token, stageLabel string) (*secretsmanager.GetSecretValueOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var secretValueOutput *secretsmanager.GetSecretValueOutput
	var err error
	if token == "" {
		secretValueOutput, err = s.Client.GetSecretValue(
			ctx,
			&secretsmanager.GetSecretValueInput{
				SecretId:     aws.String(arn),
				VersionStage: aws.String(stageLabel),
			},
		)
	} else {
		secretValueOutput, err = s.Client.GetSecretValue(
			ctx,
			&secretsmanager.GetSecretValueInput{
				SecretId:     aws.String(arn),
				VersionStage: aws.String(stageLabel),
				VersionId:    aws.String(token),
			},
		)
	}

	if err != nil {
		s.Logger.Error("error getting secret value", zap.Error(err))
		return nil, fmt.Errorf("error getting secret value: %w", err)
	}

	return secretValueOutput, nil
}

func (s *SecretsManagerClient) PutSecretValue(arn, token, value,
	stage string) (*secretsmanager.PutSecretValueOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	secretValueOutput, err := s.Client.PutSecretValue(
		ctx,
		&secretsmanager.PutSecretValueInput{
			SecretId:           aws.String(arn),
			ClientRequestToken: aws.String(token),
			SecretString:       aws.String(value),
			VersionStages:      []string{stage},
		},
	)

	if err != nil {
		s.Logger.Error("error putting secret value", zap.Error(err))
		return nil, fmt.Errorf("error putting secret value: %w", err)
	}

	return secretValueOutput, nil
}

func (s *SecretsManagerClient) GetSecret(arn string) (*secretsmanager.DescribeSecretOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	secretOutput, err := s.Client.DescribeSecret(
		ctx,
		&secretsmanager.DescribeSecretInput{
			SecretId: aws.String(arn),
		},
	)

	if err != nil {
		s.Logger.Error("error describing secret", zap.Error(err))
		return nil, fmt.Errorf("error describing secret: %w", err)
	}

	return secretOutput, nil
}

func (s *SecretsManagerClient) GetSecretValue(arn, token, stage string) (*secretsmanager.GetSecretValueOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	secretValueOutput, err := s.Client.GetSecretValue(
		ctx,
		&secretsmanager.GetSecretValueInput{
			SecretId:     aws.String(arn),
			VersionStage: aws.String(stage),
			VersionId:    aws.String(token),
		},
	)

	if err != nil {
		s.Logger.Error("error getting secret value", zap.Error(err))
		return nil, fmt.Errorf("error getting secret value: %w", err)
	}

	return secretValueOutput, nil
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
				NextToken: nextToken,
			},
		)
		if err != nil {
			s.Logger.Error("error listing secret values", zap.Error(err))
			return nil, fmt.Errorf("error listing secret values: %w", err)
		}

		s.Logger.Info(fmt.Sprintf("secrets found: %d", len(listOutput.SecretList)))
		s.Logger.Debug(fmt.Sprintf("secrets: %v", listOutput.SecretList))

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
			s.Logger.Error("error describing secret", zap.Error(err))
			return nil, fmt.Errorf("error describing secret: %w", err)
		}
		allSecrets = append(allSecrets, secretOutput)
	}

	return allSecrets, nil
}

func NewSecretsManager(cfg aws.Config, logger *zap.Logger) SecretsManager {
	return &SecretsManagerClient{
		Client: secretsmanager.NewFromConfig(cfg),
		Logger: logger,
	}
}
