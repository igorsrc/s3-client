package fs

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
)

type Client struct {
	S3client *s3.Client
}

var client Client

func GetS3Client() *s3.Client {
	return client.S3client
}

func InitClient(address, access, secret, region string) {
	log.Printf("Initializing s3 app connection: %s", address)

	s3Config, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(access, secret, "")),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL: "http://" + address,
				}, nil
			})))

	if err != nil {
		panic("Error configuring fs client: " + err.Error())
	}

	client.S3client = s3.NewFromConfig(s3Config, func(o *s3.Options) { o.UsePathStyle = true })
}

func InitBucket(name string) {

	listBucketsResult, err := client.S3client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})

	for _, bucket := range listBucketsResult.Buckets {
		if *bucket.Name == name {
			log.Printf("Bucket %s already exists. Skip bucket initialization.", name)
			return
		}
	}

	_, err = client.S3client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(name),
	})

	if err != nil {
		panic("Failed to create bucket: " + err.Error())
	}
}
