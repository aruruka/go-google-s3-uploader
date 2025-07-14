#!/bin/sh

# Start the auth-server in the background
echo "Starting auth-server..."
./auth-server/auth-server-main &

# Start the app-server in the background
echo "Starting app-server..."
./app-server/app-server-main &

# Wait for all background processes to complete
wait
