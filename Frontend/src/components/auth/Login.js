import React, { useState, useEffect } from 'react';
import { TextField, Button, Box, Typography, CircularProgress, Paper, Alert } from '@mui/material';
import { useNavigate, useLocation, Navigate } from 'react-router-dom';
import { setToken, setUserRole, setUserId, setUsername, isAuthenticated } from '../../utils/auth';

const Login = () => {
  // If already authenticated, redirect to dashboard
  if (isAuthenticated()) {
    return <Navigate to="/" replace />;
  }

  const navigate = useNavigate();
  const location = useLocation(); // Add this line to get location state
  const [credentials, setCredentials] = useState({
    username: '',
    password: '',
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(''); // Add state for error message

  // Add this useEffect to handle error message from location state
  useEffect(() => {
    if (location.state?.error) {
      setError(location.state.error);
      // Clear the location state after showing the message
      navigate(location.pathname, { replace: true, state: {} });
    }
  }, [location, navigate]);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setCredentials({ ...credentials, [name]: value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(''); // Clear previous error message
    try {
      const response = await fetch('http://localhost:8080/v1/users/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(credentials),
      });

      if (!response.ok) {
        throw new Error('Login failed');
      }
      
      const data = await response.json();

      setToken(data.data.token);
      setUserRole(data.data.role);
      setUsername(data.data.username);
      navigate('/dashboard');
    } catch (error) {
      console.error('Login error:', error);
      setError('Login failed. Please check your username and password.'); // Set error message
    } finally {
      setLoading(false);
    }
  };

  return (
    <Box display="flex" justifyContent="center" alignItems="center" height="100vh" bgcolor="#f0f2f5">
      <Paper elevation={3} style={{ padding: '2rem', maxWidth: '400px', width: '100%' }}>
        <Typography variant="h3" gutterBottom align="center" color="primary">
          Hotel Management
        </Typography>
        <Typography variant="h5" gutterBottom align="center">
          Login
        </Typography>
        {error && <Alert severity="error">{error}</Alert>} {/* Display error message */}
        {loading ? (
          <Box display="flex" justifyContent="center" alignItems="center" height="200px">
            <CircularProgress />
          </Box>
        ) : (
          <form onSubmit={handleSubmit} className="form">
            <TextField
              fullWidth
              label="Username"
              name="username"
              value={credentials.username}
              onChange={handleChange}
              margin="normal"
              required
              className="form-input"
            />
            <TextField
              fullWidth
              label="Password"
              name="password"
              value={credentials.password}
              onChange={handleChange}
              margin="normal"
              type="password"
              required
              className="form-input"
            />
            <Button type="submit" variant="contained" color="primary" className="form-submit" disabled={loading}>
              {loading ? 'Logging in...' : 'Login'}
            </Button>
          </form>
        )}
      </Paper>
    </Box>
  );
};

export default Login;
