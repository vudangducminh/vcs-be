#!/bin/bash

# Redis Cloud Health Check Script
REDIS_HOST="redis-13990.c251.east-us-mz.azure.redns.redis-cloud.com"
REDIS_PORT="13990"
REDIS_USER="default"
REDIS_PASSWORD="c3QOPfZSpiPqTmfBCINbFuPvaKUSEMM8"

# Test Redis connection
redis-cli -h "$REDIS_HOST" -p "$REDIS_PORT" -a "$REDIS_PASSWORD" --user "$REDIS_USER" ping

if [ $? -eq 0 ]; then
    echo "Redis Cloud is healthy and responding"
    exit 0
else
    echo "Redis Cloud health check failed"
    exit 1
fi
