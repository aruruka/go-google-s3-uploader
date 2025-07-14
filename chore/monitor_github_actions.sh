#!/bin/bash

# This script monitors the latest GitHub Actions workflow for the current repository.
# It requires the 'gh' CLI tool to be installed and authenticated.

REPO=$(basename "$(git rev-parse --show-toplevel)")
OWNER=$(basename "$(dirname "$(git rev-parse --show-toplevel)")") # Assuming owner is the parent directory name

echo "Monitoring latest GitHub Actions workflow for $OWNER/$REPO..."

# Get the latest workflow run ID for the main branch
RUN_ID=$(gh run list --branch main --json databaseId --jq '.[0].databaseId')

if [ -z "$RUN_ID" ]; then
  echo "No workflow runs found on the main branch. Please ensure a workflow has been triggered."
  exit 1
fi

echo "Found latest workflow run ID: $RUN_ID"

STATUS=""
CONCLUSION=""

while [[ "$STATUS" != "completed" ]]; do
  RUN_INFO=$(gh run view "$RUN_ID" --json status,conclusion --jq '{status: .status, conclusion: .conclusion}')
  STATUS=$(echo "$RUN_INFO" | jq -r '.status')
  CONCLUSION=$(echo "$RUN_INFO" | jq -r '.conclusion')

  echo "Workflow Status: $STATUS, Conclusion: $CONCLUSION (Run ID: $RUN_ID)"

  if [[ "$STATUS" == "completed" ]]; then
    if [[ "$CONCLUSION" == "success" ]]; then
      echo "✅ GitHub Actions workflow completed successfully!"
    else
      echo "❌ GitHub Actions workflow completed with status: $CONCLUSION"
    fi
    break
  fi

  sleep 10 # Wait for 10 seconds before checking again
done
