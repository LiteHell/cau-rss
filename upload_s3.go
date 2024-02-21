package main

import (
	"context"
	"mime"
	"os"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func recursiveReadDir(dir string, prefixes []string, handler func(filename string, key string)) {
	files, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if err != nil {
			panic(err)
		}

		fullpath := path.Join(dir, file.Name())
		keyComponents := append(prefixes, file.Name())
		if file.IsDir() {
			recursiveReadDir(fullpath, keyComponents, handler)
		} else {
			handler(fullpath, strings.Join(keyComponents, "/"))
		}
	}
}

func uploadS3(dir string) {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}

	client := s3.NewFromConfig(cfg)

	recursiveReadDir(dir, make([]string, 0), (func(filename string, key string) {
		reader, err := os.Open(filename)
		if err != nil {
			panic(err)
		}

		mimeType := mime.TypeByExtension(path.Ext(filename))
		if strings.HasSuffix(key, "rss") {
			mimeType = "application/rss+xml"
		} else if strings.HasSuffix(key, "atom") {
			mimeType = "application/atom+xml"
		} else if strings.HasSuffix(key, "feed") {
			mimeType = "application/feed+json"
		} else if mimeType == "" {
			mimeType = "application/octet-stream"
		}

		_, err = client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:      aws.String("rss.puang.network"),
			Key:         aws.String(key),
			Body:        reader,
			ContentType: &mimeType,
			Metadata: map[string]string{
				"Content-Type": mimeType,
			},
		})

		if err != nil {
			panic(err)
		}
	}))

}
