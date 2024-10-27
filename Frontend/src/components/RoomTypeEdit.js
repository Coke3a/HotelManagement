import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import {
  Box,
  Typography,
  TextField,
  Button,
  CircularProgress,
  Grid,
  Paper,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
} from '@mui/material';

const RoomTypeEdit = () => {
  const token = localStorage.getItem('token');
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
            'Authorization': `Bearer ${token}`,
          },
        });

        if (response.status === 401) {
          handleTokenExpiration(new Error("access token has expired"), navigate);
          return;
        }

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
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify({
          id: roomType.id,
          ...updatedFields
        }),
      });

      if (response.status === 401) {
        handleTokenExpiration(new Error("access token has expired"), navigate);
        return;
      }

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

  const handleDelete = async () => {
    if (window.confirm('Are you sure you want to delete this room type?')) {
      setLoading(true);
      try {
        const response = await fetch(`http://localhost:8080/v1/room-types/${id}`, {
          method: 'DELETE',
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json',
          },
        });

        if (response.status === 401) {
          handleTokenExpiration(new Error("access token has expired"), navigate);
          return;
        }

        if (!response.ok) {
          throw new Error('Failed to delete room type');
        }

        console.log('Room type deleted');
        navigate('/room-type');
      } catch (error) {
        console.error('Error deleting room type:', error);
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
          Edit Room Type
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
                variant="outlined"
                color="error"
                onClick={handleDelete}
                disabled={loading}
                size="medium"
              >
                Delete Room Type
              </Button>
              <Button
                type="submit"
                variant="contained"
                color="primary"
                disabled={loading}
                size="medium"
              >
                {loading ? 'Updating...' : 'Update Room Type'}
              </Button>
            </Box>
          </form>
        )}
      </Paper>
    </Box>
  );
};

export default RoomTypeEdit;
