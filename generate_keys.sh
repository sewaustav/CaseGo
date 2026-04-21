#!/bin/bash
# Генерация RSA-ключей для JWT (RS256)
# Запустить один раз: bash generate_keys.sh

set -e

mkdir -p Auth/keys

# Генерируем приватный ключ
openssl genrsa -out Auth/keys/private.pem 2048

# Публичный ключ из приватного
openssl rsa -in Auth/keys/private.pem -pubout -out Auth/keys/public.pem

echo "Ключи сгенерированы: Auth/keys/private.pem и Auth/keys/public.pem"
echo ""
echo "Скопируй содержимое для .env файлов Go-сервисов:"
echo ""
echo "=== PUBLIC_KEY ==="
awk 'NF {sub(/\r/, ""); printf "%s\\n",$0;} END{print ""}' Auth/keys/public.pem
echo ""
echo "=== PRIVATE_KEY ==="
awk 'NF {sub(/\r/, ""); printf "%s\\n",$0;} END{print ""}' Auth/keys/private.pem
