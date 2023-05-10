package rotation

import (
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
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
	pendingLabel := s.StagingLabels.Pending

	// This step tries to retrieve the secret version with the AWSPENDING stage label
	// associated with the given token. If the secret version is found,
	// it means that the AWSPENDING version is already created,
	// and there is no need to create a new version. In this case,
	//the rotation process can proceed to the next step.
	//
	// However, if the secret version with the AWSPENDING stage label is not found,
	// it indicates that the secret version with this stage label needs to be created.
	// In this case, the rotation process will create a new secret version,
	// set its stage label to AWSPENDING, and proceed with the rest of the rotation steps.
	//
	// This approach ensures that a new secret version is created only when needed,
	//avoiding unnecessary secret version creations and maintaining the integrity of the rotation process.
	_, err := s.Client.GetSecretValueByStageLabel(secretId, token, pendingLabel)
	if err != nil {
		// If the secret isn't found, that's fine, Let's create a new secret version then.
		excChars := "/@'\"\\"
		if err.Error() == "ResourceNotFoundException: Secrets Manager can't find the specified secret." {
			newSecretValue, err := s.Client.GenerateRandomPassword(excChars)
			if err != nil {
				s.Logger.Error("Error generating random password", zap.Error(err))
				return erroer.NewRotationError("Error generating random password", err)
			}

			// Create a new secret version, with the new rotated value.
			_, err = s.Client.PutSecretValue(secretId, newSecretValue, token, pendingLabel)
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