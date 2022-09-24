package main

import (
	"bytes"
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
)

type S3Info struct {
	AwsS3Region  string `yaml:"aws_s3_region"`
	AwsAccessKey string `yaml:"aws_access_key"`
	AwsSecretKey string `yaml:"aws_secret_key"`
	BucketName   string `yaml:"bucket_name"`
	S3Client     *s3.Client
}

var s3Client S3Info

func parseFromConfig() S3Info {
	file, err := ioutil.ReadFile("./resources/properties.yaml")
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = yaml.Unmarshal(file, &s3Client)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return s3Client
}

func init() {
	parseFromConfig()
	err := s3Client.SetS3ConfigByKey()
	if err != nil {
		return
	}

}

func main() {
	reader, err := ioutil.ReadFile("./resources/INfo.png")
	if err != nil {
		log.Fatalln(err.Error())
	}
	s3Client.UploadFile(bytes.NewReader(reader), "Info.png", "")

}

func (s *S3Info) SetS3ConfigByKey() error {
	creds := credentials.NewStaticCredentialsProvider(s.AwsAccessKey, s.AwsSecretKey, "")
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithCredentialsProvider(creds),
		config.WithRegion(s.AwsS3Region),
	)
	if err != nil {
		log.Printf("error: %v", err)
		panic(err)
		return errors.New(err.Error())
	}
	s.S3Client = s3.NewFromConfig(cfg)
	return nil
}

func (s *S3Info) UploadFile(file io.Reader, filename, preFix string) *manager.UploadOutput {
	uploader := manager.NewUploader(s.S3Client)
	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(filename),
		Body:   file,
	})
	if err != nil {
		log.Fatal(err)
	}
	return result
}
