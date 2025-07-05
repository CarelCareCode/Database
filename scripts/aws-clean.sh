#!/bin/bash
# Clean AWS CLI wrapper to avoid shell interference
# This script provides a clean environment for AWS CLI commands

# Set clean environment variables
export PATH="/usr/bin:/bin:/usr/sbin:/sbin:/opt/homebrew/bin"
export AWS_DEFAULT_REGION="af-south-1"

# Unset any problematic functions
unset -f head cat 2>/dev/null || true

# Execute AWS CLI with clean environment
exec /opt/homebrew/bin/aws "$@" 