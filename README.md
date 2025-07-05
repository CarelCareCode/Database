# Emergency Response System - Backend

A high-performance Go backend for emergency response services with real-time capabilities, built for AWS cloud infrastructure.

## 🏗️ Architecture Overview

This system implements a complete emergency response backend with:

- **Database**: AWS Aurora PostgreSQL with PostGIS for spatial queries
- **Real-time Cache**: AWS ElastiCache Redis for live location tracking and pub/sub
- **Event System**: AWS Kinesis for event-driven architecture
- **Backend**: Go with Gorilla WebSocket for real-time communication
- **Deployment**: AWS Fargate containers with automated CI/CD
- **Spatial Indexing**: H3 hexagonal grid system for efficient location queries

## 🚀 Features

### Core Functionality
- ✅ User registration and authentication (JWT)
- ✅ Medical profile management with encryption
- ✅ Emergency incident creation and tracking
- ✅ Real-time paramedic assignment using H3 spatial indexing
- ✅ WebSocket-based real-time communication
- ✅ Chat system (dispatcher ↔ paramedic, dispatcher ↔ client)
- ✅ Live location tracking with Redis
- ✅ Event-driven workflow with AWS Kinesis

### Technical Features
- ✅ PostgreSQL with PostGIS for spatial data
- ✅ H3 hexagonal indexing for proximity matching
- ✅ Redis pub/sub for real-time updates
- ✅ AWS Kinesis event streaming
- ✅ Docker containerization
- ✅ GitHub Actions CI/CD pipeline
- ✅ CloudFormation infrastructure as code

## 🛠️ Technology Stack

| Component | Technology |
|-----------|------------|
| **Backend** | Go 1.21, Gorilla Mux, Gorilla WebSocket |
| **Database** | AWS Aurora PostgreSQL + PostGIS |
| **Cache** | AWS ElastiCache Redis |
| **Event Bus** | AWS Kinesis |
| **Compute** | AWS Fargate |
| **Spatial** | H3 Hexagonal Grid System |
| **Auth** | JWT with bcrypt |
| **Deploy** | Docker, GitHub Actions, CloudFormation |

## 📦 Project Structure

```
emergency-response-backend/
├── .github/workflows/          # GitHub Actions CI/CD
├── infrastructure/             # AWS CloudFormation templates
├── internal/                   # Go application code
│   ├── config/                # Configuration management
│   ├── database/              # Database connection and utilities
│   ├── handlers/              # HTTP and WebSocket handlers
│   ├── middleware/            # Authentication and logging middleware
│   ├── models/                # Data models and structs
│   ├── redis/                 # Redis client and utilities
│   ├── server/                # HTTP server setup
│   └── websocket/             # WebSocket hub management
├── migrations/                # Database migration scripts
├── scripts/                   # Deployment and utility scripts
├── Context/                   # Project documentation
├── docker-compose.yml         # Local development environment
├── Dockerfile                 # Production container
├── go.mod                     # Go module dependencies
└── main.go                    # Application entry point
```

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Docker and Docker Compose
- AWS CLI configured
- PostgreSQL with PostGIS (for local development)
- Redis (for local development)

### Local Development

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd emergency-response-backend
   ```

2. **Set up environment variables**
   ```bash
   cp env.example .env
   # Edit .env with your configuration
   ```

3. **Start local services**
   ```bash
   docker-compose up -d postgres redis
   ```

4. **Run database migrations**
   ```bash
   go run scripts/migrate.go
   ```

5. **Start the application**
   ```bash
   go run main.go
   ```

6. **Access the services**
   - API: http://localhost:8080
   - Health Check: http://localhost:8080/health
   - WebSocket: ws://localhost:8080/api/v1/ws
   - Adminer (DB UI): http://localhost:8081

### Using Docker Compose

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f backend

# Stop services
docker-compose down
```

## 🔧 API Endpoints

### Authentication
- `POST /api/v1/register` - User registration
- `POST /api/v1/login` - User login

### User Management
- `GET /api/v1/profile` - Get user profile
- `POST /api/v1/medical` - Create/update medical profile
- `GET /api/v1/medical/{user_id}` - Get medical profile

### Emergency Services
- `POST /api/v1/emergency` - Create emergency incident
- `GET /api/v1/incidents` - Get incidents (dispatcher)
- `GET /api/v1/incidents/{id}` - Get specific incident
- `POST /api/v1/incidents/{id}/assign` - Assign paramedic
- `PUT /api/v1/incidents/{id}/status` - Update incident status

### Real-time Communication
- `GET /api/v1/ws` - WebSocket connection
- `POST /api/v1/chat` - Send chat message
- `GET /api/v1/chat/{incident_id}` - Get chat history

### Paramedic Services
- `GET /api/v1/paramedics` - Get available paramedics
- `POST /api/v1/paramedics/location` - Update paramedic location

