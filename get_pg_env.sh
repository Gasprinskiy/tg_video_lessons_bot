#!/bin/bash

echo "Get postgres env"

PARENT_ENV="../.env"
LOCAL_ENV="./.env"

VARS="POSTGRES_USER POSTGRES_PASSWORD POSTGRES_DB POSTGRES_INNER_PORT POSTGRES_HOST"

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