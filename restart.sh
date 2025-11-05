#!/bin/bash

COMPOSE_FILE="docker-compose.yml"

echo "--- 🛑 Menghentikan stack Docker Compose yang sedang berjalan... ---"
docker compose down -v

if [ $? -ne 0 ]; then
    echo "❌ Gagal menghentikan stack Docker Compose. Mohon periksa error di atas."
    exit 1
fi

echo "--- ✅ Stack berhasil dihentikan. ---"

echo "--- ⬆️ Menjalankan stack Docker Compose dalam mode detached... ---"
docker compose up -d

if [ $? -ne 0 ]; then
    echo "❌ Gagal menjalankan stack Docker Compose. Mohon periksa error di atas."
    exit 1
fi

echo "--- ✨ Deployment Selesai! Stack Traefik dan Whoami sekarang berjalan. ---"
docker compose ps