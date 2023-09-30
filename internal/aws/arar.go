package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type AssumeRoleAndRun struct {
	Region          string
	RoleARN         string
	SessionName     string
	UsernameSession bool
}

type Session struct {
	AssumedRoleUserARN string
	AssumedRoleID      string
	AccessKeyID        string
	SecretAccessKey    string
	SessionToken       string
}

func (a *AssumeRoleAndRun) AssumeRole(ctx context.Context) (*Session, error) {
	var cfgOpts []func(*config.LoadOptions) error
	if a.Region != "" {
		cfgOpts = append(cfgOpts, config.WithRegion(a.Region))
	}
	cfg, err := config.LoadDefaultConfig(
		ctx,
		cfgOpts...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load default config: %w", err)
	}
	s := sts.NewFromConfig(cfg)
	sessionName := a.SessionName
	if a.UsernameSession {
		output, err := s.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
		if err != nil {
			return nil, fmt.Errorf("failed to call get caller identity: %w", err)
		}
		sessionName = *output.UserId
	}
	output, err := s.AssumeRole(ctx, &sts.AssumeRoleInput{
		RoleArn:         &a.RoleARN,
		RoleSessionName: &sessionName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call assume role: %w", err)
	}
	return &Session{
		AssumedRoleUserARN: *output.AssumedRoleUser.Arn,
		AssumedRoleID:      *output.AssumedRoleUser.AssumedRoleId,
		AccessKeyID:        *output.Credentials.AccessKeyId,
		SecretAccessKey:    *output.Credentials.SecretAccessKey,
		SessionToken:       *output.Credentials.SessionToken,
	}, nil
}
