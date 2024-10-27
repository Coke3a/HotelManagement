import React, { useState, useEffect } from 'react';
import { TextField, Button, Box, Typography, CircularProgress, Select, MenuItem, FormControl, InputLabel, Grid, Paper } from '@mui/material';
import { useNavigate, useParams } from 'react-router-dom';

const UserEdit = () => {
  const navigate = useNavigate();
  const { id } = useParams();
  const [user, setUser] = useState({
    id: '',
    username: '',
    role: '',
    rank: '',
    status: '',
  });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchUser = async () => {
      setLoading(true);
      try {
        const token = localStorage.getItem('token');
        const response = await fetch(`http://localhost:8080/v1/users/${id}`, {
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`,
          },
          credentials: 'include'
        });

        if (!response.ok) {
          throw new Error('Failed to fetch user');
        }

        const data = await response.json();
        setUser(data.data);
      } catch (error) {
        console.error('Error fetching user:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchUser();
  }, [id]);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setUser({ ...user, [name]: value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);

    const token = localStorage.getItem('token');

    try {
      const response = await fetch(`http://localhost:8080/v1/users/`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify(user),
      });

      if (!response.ok) {
        throw new Error('Failed to update user');
      }

      const data = await response.json();
      console.log('User updated:', data);
      navigate('/user');
    } catch (error) {
      console.error('Error updating user:', error);
      alert(error.message);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async () => {
    if (window.confirm('Are you sure you want to delete this user?')) {
      setLoading(true);

      const token = localStorage.getItem('token');
      try {
        const response = await fetch(`http://localhost:8080/v1/users/${id}`, {
          method: 'DELETE',
          headers: {
            'Authorization': `Bearer ${token}`,
          }
        });

        if (response.status === 401) {
          handleTokenExpiration(new Error("access token has expired"), navigate);
          return;
        }

        if (!response.ok) {
          throw new Error('Failed to delete user');
        }

        console.log('User deleted');
        navigate('/user');
      } catch (error) {
        console.error('Error deleting user:', error);
        alert(error.message);
      } finally {
        setLoading(false);
      }
    }
  };

  return (
    <Box sx={{ maxWidth: '800px', margin: '0 auto', padding: '20px' }}>
      <Paper elevation={3} sx={{ padding: '24px', backgroundColor: '#f8f9fa' }}>
        <Typography variant="h5" gutterBottom className="form-title">
          Edit User
        </Typography>
        {loading ? (
          <Box display="flex" justifyContent="center" alignItems="center" height="300px">
            <CircularProgress />
          </Box>
        ) : (
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
            <Box display="flex" justifyContent="flex-end" mt={3} gap={2}>
              <Button
                variant="outlined"
                color="error"
                onClick={handleDelete}
                disabled={loading}
                size="medium"
              >
                Delete User
              </Button>
              <Button
                type="submit"
                variant="contained"
                color="primary"
                disabled={loading}
                size="medium"
              >
                {loading ? 'Updating...' : 'Update User'}
              </Button>
            </Box>
          </form>
        )}
      </Paper>
    </Box>
  );
};

export default UserEdit;
