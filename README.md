# Cloud Video Processing Pipeline

A serverless video processing pipeline built on AWS.

Upload a video to S3 and the pipeline automatically starts an ECS Fargate task that transcodes it to 720p H.264 using FFmpeg, then uploads the processed video back to S3.

## Architecture

S3 Upload
→ EventBridge
→ Lambda
→ ECS Fargate
→ FFmpeg
→ S3 Processed Output

## Tech

- Go
- Python
- Docker
- AWS S3
- AWS Lambda
- Amazon ECS / Fargate
- Amazon ECR
- EventBridge
- CloudWatch
- GitHub Actions

## How it works

1. A video is uploaded to the `uploads/` prefix in S3.
2. EventBridge detects the upload.
3. A Lambda function starts an ECS Fargate task.
4. The container downloads the video from S3.
5. FFmpeg transcodes it to 720p H.264.
6. The processed video is uploaded to `processed/`.
7. Logs are sent to CloudWatch.

## CI/CD

Pushes to `main` trigger a GitHub Actions workflow that:

- builds the Docker image
- authenticates to AWS using OIDC
- pushes the image to Amazon ECR

No long-lived AWS credentials are stored in GitHub.

## Running locally

Build the image:

docker build -t video-processor .

Run it with the required AWS credentials and environment variables configured.

## Project Structure

- `main.go` - video processing application
- `lambda_function.py` - starts the Fargate task
- `Dockerfile` - container image
- `.github/workflows/deploy.yml` - CI/CD workflow
- `event-pattern.json` - EventBridge configuration