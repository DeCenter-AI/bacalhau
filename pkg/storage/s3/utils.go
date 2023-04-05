package s3

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/rs/zerolog/log"
)

var isEC2Instance bool
var isEC2InstanceOnce sync.Once

func DefaultAWSConfig() (aws.Config, error) {
	imdsMode := imds.ClientEnabled
	if !IsEC2Instance(context.Background()) {
		imdsMode = imds.ClientDisabled
	}
	return config.LoadDefaultConfig(context.Background(),
		config.WithEC2IMDSClientEnableState(imdsMode),
	)
}

// HasValidCredentials returns true if the AWS config has valid credentials.
func HasValidCredentials(config aws.Config) bool {
	credentials, err := config.Credentials.Retrieve(context.Background())
	if err != nil {
		log.Debug().Err(err).Msg("Failed to check if we have valid AWS credentials")
		return false
	}
	return credentials.HasKeys()
}

// IsEC2Instance returns true if the current process is running on an EC2 instance.
// This requires checking the EC2 instance metadata service, which takes a few seconds to resolve. This is why we are only calling it
// once and caching the result.
func IsEC2Instance(ctx context.Context) bool {
	isEC2InstanceOnce.Do(func() {
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			isEC2Instance = false
			return
		}

		client := imds.NewFromConfig(cfg)
		_, err = client.GetMetadata(ctx, &imds.GetMetadataInput{
			Path: "instance-id",
		})
		if err != nil {
			isEC2Instance = false
		}
		isEC2Instance = true
	})
	return isEC2Instance
}