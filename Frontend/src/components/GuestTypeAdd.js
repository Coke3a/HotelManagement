import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Box,
  Typography,
  TextField,
  Button,
  CircularProgress,
  Paper,
  Grid,
} from '@mui/material';

const GuestTypeAdd = () => {
  const token = localStorage.getItem('token');
  const navigate = useNavigate();
  const [guestType, setGuestType] = useState({
    name: '',
    description: '',
  });
  const [loading, setLoading] = useState(false);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setGuestType({ ...guestType, [name]: value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    try {
      const response = await fetch('http://localhost:8080/v1/customer-types/', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify(guestType),
      });

      if (response.status === 401) {
        handleTokenExpiration(new Error("access token has expired"), navigate);
        return;
      }

      if (response.status === 409) {
        throw new Error('Guest type already exists. Please use a different name.');
      }

      if (!response.ok) {
        throw new Error('Failed to create guest type');
      }

      const data = await response.json();
      console.log('Guest type created:', data);
      navigate('/guest-type');
    } catch (error) {
      console.error('Error adding guest type:', error);
      alert(error.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Box sx={{ maxWidth: '800px', margin: '0 auto', padding: '20px' }}>
      <Paper elevation={3} sx={{ padding: '24px', backgroundColor: '#f8f9fa' }}>
        <Typography variant="h5" gutterBottom className="form-title">
          Add New Guest Type
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
                  label="Name"
                  name="name"
                  value={guestType.name}
                  onChange={handleChange}
                  required
                  className="form-input"
                  size="small"
                />
              </Grid>
              <Grid item xs={12}>
                <TextField
                  fullWidth
                  label="Description"
                  name="description"
                  value={guestType.description}
                  onChange={handleChange}
                  multiline
                  rows={3}
                  className="form-input"
                  size="small"
                />
              </Grid>
            </Grid>
            <Box display="flex" justifyContent="flex-end" mt={3} gap={2}>
              <Button
                type="submit"
                variant="contained"
                color="primary"
                disabled={loading}
                size="medium"
              >
                {loading ? 'Adding...' : 'Add Guest Type'}
              </Button>
            </Box>
          </form>
        )}
      </Paper>
    </Box>
  );
};

export default GuestTypeAdd;
