import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import {
  Box,
  Typography,
  TextField,
  Button,
  CircularProgress,
  Grid,
} from '@mui/material';

const RoomTypeEdit = () => {
  const navigate = useNavigate();
  const { id } = useParams();
  const [roomType, setRoomType] = useState({
    id: '',
    name: '',
    description: '',
    capacity: '',
    default_price: '',
    created_at: '',
    updated_at: '',
  });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchRoomType = async () => {
      setLoading(true);
      try {
        const response = await fetch(`http://localhost:8080/v1/room-types/${id}`, {
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include'
        });

        if (!response.ok) {
          throw new Error('Failed to fetch room type');
        }

        const data = await response.json();
        setRoomType(data.data);
      } catch (error) {
        console.error('Error fetching room type:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchRoomType();
  }, [id]);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setRoomType({ ...roomType, [name]: value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    try {
      const updatedFields = {};
      if (roomType.name !== '') updatedFields.name = roomType.name;
      if (roomType.description !== '') updatedFields.description = roomType.description;
      if (roomType.capacity !== '') updatedFields.capacity = Number(roomType.capacity);
      if (roomType.default_price !== '') updatedFields.default_price = Number(roomType.default_price);

      const response = await fetch(`http://localhost:8080/v1/room-types/`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          id: roomType.id,
          ...updatedFields
        }),
        credentials: 'include'
      });

      if (!response.ok) {
        throw new Error('Failed to update room type');
      }

      const data = await response.json();
      console.log('Room type updated:', data);
      navigate('/room-type');
    } catch (error) {
      console.error('Error updating room type:', error);
      alert(error.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Box className="form-container">
      <Typography variant="h4" gutterBottom className="form-title">
        Edit Room Type
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
                {loading ? 'Updating...' : 'Update Room Type'}
              </Button>
            </Grid>
          </Grid>
        </form>
      )}
    </Box>
  );
};

export default RoomTypeEdit;
