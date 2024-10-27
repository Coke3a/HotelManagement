import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Box,
  Typography,
  TextField,
  Button,
  CircularProgress,
  Grid,
  Paper,
} from '@mui/material';

const RoomTypeAdd = () => {
  const navigate = useNavigate();
  const [roomType, setRoomType] = useState({
    name: '',
    description: '',
    capacity: '',
    default_price: '',
  });
  const [loading, setLoading] = useState(false);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setRoomType({
      ...roomType,
      [name]: name === 'capacity' || name === 'default_price' ? Number(value) : value
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    try {
      const token = localStorage.getItem('token');
      const response = await fetch('http://localhost:8080/v1/room-types/', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify(roomType),
      });

      if (response.status === 401) {
        handleTokenExpiration(new Error("access token has expired"), navigate);
        return;
      }

      if (response.status === 409) {
        throw new Error('Room type already exists. Please use a different name.');
      }

      if (!response.ok) {
        throw new Error('Failed to create room type');
      }

      const data = await response.json();
      console.log('Room type created:', data);
      navigate('/room-type');
    } catch (error) {
      console.error('Error adding room type:', error);
      alert(error.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Box sx={{ maxWidth: '800px', margin: '0 auto', padding: '20px' }}>
      <Paper elevation={3} sx={{ padding: '24px', backgroundColor: '#f8f9fa' }}>
        <Typography variant="h5" gutterBottom className="form-title">
          Add New Room Type
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
                  value={roomType.name}
                  onChange={handleChange}
                  required
                  className="form-input"
                  size="small"
                />
              </Grid>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label="Capacity"
                  name="capacity"
                  value={roomType.capacity}
                  onChange={handleChange}
                  type="number"
                  required
                  className="form-input"
                  size="small"
                  inputProps={{ min: 1 }}
                />
              </Grid>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label="Default Price"
                  name="default_price"
                  value={roomType.default_price}
                  onChange={handleChange}
                  type="number"
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
                  value={roomType.description}
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
                {loading ? 'Adding...' : 'Add Room Type'}
              </Button>
            </Box>
          </form>
        )}
      </Paper>
    </Box>
  );
};

export default RoomTypeAdd;
