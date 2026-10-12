package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	bucket = "cloud-video-transcoder"

	localInput  = "input/video.mp4"
	localOutput = "output/video.mp4"
)

func main() {
	// Create local working directories
	err := os.MkdirAll("input", 0755)
	if err != nil {
		fmt.Println("Could not create input directory:", err)
		return
	}

	err = os.MkdirAll("output", 0755)
	if err != nil {
		fmt.Println("Could not create output directory:", err)
		return
	}

	inputKey := os.Getenv("INPUT_KEY")
	if inputKey == "" {
		inputKey = "uploads/video.mp4"
	}

	filename := filepath.Base(inputKey)
	outputKey := "processed/" + filename

	// 1. Download ONE video from S3
	fmt.Println("Downloading from S3...")

	err = downloadFromS3(bucket, inputKey, localInput)
	if err != nil {
		fmt.Println("Download failed:", err)
		return
	}

	fmt.Println("Downloaded:", inputKey)

	// 2. Transcode that ONE video
	fmt.Println("\nTranscoding video...")

	cmd := exec.Command(
		"ffmpeg",
		"-y",
		"-loglevel", "error",
		"-i", localInput,
		"-vf", "scale=-2:720",
		"-c:v", "libx264",
		"-crf", "23",
		localOutput,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Println("Transcoding failed:", err)
		return
	}

	fmt.Println("Transcoding complete:", localOutput)

	// 3. Upload processed video back to S3
	fmt.Println("\nUploading processed video to S3...")

	err = uploadToS3(bucket, outputKey, localOutput)
	if err != nil {
		fmt.Println("Upload failed:", err)
		return
	}

	fmt.Println("Uploaded:", outputKey)
	fmt.Println("\nPipeline complete!")
}

func downloadFromS3(bucket string, key string, destination string) error {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	client := s3.NewFromConfig(cfg)

	result, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return err
	}
	defer result.Body.Close()

	file, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, result.Body)
	return err
}

func uploadToS3(bucket string, key string, source string) error {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	client := s3.NewFromConfig(cfg)

	file, err := os.Open(source)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   file,
	})

	return err
}
