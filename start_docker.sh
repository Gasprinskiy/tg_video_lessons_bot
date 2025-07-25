#!/bin/bash

if [[ "$1" == "-prod" ]]; then
  DETACH_FLAG="-d"
else
  DETACH_FLAG=""
fi


./get_pg_env.sh

echo "Run bot_api docker"

# Запуск
docker compose down
docker compose build
docker compose -p bot_api up --force-recreate --remove-orphans $DETACH_FLAG
