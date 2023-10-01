package aws

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/ikedam/arar/internal/log"
)

type AssumeRoleAndRun struct {
	Region          string
	RoleARN         string
	SessionName     string
	UsernameSession bool
	SerialNumber    string
	DurationSeconds int32
}

type Session struct {
	AssumedRoleUserARN string
	AssumedRoleID      string
	AccessKeyID        string
	SecretAccessKey    string
	SessionToken       string
	Expiration         time.Time
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
	input := &sts.AssumeRoleInput{
		RoleArn:         &a.RoleARN,
		RoleSessionName: &a.SessionName,
	}
	if a.UsernameSession {
		output, err := s.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
		if err != nil {
			return nil, fmt.Errorf("failed to call get caller identity: %w", err)
		}
		log.Debug(
			ctx,
			"get-caller-identity succeeded",
			log.WithString("arn", *output.Arn),
		)
		parts := strings.Split(*output.Arn, "/")
		input.RoleSessionName = &parts[len(parts)-1]
	}
	if a.SerialNumber != "" {
		input.SerialNumber = &a.SerialNumber
		var v string
		log.Debug(
			ctx,
			"get token code",
			log.WithString("serial_number", a.SerialNumber),
		)
		fmt.Fprintf(os.Stderr, "Assume Role MFA token code: ")
		_, err := fmt.Scanln(&v)
		if err != nil {
			return nil, fmt.Errorf("failed to get MFA token code: %w", err)
		}
		input.TokenCode = &v
	}
	if a.DurationSeconds != 0 {
		input.DurationSeconds = &a.DurationSeconds
	}

	log.Debug(
		ctx,
		"assume-role",
		log.WithString("role_arn", *input.RoleArn),
		log.WithString("role_session_name", *input.RoleSessionName),
	)
	output, err := s.AssumeRole(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to call assume role: %w", err)
	}
	log.Debug(
		ctx,
		"assume-role succeeded",
		log.WithString("arn", *output.AssumedRoleUser.Arn),
		log.WithString("role_id", *output.AssumedRoleUser.AssumedRoleId),
		log.WithString("access_key_id", *output.Credentials.AccessKeyId),
		log.WithTime("expiration", *output.Credentials.Expiration),
	)
	return &Session{
		AssumedRoleUserARN: *output.AssumedRoleUser.Arn,
		AssumedRoleID:      *output.AssumedRoleUser.AssumedRoleId,
		AccessKeyID:        *output.Credentials.AccessKeyId,
		SecretAccessKey:    *output.Credentials.SecretAccessKey,
		SessionToken:       *output.Credentials.SessionToken,
		Expiration:         *output.Credentials.Expiration,
	}, nil
}
