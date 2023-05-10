package rotation

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"github.com/excoriate/aws-secrets-rotation-lambda/internal/client"
	"github.com/excoriate/aws-secrets-rotation-lambda/internal/erroer"
	"go.uber.org/zap"
)

type StepExecutioner interface {
	CreateSecretStep() error
	SetSecretStep() error
	TestSecretStep() error
	FinishSecretStep() error
}

type StepsClient struct {
	Logger *zap.Logger
	Client client.SecretsManager
	// Common in all the step functions, and thus, it's convenient to have it
	// here, set at the object/struct's state.
	SecretEvent   Event
	SecretData    *secretsmanager.DescribeSecretOutput
	StagingLabels StagingLabels
}

func (s *StepsClient) CreateSecretStep() error {
	secretId := *s.SecretData.ARN
	token := *s.SecretEvent.Token
	stagePending := s.StagingLabels.Pending

	// This step tries to retrieve the secret version with the AWSPENDING stage label
	// associated with the given token. If the secret version is found,
	// it means that the AWSPENDING version is already created,
	// and there is no need to create a new version. In this case,
	// the rotation process can proceed to the next step.
	//
	// However, if the secret version with the AWSPENDING stage label is not found,
	// it indicates that the secret version with this stage label needs to be created.
	// In this case, the rotation process will create a new secret version,
	// set its stage label to AWSPENDING, and proceed with the rest of the rotation steps.
	//
	// This approach ensures that a new secret version is created only when needed,
	//avoiding unnecessary secret version creations and maintaining the integrity of the rotation process.
	_, err := s.Client.GetSecretValueByStageLabel(secretId, token, stagePending)
	if err != nil {
		var resourceNotFoundError *types.ResourceNotFoundException
		// If the secret isn't found, that's fine, Let's create a new secret version then.
		excChars := "/@'\"\\"
		if errors.As(err, &resourceNotFoundError) {
			newSecretValue, err := s.Client.GenerateRandomPassword(excChars)
			if err != nil {
				s.Logger.Error("Error generating random password", zap.Error(err))
				return erroer.NewRotationError("Error generating random password", err)
			}

			// Create a new secret version, with the new rotated value.
			_, err = s.Client.PutSecretValue(secretId, token, newSecretValue, stagePending)
			if err != nil {
				s.Logger.Error("Error creating new secret version", zap.Error(err))
				return erroer.NewRotationError("Error creating new secret version", err)
			}

		}
	}

	return nil
}

func (s *StepsClient) SetSecretStep() error {
	return nil
}

func (s *StepsClient) TestSecretStep() error {
	return nil
}

func (s *StepsClient) FinishSecretStep() error {
	arn := *s.SecretData.ARN
	token := *s.SecretEvent.Token
	s.Logger.Info("Finishing secret rotation", zap.String("ARN", arn))

	secret, err := s.Client.GetSecret(arn)
	if err != nil {
		s.Logger.Error(fmt.Sprintf("Error describing secret with arn %s", arn), zap.Error(err))
		return erroer.NewRotationError(fmt.Sprintf("Error describing secret with arn %s", arn), err)
	}

	var currentVersion string
	// Check versions
	for version, _ := range secret.VersionIdsToStages {
		if secret.VersionIdsToStages[version] != nil {
			if version == token {
				// The correct version is already marked as current, no need to do anything.
				s.Logger.Info("The correct version %s is already marked as current, no need to do anything.", zap.String("version", version))
				return nil
			}

			s.Logger.Info("Setting version %s as current", zap.String("version", version))
			currentVersion = version
			break
		}
	}

	// Finalize by staging the secret version current
	currentStage := s.StagingLabels.Current
	updatedSecret, finErr := s.Client.UpdateSecretVersion(arn, token, currentStage, currentVersion)

	if finErr != nil {
		s.Logger.Error(fmt.Sprintf("Error finalizing secret rotation with arn %s, "+
			"while updating the secret with token %s to stage %s and version %s", arn, token,
			currentStage, currentVersion), zap.Error(finErr))

		return erroer.NewRotationError(fmt.Sprintf("Error finalizing secret rotation with arn %s,"+
			" "+"while updating the secret with token %s to stage %s and version %s", arn, token,
			currentStage, currentVersion), finErr)
	}

	s.Logger.Info("Secret rotation finished successfully for secret with ARN %s, and token %s", zap.String("ARN", *updatedSecret.ARN), zap.String("token", token))
	return nil
}

func NewStepExecutionerClient(logger *zap.Logger, client client.SecretsManager,
	secretEvent Event, secretData *secretsmanager.DescribeSecretOutput) *StepsClient {
	return &StepsClient{
		Logger:        logger,
		Client:        client,
		SecretEvent:   secretEvent,
		SecretData:    secretData,
		StagingLabels: GetStagingLabels(),
	}
}
