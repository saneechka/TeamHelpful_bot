#!/bin/bash

# Check for .env file
if [ ! -f .env ]; then
    echo "Error: .env file not found!"
    echo "Please create .env file with BOT_TOKEN=your_token_here and DB_PATH=users.db"
    exit 1
fi

# Stop any running containers
echo "Stopping any running containers..."
docker-compose down

# Build and start containers
echo "Building and starting containers..."
docker-compose up -d --build

# Check if container is running
if docker-compose ps | grep -q "Up"; then
    echo "HelpBot successfully deployed!"
else
    echo "Error: HelpBot failed to start. Check logs with 'docker-compose logs'"
    exit 1
fi


chmod +x deploy.sh 