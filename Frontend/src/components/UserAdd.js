import React, { useState } from 'react';
import { TextField, Button, Box, Typography, CircularProgress, Select, MenuItem, FormControl, InputLabel, Grid, Paper } from '@mui/material';
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
      const token = localStorage.getItem('token');
      const response = await fetch('http://localhost:8080/v1/users/', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify(user),
      });

      if (response.status === 409) {
        throw new Error('User already exists. Please use a different username or email.');
      }

      if (response.status === 401) {
        handleTokenExpiration(new Error("access token has expired"), navigate);
        return;
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
    <Box sx={{ maxWidth: '800px', margin: '0 auto', padding: '20px' }}>
      <Paper elevation={3} sx={{ padding: '24px', backgroundColor: '#f8f9fa' }}>
        <Typography variant="h5" gutterBottom className="form-title">
          Add New User
        </Typography>
        <form onSubmit={handleSubmit} className="form">
          <Grid container spacing={3}>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Username"
                name="username"
                value={user.username}
                onChange={handleChange}
                required
                className="form-input"
                size="small"
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <TextField
                fullWidth
                label="Password"
                name="password"
                type="password"
                value={user.password}
                onChange={handleChange}
                required
                className="form-input"
                size="small"
              />
            </Grid>
            <Grid item xs={12} md={6}>
              <FormControl fullWidth required size="small">
                <InputLabel>Role</InputLabel>
                <Select
                  name="role"
                  value={user.role}
                  onChange={handleChange}
                  label="Role"
                >
                  <MenuItem value={0}>Staff</MenuItem>
                  <MenuItem value={1}>Admin</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} md={6}>
              <FormControl fullWidth required size="small">
                <InputLabel>Rank</InputLabel>
                <Select
                  name="rank"
                  value={user.rank}
                  onChange={handleChange}
                  label="Rank"
                >
                  <MenuItem value="1">1</MenuItem>
                  <MenuItem value="2">2</MenuItem>
                  <MenuItem value="3">3</MenuItem>
                </Select>
              </FormControl>
            </Grid>
          </Grid>
          <Box display="flex" justifyContent="flex-end" mt={3}>
            <Button
              type="button"
              variant="outlined"
              onClick={() => navigate('/user')}
              sx={{ mr: 2 }}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              variant="contained"
              disabled={loading}
            >
              {loading ? <CircularProgress size={24} /> : 'Add User'}
            </Button>
          </Box>
        </form>
      </Paper>
    </Box>
  );
};

export default UserAdd;
