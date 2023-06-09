package daggerio

import (
	"context"
	"dagger.io/dagger"
	"os"
)

func NewClient(ctx context.Context) (*dagger.Client, error) {
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewClientWithWorkDir(ctx context.Context, workDir string) (*dagger.Client, error) {
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout), dagger.WithWorkdir(workDir))
	if err != nil {
		return nil, err
	}
	return client, nil
}
