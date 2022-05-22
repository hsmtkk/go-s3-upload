package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/hsmtkk/go-s3-upload/upload"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func main() {
	mustEnvVar()
	cmd := &cobra.Command{
		Use:  "go-s3-upload bucket directory",
		Run:  run,
		Args: cobra.ExactArgs(2),
	}
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func run(cmd *cobra.Command, args []string) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("failed to init logger; %v", err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	bucket := args[0]
	directory := args[1]
	uploader := upload.New(bucket, directory)
	entries, err := ioutil.ReadDir(directory)
	if err != nil {
		sugar.Fatalw("failed to read directory", "directory", directory, "error", err)
	}
	for _, entry := range entries {
		name := entry.Name()
		sugar.Infow("start upload", "name", name)
		location, err := uploader.Upload(name)
		if err != nil {
			sugar.Error(err)
		} else {
			sugar.Infow("finish upload", "location", location)
		}
	}
}

func mustEnvVar() {
	envs := []string{"AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "AWS_REGION"}
	for _, e := range envs {
		v := os.Getenv(e)
		if v == "" {
			log.Fatalf("you must define %s env var", e)
		}
	}
}
