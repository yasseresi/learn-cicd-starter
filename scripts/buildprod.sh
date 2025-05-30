#!/bin/bash

echo "🔧 Building Go binary..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o notely

# Authenticate with Google Cloud
echo "🔐 Setting up gcloud CLI..."
# Make sure gcloud is installed on your machine, or install it beforehand
gcloud auth login
gcloud config set project notely-461315 # Replace with your actual project ID
gcloud config set artifacts/location us-central1 # e.g., us-central1

# Enable Artifact Registry if not already enabled
gcloud services enable artifactregistry.googleapis.com

# Docker auth to Artifact Registry
echo "🔑 Configuring Docker to use Artifact Registry..."
gcloud auth configure-docker us-central1-docker.pkg.dev

# Build and push Docker image
echo "🐳 Building and pushing Docker image..."
gcloud builds submit --tag us-central1-docker.pkg.dev/notely-461315/notely-ar-repo/notely:latest .

echo "✅ Build and deploy completed."
