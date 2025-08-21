#!/bin/bash

# Hentikan script jika terjadi error
set -e

# Periksa apakah variabel SERVICE_CMD_PATH sudah di-set
if [ -z "$SERVICE_CMD_PATH" ]; then
  echo "Error: SERVICE_CMD_PATH environment variable is not set."
  exit 1
fi

# Lakukan build menggunakan path dari variabel
echo "Building service: $SERVICE_CMD_PATH"
go build -o ./tmp/main $SERVICE_CMD_PATH