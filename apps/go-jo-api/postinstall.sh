#!/bin/bash
# go-jo-api postinstall script

# Create directory for configuration
mkdir -p /etc/go-jo-api

# Enable the service (don't start it automatically)
systemctl enable go-jo-api || true

echo "go-jo-api service enabled. Configure /etc/go-jo-api/.env and start with: systemctl start go-jo-api" 