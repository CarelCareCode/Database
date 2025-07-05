#!/bin/bash

# Emergency Response System - AWS Deployment Script
# This script deploys the complete AWS infrastructure

set -e

# Configuration
ENVIRONMENT=${1:-development}
AWS_REGION=${AWS_REGION:-us-east-1}
STACK_NAME="emergency-response-infrastructure"
ECR_REPOSITORY="emergency-response-backend"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Logging function
log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

warn() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] WARNING: $1${NC}"
}

error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] ERROR: $1${NC}"
    exit 1
}

# Check prerequisites
check_prerequisites() {
    log "Checking prerequisites..."
    
    # Check if AWS CLI is installed
    if ! command -v aws &> /dev/null; then
        error "AWS CLI is not installed. Please install it first."
    fi
    
    # Check if Docker is installed
    if ! command -v docker &> /dev/null; then
        error "Docker is not installed. Please install it first."
    fi
    
    # Check if jq is installed
    if ! command -v jq &> /dev/null; then
        error "jq is not installed. Please install it first."
    fi
    
    # Check AWS credentials
    if ! aws sts get-caller-identity &> /dev/null; then
        error "AWS credentials not configured. Please run 'aws configure'"
    fi
    
    log "Prerequisites check completed successfully"
}

# Create ECR repository if it doesn't exist
create_ecr_repository() {
    log "Creating ECR repository..."
    
    aws ecr describe-repositories --repository-names $ECR_REPOSITORY --region $AWS_REGION &> /dev/null || {
        log "Creating ECR repository: $ECR_REPOSITORY"
        aws ecr create-repository --repository-name $ECR_REPOSITORY --region $AWS_REGION
    }
    
    log "ECR repository ready"
}

# Build and push Docker image
build_and_push_image() {
    log "Building and pushing Docker image..."
    
    # Get ECR login token
    aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin $(aws sts get-caller-identity --query Account --output text).dkr.ecr.$AWS_REGION.amazonaws.com
    
    # Build image
    docker build -t $ECR_REPOSITORY .
    
    # Tag image
    ECR_URI=$(aws sts get-caller-identity --query Account --output text).dkr.ecr.$AWS_REGION.amazonaws.com/$ECR_REPOSITORY:latest
    docker tag $ECR_REPOSITORY:latest $ECR_URI
    
    # Push image
    docker push $ECR_URI
    
    log "Docker image pushed successfully to $ECR_URI"
}

# Deploy CloudFormation stack
deploy_infrastructure() {
    log "Deploying CloudFormation stack..."
    
    # Check if secrets exist
    if [ -z "$DATABASE_PASSWORD" ]; then
        error "DATABASE_PASSWORD environment variable is required"
    fi
    
    if [ -z "$JWT_SECRET" ]; then
        error "JWT_SECRET environment variable is required"
    fi
    
    # Deploy stack
    aws cloudformation deploy \
        --template-file infrastructure/cloudformation.yml \
        --stack-name $STACK_NAME \
        --capabilities CAPABILITY_IAM \
        --parameter-overrides \
            Environment=$ENVIRONMENT \
            DatabasePassword=$DATABASE_PASSWORD \
            JWTSecret=$JWT_SECRET \
        --region $AWS_REGION
    
    log "CloudFormation stack deployed successfully"
}

# Get stack outputs
get_stack_outputs() {
    log "Getting stack outputs..."
    
    # Get stack outputs
    OUTPUTS=$(aws cloudformation describe-stacks --stack-name $STACK_NAME --region $AWS_REGION --query 'Stacks[0].Outputs' --output json)
    
    # Extract important endpoints
    DB_ENDPOINT=$(echo $OUTPUTS | jq -r '.[] | select(.OutputKey=="DatabaseEndpoint") | .OutputValue')
    REDIS_ENDPOINT=$(echo $OUTPUTS | jq -r '.[] | select(.OutputKey=="RedisEndpoint") | .OutputValue')
    ECS_CLUSTER=$(echo $OUTPUTS | jq -r '.[] | select(.OutputKey=="ECSClusterName") | .OutputValue')
    
    log "Database Endpoint: $DB_ENDPOINT"
    log "Redis Endpoint: $REDIS_ENDPOINT"
    log "ECS Cluster: $ECS_CLUSTER"
    
    # Save outputs to file
    echo "DB_ENDPOINT=$DB_ENDPOINT" > deployment-outputs.env
    echo "REDIS_ENDPOINT=$REDIS_ENDPOINT" >> deployment-outputs.env
    echo "ECS_CLUSTER=$ECS_CLUSTER" >> deployment-outputs.env
    
    log "Outputs saved to deployment-outputs.env"
}

# Run database migrations
run_migrations() {
    log "Running database migrations..."
    
    # This would typically be done via a migration job
    # For now, we'll create a simple task that runs migrations
    warn "Database migrations should be run manually or via a separate migration job"
    log "Migration files are available in the migrations/ directory"
}

# Main deployment function
main() {
    log "Starting deployment for environment: $ENVIRONMENT"
    
    check_prerequisites
    create_ecr_repository
    build_and_push_image
    deploy_infrastructure
    get_stack_outputs
    run_migrations
    
    log "Deployment completed successfully!"
    log "Your Emergency Response System is now running in AWS"
    log "Check the AWS Console for service status and endpoints"
}

# Help function
show_help() {
    cat << EOF
Emergency Response System - AWS Deployment Script

Usage: $0 [ENVIRONMENT]

ENVIRONMENT: development, staging, or production (default: development)

Required environment variables:
- DATABASE_PASSWORD: Password for Aurora PostgreSQL database
- JWT_SECRET: Secret key for JWT token generation (minimum 32 characters)

Optional environment variables:
- AWS_REGION: AWS region (default: us-east-1)

Examples:
  $0 development
  $0 production
  DATABASE_PASSWORD=mypassword JWT_SECRET=mysecretkey $0 production

EOF
}

# Check if help is requested
if [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
    show_help
    exit 0
fi

# Run main function
main 