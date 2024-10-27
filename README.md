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

1. **Clone the repository**

   ```bash
   git clone <repository-url>
   cd <repository-directory>
   ```

2. **Set up the backend**

   ```bash
   cd Backend
   go mod download
   ```

3. **Set up the frontend**

   ```bash
   cd Frontend
   npm install
   ```

### Running the Application

1. **Start the backend**

   ```bash
   cd Backend
   go run main.go
   ```

2. **Start the frontend**

   ```bash
   cd Frontend
   npm start
   ```

3. **Using Docker**

   You can also use Docker to run the application. Ensure Docker is installed and running on your machine.

   ```bash
   docker-compose up --build
   ```

### Environment Variables

Ensure you have a `.env` file in the `Backend` directory with the necessary environment variables. Refer to `.env.example` for guidance.

### Additional Information

- **Frontend Configuration**: The frontend is configured using Webpack and Babel. Tailwind CSS is used for styling, and the configuration can be found in `tailwind.config.js` and `.babelrc`.

- **Backend Configuration**: The backend uses Go modules for dependency management. The configuration is handled using environment variables, and the `config.go` file is responsible for loading these variables.

- **Database Migrations**: Database migrations are managed using `migrate`. Ensure your PostgreSQL database is set up and running before starting the backend.

### Contributing

Contributions are welcome! Please fork the repository and submit a pull request for any improvements or bug fixes.

### License

This project is licensed under the ISC License.
