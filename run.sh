#!/bin/bash

echo "Starting Backend (Go) and Frontend (Next.js)..."

# Trap SIGINT to kill both processes on exit
trap 'kill $(jobs -p)' SIGINT EXIT

# Start Backend
echo "Starting Go backend..."
cd backend
go run main.go &
cd ..

# Start Frontend
echo "Starting Next.js frontend..."
cd frontend-app
npm run dev -- -H 0.0.0.0 &
cd ..

echo "Both services are running! Press Ctrl+C to stop."
wait
