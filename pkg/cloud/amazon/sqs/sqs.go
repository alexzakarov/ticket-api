package sqs

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// NewAwsSqs Return new AWS Simple Queue System instance
func NewAwsSqs(ctx context.Context, region string) *sqs.Client {

	cfg_sqs, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := sqs.NewFromConfig(cfg_sqs)
	println("SQS Client Initialized")

	return client
}
