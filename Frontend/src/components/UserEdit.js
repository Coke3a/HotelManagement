import React, { useState, useEffect } from 'react';
import { TextField, Button, Box, Typography, CircularProgress, Select, MenuItem, FormControl, InputLabel, Grid } from '@mui/material';
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
        const response = await fetch(`http://localhost:8080/v1/users/${id}`, {
          headers: {
            'Content-Type': 'application/json',
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
    try {
      const response = await fetch(`http://localhost:8080/v1/users/`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Access-Control-Expose-Headers': 'X-My-Custom-Header, X-Another-Custom-Header'
        },
        body: JSON.stringify(user),
        credentials: 'include'
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
      try {
        const response = await fetch(`http://localhost:8080/v1/users/${id}`, {
          method: 'DELETE',
          credentials: 'include'
        });

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
    <Box className="form-container">
      <Typography variant="h4" gutterBottom className="form-title">
        Edit User
      </Typography>
      {loading ? (
        <Box display="flex" justifyContent="center" alignItems="center" height="300px">
          <CircularProgress />
        </Box>
      ) : (
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
          <Box display="flex" justifyContent="center" mt={2}>
            <Button
              variant="contained"
              color="error"
              onClick={handleDelete}
              disabled={loading}
              size="medium"
              style={{ marginRight: '8px' }}
            >
              Delete User
            </Button>
            <Button type="submit" variant="contained" color="primary" disabled={loading} size="medium">
              {loading ? 'Updating...' : 'Update User'}
            </Button>
          </Box>
        </form>
      )}
    </Box>
  );
};

export default UserEdit;
