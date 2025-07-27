#!/bin/bash

PARENT_ENV="../.env"

if [ ! -f "$PARENT_ENV" ]; then
  echo "Global .env file not found, local .env file will be used, insure that all variables in it"
  exit 0
fi

echo "Get postgres env"

LOCAL_ENV="./.env"

VARS="POSTGRES_USER POSTGRES_PASSWORD POSTGRES_DB POSTGRES_INNER_PORT POSTGRES_HOST REDIS_PASSWORD REDIS_PORT"

for VAR in $VARS; do
  LINE=$(grep -E "^[[:space:]]*${VAR}[[:space:]]*=" "$PARENT_ENV" | grep -v '^[[:space:]]*#' | tail -n 1)

  if [ -n "$LINE" ]; then
    NAME=$(echo "$LINE" | cut -d '=' -f 1 | tr -d ' ')
    RAW_VAL=$(echo "$LINE" | cut -d '=' -f 2- | sed 's/^[[:space:]]*//' | sed 's/[[:space:]]*$//')
    VALUE=$(echo "$RAW_VAL" | sed 's/^"\(.*\)"$/\1/')

    sed -i "/^[[:space:]]*${VAR}[[:space:]]*=/d" "$LOCAL_ENV"

    # Добавить новое
    echo "${NAME}=${VALUE}" >> "$LOCAL_ENV"
  else
    echo "⚠️  Переменная $VAR не найдена в $PARENT_ENV"
  fi
done