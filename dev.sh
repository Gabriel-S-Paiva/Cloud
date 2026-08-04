#!/usr/bin/env bash

# Kill background processes on Ctrl+C
trap 'kill 0' EXIT

echo "Starting Go backend..."
(cd backend && go run .) &

echo "Starting Svelte frontend..."
(cd frontend && pnpm run dev) &

wait
