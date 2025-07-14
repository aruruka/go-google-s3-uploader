#!/bin/bash

# This script monitors the latest AWS App Runner deployment for a specific service.
# It requires the 'aws' CLI tool to be installed and configured with the 'local-dev-admin' profile.

SERVICE_NAME="go-s3-uploader-service" # Replace with your App Runner service name
AWS_REGION="ap-northeast-1" # Replace with your AWS region
AWS_PROFILE="local-dev-admin" # Your AWS CLI profile

echo "Monitoring AWS App Runner service: $SERVICE_NAME in region $AWS_REGION with profile $AWS_PROFILE"

# Get the service ARN
SERVICE_ARN=$(aws apprunner list-services --region "$AWS_REGION" --profile "$AWS_PROFILE" \
  --query "ServiceSummaryList[?ServiceName=='$SERVICE_NAME'].ServiceArn" --output text)

if [ -z "$SERVICE_ARN" ]; then
  echo "Error: App Runner service '$SERVICE_NAME' not found in region '$AWS_REGION'."
  exit 1
fi

echo "Service ARN: $SERVICE_ARN"

STATUS=""
OPERATION_ID=""

# Loop to find the latest operation ID (deployment)
echo "Attempting to find the latest deployment operation..."
for i in {1..10}; do
  OPERATION_ID=$(aws apprunner list-operations --service-arn "$SERVICE_ARN" --region "$AWS_REGION" --profile "$AWS_PROFILE" \
    --query "OperationSummaryList[?Type=='UPDATE_SERVICE' || Type=='CREATE_SERVICE'].OperationId" --output text | head -n 1)
  
  if [ -n "$OPERATION_ID" ]; then
    echo "Found latest operation ID: $OPERATION_ID"
    break
  fi
  echo "No active deployment found yet, retrying in 5 seconds..."
  sleep 5
done

if [ -z "$OPERATION_ID" ]; then
  echo "Error: Could not find any recent deployment operations for service '$SERVICE_NAME'."
  exit 1
fi

# Monitor the operation status
while [[ "$STATUS" != "SUCCEEDED" && "$STATUS" != "FAILED" && "$STATUS" != "ABORTED" ]]; do
  OPERATION_INFO=$(aws apprunner describe-operation --operation-id "$OPERATION_ID" --service-arn "$SERVICE_ARN" \
    --region "$AWS_REGION" --profile "$AWS_PROFILE" \
    --query '{Status: Operation.Status, Type: Operation.Type, StartTime: Operation.StartedAt, EndTime: Operation.EndedAt}')
  
  STATUS=$(echo "$OPERATION_INFO" | jq -r '.Status')
  TYPE=$(echo "$OPERATION_INFO" | jq -r '.Type')
  START_TIME=$(echo "$OPERATION_INFO" | jq -r '.StartTime')
  END_TIME=$(echo "$OPERATION_INFO" | jq -r '.EndTime')

  echo "Operation Type: $TYPE, Status: $STATUS (Operation ID: $OPERATION_ID)"

  if [[ "$STATUS" == "SUCCEEDED" ]]; then
    echo "✅ App Runner deployment completed successfully! Start Time: $START_TIME, End Time: $END_TIME"
    break
  elif [[ "$STATUS" == "FAILED" || "$STATUS" == "ABORTED" ]]; then
    echo "❌ App Runner deployment failed or was aborted! Status: $STATUS, Start Time: $START_TIME, End Time: $END_TIME"
    exit 1
  fi

  sleep 15 # Wait for 15 seconds before checking again
done
