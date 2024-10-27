import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Box,
  Typography,
  TextField,
  Button,
  CircularProgress,
  Grid,
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
      const response = await fetch('http://localhost:8080/v1/room-types/', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(roomType),
        credentials: 'include'
      });

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
    <Box className="form-container">
      <Typography variant="h4" gutterBottom className="form-title">
        Add New Room Type
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
                label="Name"
                name="name"
                value={roomType.name}
                onChange={handleChange}
                margin="normal"
                required
                className="form-input"
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Capacity"
                name="capacity"
                value={roomType.capacity}
                onChange={handleChange}
                margin="normal"
                type="number"
                inputProps={{ min: 1 }}
                required
                className="form-input"
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Default Price"
                name="default_price"
                value={roomType.default_price}
                onChange={handleChange}
                margin="normal"
                type="number"
                required
                className="form-input"
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Description"
                name="description"
                value={roomType.description}
                onChange={handleChange}
                margin="normal"
                multiline
                rows={3}
                className="form-input"
              />
            </Grid>
            <Grid item xs={12}>
              <Button type="submit" variant="contained" color="primary" className="form-submit" disabled={loading}>
                {loading ? 'Adding...' : 'Add Room Type'}
              </Button>
            </Grid>
          </Grid>
        </form>
      )}
    </Box>
  );
};

export default RoomTypeAdd;