package adapter

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func NewAWS(region string) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(lo *config.LoadOptions) error {
		lo.DefaultRegion = region
		return nil
	})
	if err != nil {
		return aws.Config{}, err
	}

	return cfg, nil
}
