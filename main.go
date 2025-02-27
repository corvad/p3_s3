package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()
	app.OnBootstrap().BindFunc(func(e *core.BootstrapEvent) error {
		if err := e.Next(); err != nil {
			return err
		}
		if e.App.Settings().S3.Enabled {
			e.App.Logger().Info("S3 Enabled: Using Direct CDN for Downloads.")
			cfg, err := config.LoadDefaultConfig(context.TODO(),
				config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(e.App.Settings().S3.AccessKey, e.App.Settings().S3.Secret, "")),
				config.WithRegion(e.App.Settings().S3.Region),
			)
			client := s3.NewFromConfig(cfg, func(o *s3.Options) {
				o.BaseEndpoint = aws.String(e.App.Settings().S3.Endpoint)
			})
			e.App.Store().Set("s3", s3.NewPresignClient(client))

			if err != nil {
				log.Fatal(err)
			}
		}

		return nil
	})

	app.OnFileDownloadRequest().BindFunc(func(e *core.FileDownloadRequestEvent) error {
		if e.App.Settings().S3.Enabled {
			presignResult, err := e.App.Store().Get("s3").(*s3.PresignClient).PresignGetObject(context.TODO(), &s3.GetObjectInput{
				Bucket: aws.String(e.App.Settings().S3.Bucket),
				Key:    aws.String(e.ServedPath),
			})
			if err != nil {
				e.App.Logger().Error("Error Signing S3 URL")
			}
			e.App.Logger().Info("Generated Direct Pre-Signed S3 URL")
			return e.Redirect(302, presignResult.URL)
		}
		return e.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
