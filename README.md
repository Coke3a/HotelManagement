# Hotel Management System

This project is a comprehensive hotel management system with a React frontend and a Go backend.

## Project Structure

The project is divided into two main parts:

1. **Frontend (React)**
2. **Backend (Go)**

### Frontend

The frontend is built using React and includes the following key features:

- Material-UI components for a modern UI
- React Router for navigation
- Recharts for data visualization
- Tailwind CSS for additional styling

### Backend

The backend is built using Go and includes:

- Gin web framework
- PostgreSQL database
- PASETO for authentication
- Swagger for API documentation

## Prerequisites

- Docker and Docker Compose
- Git

## Running with Docker

### Development Environment

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd <repository-directory>
   ```

2. Create a `.env` file in the Backend directory:
   ```env
   # Application
   APP_NAME=hotel_management
   APP_ENV=development

   # Database
   DB_CONNECTION=postgres
   DB_HOST=postgres
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_NAME=hotel_db

   # HTTP
   HTTP_PORT=8080
   HTTP_ALLOWED_ORIGINS=http://localhost:9500

   # Token
   TOKEN_DURATION=24h
   ```

3. Start the development environment:
   ```bash
   docker-compose up --build
   ```

   The application will be available at:
   - Frontend: http://localhost:9500
   - Backend: http://localhost:8080
   - Database: localhost:5432

4. Development features:
   - Hot reloading for frontend changes
   - Live reload for backend changes
   - Direct database access
   - Automatic migration execution

### Production Environment

1. Create a `.env.prod` file in the Backend directory with your production settings.

2. Deploy using production configuration:
   ```bash
   docker-compose -f docker-compose.prod.yml up --build -d
   ```

3. Production features:
   - Optimized builds
   - Secure environment configurations
   - Automatic container restart
   - Production-grade PostgreSQL settings

### Common Docker Commands

```bash
# View running containers
docker-compose ps

# View logs
docker-compose logs -f

# Stop the application
docker-compose down

# Remove volumes (database data)
docker-compose down -v

# Rebuild specific service
docker-compose build <service-name>
```

## Additional Information

- **Frontend Configuration**: The frontend is configured using Webpack and Babel. Tailwind CSS is used for styling.
- **Backend Configuration**: The backend uses Go modules for dependency management.
- **Database Migrations**: Migrations run automatically when the backend container starts.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request for any improvements or bug fixes.

## License

This project is licensed under the ISC License.
