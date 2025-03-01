#!/bin/bash

# Проверка наличия .env файла
if [ ! -f .env ]; then
    echo "Error: .env file not found!"
    echo "Please create .env file with BOT_TOKEN=your_token_here"
    exit 1
fi

# Сборка и запуск контейнера
docker-compose up -d --build

echo "HelpBot successfully deployed!"

chmod +x deploy.sh 