package s3

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"path"

	"github.com/jcfug8/daylear/server/core/file"
	pConfig "github.com/jcfug8/daylear/server/ports/config"
	"github.com/jcfug8/daylear/server/ports/fileretriever"
	"github.com/rs/zerolog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/jcfug8/daylear/server/core/logutil"
	"go.uber.org/fx"
)

type Client struct {
	log zerolog.Logger

	s3Client          *s3.Client
	presignClient     *s3.PresignClient
	bucket            string
	publicEndpointURL *url.URL
}

type NewClientParams struct {
	fx.In

	L            zerolog.Logger
	ConfigClient pConfig.Client
}

func NewClient(params NewClientParams) (*Client, error) {
	ctx := context.Background()

	c := params.ConfigClient.GetConfig()["s3"].(map[string]interface{})
	accessKey := c["accesskey"].(string)
	secret := c["secret"].(string)
	endpoint := c["endpoint"].(string)
	publicEndpoint := c["publicendpoint"].(string)
	bucket := c["bucket"].(string)

	publicEndpointURL, err := url.Parse(publicEndpoint)
	if err != nil {
		return nil, fileretriever.ErrInvalidArgument{Msg: "unable to parse public endpoint"}
	}

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secret, "")),
	)
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})
	presignClient := s3.NewPresignClient(s3Client)

	client := &Client{
		log:               params.L,
		s3Client:          s3Client,
		presignClient:     presignClient,
		bucket:            bucket,
		publicEndpointURL: publicEndpointURL,
	}

	// Ensure the bucket exists
	if err := client.ensureBucketExists(ctx); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) ensureBucketExists(ctx context.Context) error {
	log := logutil.EnrichLoggerWithContext(c.log, ctx)
	log.Info().Msgf("Ensuring bucket %s exists", c.bucket)
	_, err := c.s3Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(c.bucket),
	})
	if err != nil {
		log.Info().Msgf("Bucket %s does not exist, creating", c.bucket)
		_, err := c.s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(c.bucket),
		})
		if err != nil {
			log.Error().Err(err).Msg("failed to create bucket")
			return fmt.Errorf("failed to create bucket: %w", err)
		}
	}
	return nil
}

func (c *Client) UploadPublicFile(ctx context.Context, filePath string, file file.File) (string, error) {
	log := logutil.EnrichLoggerWithContext(c.log, ctx)
	_, err := c.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(c.bucket),
		Key:           aws.String(filePath),
		Body:          file,
		ContentType:   aws.String(file.ContentType),
		ContentLength: aws.Int64(file.ContentLength),
		ACL:           types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to upload public file")
	}
	publicEndpointURL := *c.publicEndpointURL
	publicEndpointURL.Path = path.Join(c.bucket, filePath)
	return publicEndpointURL.String(), err
}

func (c *Client) GetFile(ctx context.Context, filePath string) (io.ReadCloser, error) {
	log := logutil.EnrichLoggerWithContext(c.log, ctx)
	resp, err := c.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(filePath),
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to get file")
		return nil, err
	}
	return resp.Body, nil
}

func (c *Client) DeleteFile(ctx context.Context, path string) error {
	log := logutil.EnrichLoggerWithContext(c.log, ctx)
	_, err := c.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to delete file")
	}
	return err
}
