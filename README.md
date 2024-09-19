# Hotel Management System

This project is a comprehensive hotel management system with a React frontend and a Go backend.

## Project Structure

The project is divided into two main parts:

1. Frontend (React)
2. Backend (Go)

### Frontend

The frontend is built using React and includes the following key features:

- Material-UI components for a modern UI
- React Router for navigation
- Recharts for data visualization
- Tailwind CSS for additional styling

#### Key Components

- Dashboard
- Guest management
- Room management
- Booking management
- User management
- Rate and pricing management

### Backend

The backend is built using Go and includes:

- Gin web framework
- PostgreSQL database
- PASETO for authentication
- Swagger for API documentation

## Getting Started

### Prerequisites

- Node.js and npm for the frontend
- Go 1.22.3 or later for the backend
- PostgreSQL database

### Installation

1. Clone the repository
2. Set up the backend:
   ```
   cd Backend
   go mod download
   ```
3. Set up the frontend:
   ```
   cd Frontend
   npm install
   ```

### Running the Application

1. Start the backend server:
   ```
   cd Backend
   go run main.go
   ```
2. Start the frontend development server:
   ```
   cd Frontend
   npm start
   ```

The frontend will be available at `http://localhost:9500`.

## Features

- User authentication and authorization
- Guest management
- Room management
- Booking system
- Rate and pricing management
- Dashboard with analytics

## Contributing

Please read CONTRIBUTING.md for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the [MIT License](LICENSE).
