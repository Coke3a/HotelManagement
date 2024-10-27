import React, { useState } from 'react';
import { TextField, Button, Box, Typography, CircularProgress, Select, MenuItem, FormControl, InputLabel, Grid } from '@mui/material';
import { useNavigate } from 'react-router-dom';

const UserAdd = () => {
  const navigate = useNavigate();
  const [user, setUser] = useState({
    username: '',
    password: '',
    role: '',
    rank: '',
    status: 'active',
  });
  const [loading, setLoading] = useState(false);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setUser({ ...user, [name]: value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    try {
      const response = await fetch('http://localhost:8080/v1/users/', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Access-Control-Expose-Headers': 'X-My-Custom-Header, X-Another-Custom-Header'
        },
        body: JSON.stringify(user),
        credentials: 'include'
      });

      if (response.status === 409) {
        throw new Error('User already exists. Please use a different username or email.');
      }

      if (!response.ok) {
        throw new Error('Failed to create user');
      }

      const data = await response.json();
      console.log('User created:', data);
      navigate('/user');
    } catch (error) {
      console.error('Error adding user:', error);
      // Display error message to the user
      alert(error.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Box className="form-container">
      <Typography variant="h4" gutterBottom className="form-title">
        Add User
      </Typography>
      <form onSubmit={handleSubmit} className="form">
        <Grid container spacing={2}>
          <Grid item xs={12} sm={6}>
            <TextField
              fullWidth
              label="Username"
              name="username"
              value={user.username}
              onChange={handleChange}
              margin="normal"
              required
              className="form-input"
              size="small"
            />
          </Grid>
          <Grid item xs={12} sm={6}>
            <TextField
              fullWidth
              label="Password"
              name="password"
              value={user.password}
              onChange={handleChange}
              margin="normal"
              required
              type="password"
              className="form-input"
              size="small"
            />
          </Grid>
          <Grid item xs={12} sm={6}>
            <FormControl fullWidth margin="normal" required className="form-input">
              <InputLabel>Role</InputLabel>
              <Select
                name="role"
                value={user.role}
                onChange={handleChange}
                label="Role"
                size="small"
              >
                <MenuItem value={0}>Staff</MenuItem>
                <MenuItem value={1}>Admin</MenuItem>
              </Select>
            </FormControl>
          </Grid>
          <Grid item xs={12} sm={6}>
            <FormControl fullWidth margin="normal" required className="form-input">
              <InputLabel>Rank</InputLabel>
              <Select
                name="rank"
                value={user.rank}
                onChange={handleChange}
                label="Rank"
                size="small"
              >
                <MenuItem value="1">1</MenuItem>
                <MenuItem value="2">2</MenuItem>
                <MenuItem value="3">3</MenuItem>
              </Select>
            </FormControl>
          </Grid>
        </Grid>
        <Box display="flex" justifyContent="flex-end" mt={2}>
          <Button type="submit" variant="contained" color="primary" className="form-button" size="small">
            Add User
          </Button>
        </Box>
      </form>
    </Box>
  );
};

export default UserAdd;
