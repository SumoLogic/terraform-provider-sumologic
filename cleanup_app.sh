#!/bin/bash

# Cleanup script for orphaned Varnish app instance
# This script will uninstall the app that's blocking the tests

UUID="d2ef33c3-67f2-4438-9124-14a30ec2ecf3"

# Check required environment variables
if [ -z "$SUMOLOGIC_ACCESSID" ] || [ -z "$SUMOLOGIC_ACCESSKEY" ]; then
    echo "Error: Please set SUMOLOGIC_ACCESSID and SUMOLOGIC_ACCESSKEY environment variables"
    echo ""
    echo "Example:"
    echo "  export SUMOLOGIC_ACCESSID=\"your_access_id\""
    echo "  export SUMOLOGIC_ACCESSKEY=\"your_access_key\""
    echo "  export SUMOLOGIC_ENVIRONMENT=\"us2\"  # or set SUMOLOGIC_BASE_URL"
    exit 1
fi

# Construct base URL
if [ -n "$SUMOLOGIC_BASE_URL" ]; then
    # Remove /api suffix if present, and any trailing slashes
    BASE_URL="${SUMOLOGIC_BASE_URL%/api}"
    BASE_URL="${BASE_URL%/}"
elif [ -n "$SUMOLOGIC_ENVIRONMENT" ]; then
    BASE_URL="https://api.${SUMOLOGIC_ENVIRONMENT}.sumologic.com/api"
else
    echo "Error: Either SUMOLOGIC_BASE_URL or SUMOLOGIC_ENVIRONMENT must be set"
    exit 1
fi

echo "Using Base URL: $BASE_URL"
echo "Uninstalling app with UUID: $UUID"
echo ""

# Call the uninstall API
RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X POST \
  -u "${SUMOLOGIC_ACCESSID}:${SUMOLOGIC_ACCESSKEY}" \
  -H "Content-Type: application/json" \
  "${BASE_URL}/v2/apps/${UUID}/uninstall")

# Extract HTTP status and response body
HTTP_BODY=$(echo "$RESPONSE" | sed -e 's/HTTP_STATUS\:.*//g')
HTTP_STATUS=$(echo "$RESPONSE" | tr -d '\n' | sed -e 's/.*HTTP_STATUS://')

echo "HTTP Status: $HTTP_STATUS"
echo "Response: $HTTP_BODY"
echo ""

if [ "$HTTP_STATUS" = "200" ] || [ "$HTTP_STATUS" = "202" ]; then
    if echo "$HTTP_BODY" | grep -q "jobId"; then
        echo "✓ App uninstall job started successfully"
        echo "Waiting 15 seconds for job to complete..."
        sleep 15
        echo "✓ Cleanup complete. You can now run the tests."
    else
        echo "⚠ Uninstall initiated but no jobId in response"
    fi
elif [ "$HTTP_STATUS" = "404" ]; then
    echo "✓ App is already uninstalled (404 Not Found)"
else
    echo "✗ Uninstall failed with HTTP status $HTTP_STATUS"
    echo "Please check your credentials and try again, or manually uninstall via the UI"
fi