### System
- `GET /health` - Health check endpoint

## 🌍 AWS Deployment

### Prerequisites

- AWS CLI configured with appropriate permissions
- Docker installed locally
- Environment variables set:
  - `DATABASE_PASSWORD`: Aurora PostgreSQL password
  - `JWT_SECRET`: JWT secret key (32+ characters)

### Automated Deployment

```bash
# Deploy to development
DATABASE_PASSWORD=mypassword JWT_SECRET=mysecretkey ./scripts/deploy.sh development

# Deploy to production
DATABASE_PASSWORD=mypassword JWT_SECRET=mysecretkey ./scripts/deploy.sh production
```

### Manual Deployment

1. **Create ECR repository**
   ```bash
   aws ecr create-repository --repository-name emergency-response-backend --region us-east-1
   ```

2. **Build and push Docker image**
   ```bash
   docker build -t emergency-response-backend .
   docker tag emergency-response-backend:latest <account-id>.dkr.ecr.us-east-1.amazonaws.com/emergency-response-backend:latest
   docker push <account-id>.dkr.ecr.us-east-1.amazonaws.com/emergency-response-backend:latest
   ```

3. **Deploy infrastructure**
   ```bash
   aws cloudformation deploy \
     --template-file infrastructure/cloudformation.yml \
     --stack-name emergency-response-infrastructure \
     --capabilities CAPABILITY_IAM \
     --parameter-overrides \
       Environment=production \
       DatabasePassword=<your-password> \
       JWTSecret=<your-jwt-secret>
   ```

### GitHub Actions CI/CD

The repository includes automated CI/CD pipeline that:

1. **Tests**: Runs Go tests with PostgreSQL and Redis
2. **Builds**: Creates Docker image and pushes to ECR
3. **Deploys**: Updates CloudFormation stack and ECS service

Required GitHub Secrets:
- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`
- `DATABASE_PASSWORD`
- `JWT_SECRET`

## 🗄️ Database Schema

### Core Tables

- `users` - User accounts and authentication
- `medical_profiles` - Medical information (encrypted)
- `clients` - Organization/client information
- `zones` - Geographic service areas
- `paramedics` - Paramedic information and status
- `incidents` - Emergency incidents with H3 indexing
- `chat_messages` - Real-time chat messages

### Spatial Features

- PostGIS extensions for geographic queries
- H3 hexagonal indexing for efficient proximity searches
- Location geometry storage for incidents and boundaries

## 🔒 Security

- JWT authentication with bcrypt password hashing
- Medical data encryption at rest using AWS KMS
- VPC isolation for database and cache
- Security groups with minimal access rules
- Non-root Docker container execution
- AWS IAM roles with least privilege

## 🎯 Performance

- H3 spatial indexing for O(1) proximity queries
- Redis caching for frequently accessed data
- Connection pooling for database efficiency
- Horizontal scaling with AWS Fargate
- CloudFront CDN for static assets (when added)

## 🔧 Configuration

### Environment Variables

See `env.example` for all available configuration options.

### Key Configuration Areas

- **Database**: Aurora PostgreSQL connection settings
- **Redis**: ElastiCache connection and pub/sub settings
- **AWS**: Kinesis stream and SNS topic configuration
- **Security**: JWT secrets and encryption keys
- **Spatial**: H3 resolution and PostGIS settings

## 📊 Monitoring

- CloudWatch logs for application monitoring
- Health check endpoints for load balancer
- Database and Redis metrics via CloudWatch
- ECS task monitoring and auto-scaling
- Kinesis stream monitoring

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -v -race -coverprofile=coverage.out ./...

# Run specific test
go test -v ./internal/handlers
```

## 📋 TODO

- [ ] Add comprehensive unit tests
- [ ] Implement AWS Cognito integration
- [ ] Add API rate limiting
- [ ] Implement incident analytics
- [ ] Add Prometheus metrics
- [ ] Create performance benchmarks
- [ ] Add OpenAPI/Swagger documentation

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run tests and ensure they pass
6. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

For issues and support:
1. Check the GitHub Issues page
2. Review the deployment logs in CloudWatch
3. Verify AWS service status
4. Check the health endpoint: `/health`

## 🔗 Related Resources

- [AWS Aurora PostgreSQL Documentation](https://docs.aws.amazon.com/AmazonRDS/latest/AuroraUserGuide/Aurora.AuroraPostgreSQL.html)
- [PostGIS Documentation](https://postgis.net/documentation/)
- [H3 Spatial Indexing](https://h3geo.org/)
- [Gorilla WebSocket Documentation](https://pkg.go.dev/github.com/gorilla/websocket)
- [AWS Kinesis Documentation](https://docs.aws.amazon.com/kinesis/) 