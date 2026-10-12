import boto3
import os
import urllib.parse

ecs = boto3.client("ecs")

def lambda_handler(event, context):
    bucket = event["detail"]["bucket"]["name"]
    key = urllib.parse.unquote_plus(event["detail"]["object"]["key"])

    print(f"New upload: s3://{bucket}/{key}")

    response = ecs.run_task(
        cluster=os.environ["ECS_CLUSTER"],
        taskDefinition=os.environ["TASK_DEFINITION"],
        launchType="FARGATE",
        networkConfiguration={
            "awsvpcConfiguration": {
                "subnets": [os.environ["SUBNET_ID"]],
                "securityGroups": [os.environ["SECURITY_GROUP_ID"]],
                "assignPublicIp": "ENABLED"
            }
        },
        overrides={
            "containerOverrides": [
                {
                    "name": "video-processor",
                    "environment": [
                        {
                            "name": "INPUT_KEY",
                            "value": key
                        }
                    ]
                }
            ]
        }
    )

    print(f"Started ECS task: {response['tasks'][0]['taskArn']}")

    return {
        "statusCode": 200,
        "inputKey": key
    }
