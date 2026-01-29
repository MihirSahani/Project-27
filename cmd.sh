#!/bin/zsh

# 1. No spaces around the '=' in assignments
MIGRATION_PATH="migration"

# 2. Ensure POSTGRES_ADDRESS is set before proceeding
if [ -z "$POSTGRES_ADDRESS" ]; then
  echo "Error: POSTGRES_ADDRESS environment variable is not set."
  exit 1
fi

case "$1" in
  "up")
    # 3. Use $VARIABLE or ${VARIABLE} to access vars
    # 4. Use double quotes around variables to handle spaces/special chars
    if [ -z "$2" ]; then
      migrate -path "$MIGRATION_PATH" -database "$POSTGRES_ADDRESS" up
    else
      migrate -path "$MIGRATION_PATH" -database "$POSTGRES_ADDRESS" up "$2"
    fi
    ;;
  "down")
    if [ -z "$2" ]; then
      migrate -path "$MIGRATION_PATH" -database "$POSTGRES_ADDRESS" down
    else
      migrate -path "$MIGRATION_PATH" -database "$POSTGRES_ADDRESS" down "$2"
    fi
    # 5. Fixed the double $$ prefix to a single $
    ;;
  "create")
    if [ -z "$2" ]; then
      echo "Please provide a name for the migration."
      exit 1
    fi
    migrate create -ext sql -dir "$MIGRATION_PATH" -seq "$2"
    ;;
  *)
    echo "Usage: $0 {up|down|create <name>}"
    exit 1
    ;;
esac

exit 0