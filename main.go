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
	"net/http"
	"path/filepath"
	"strings"
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
	//reader, err := ioutil.ReadFile("./resources/INfo.png")
	//if err != nil {
	//	log.Fatalln(err.Error())
	//}
	res, err := http.Get("https://postfiles.pstatic.net/MjAyMjA0MDRfMjAy/MDAxNjQ5MDY2MDY3MTUw.22jLgHl_oh3zJtif3QuD4qKaz96MMq8AQ6VBxW_d_Jkg.vkHSQE7PR_6Q4zkV7jxXsFo1ct79V0TRwZIxl0TdTl4g.JPEG.dsmhs2022/02.jpg?type=w773")
	if err != nil {
		panic(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)

	reader, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	s3Client.UploadFile(bytes.NewReader(reader), "test.png", "")

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
	fileType := "image/" + findExtension(filename)
	uploader := manager.NewUploader(s.S3Client)
	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s.BucketName),
		Key:         aws.String(filename),
		Body:        file,
		ContentType: &fileType,
	})
	if err != nil {
		log.Fatal(err)
	}
	println(result.Location)
	return result
}

func findExtension(path string) string {
	ext := filepath.Ext(path)
	_, result, _ := strings.Cut(ext, ".")
	return result
}
