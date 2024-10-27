import React, { useState, useEffect } from 'react';
import { TextField, Button, Box, Typography, CircularProgress, Select, MenuItem, FormControl, InputLabel, Grid, Paper } from '@mui/material';
import { useNavigate } from 'react-router-dom';

const RoomAdd = () => {
  const token = localStorage.getItem('token');
  const navigate = useNavigate();
  const [room, setRoom] = useState({
    room_number: '',
    type_id: '',
    description: '',
    status: '',
    floor: '',
  });
  const [loading, setLoading] = useState(false);
  const [roomTypes, setRoomTypes] = useState([]);

  useEffect(() => {
    const fetchRoomTypes = async () => {
      try {
        const response = await fetch('http://localhost:8080/v1/room-types/', {
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
          throw new Error('Failed to fetch room types');
        }

        const data = await response.json();
        setRoomTypes(data.data || []);
      } catch (error) {
        console.error('Error fetching room types:', error);
      }
    };

    fetchRoomTypes();
  }, []);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setRoom({ ...room, [name]: value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    try {
      const response = await fetch('http://localhost:8080/v1/rooms/', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify({
          ...room,
          status: parseInt(room.status),
          floor: parseInt(room.floor),
        }),
      });

      if (response.status === 409) {
        throw new Error('Room already exists. Please use a different room number.');
      }

      if (response.status === 401) {
        handleTokenExpiration(new Error("access token has expired"), navigate);
        return;
      }

      if (!response.ok) {
        throw new Error('Failed to create room');
      }

      const data = await response.json();
      console.log('Room created:', data);
      navigate('/room');
    } catch (error) {
      console.error('Error adding room:', error);
      alert(error.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Box sx={{ maxWidth: '800px', margin: '0 auto', padding: '20px' }}>
      <Paper elevation={3} sx={{ padding: '24px', backgroundColor: '#f8f9fa' }}>
        <Typography variant="h5" gutterBottom className="form-title">
          Add New Room
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
                  label="Room Number"
                  name="room_number"
                  value={room.room_number}
                  onChange={handleChange}
                  required
                  className="form-input"
                  size="small"
                />
              </Grid>
              <Grid item xs={12} md={6}>
                <FormControl fullWidth required size="small">
                  <InputLabel>Room Type</InputLabel>
                  <Select
                    name="type_id"
                    value={room.type_id}
                    onChange={handleChange}
                    label="Room Type"
                  >
                    {roomTypes.map((type) => (
                      <MenuItem key={type.id} value={type.id}>
                        {type.name}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
              </Grid>
              <Grid item xs={12} md={6}>
                <FormControl fullWidth required size="small">
                  <InputLabel>Status</InputLabel>
                  <Select
                    name="status"
                    value={room.status}
                    onChange={handleChange}
                  >
                    <MenuItem value={1}>Available</MenuItem>
                    <MenuItem value={2}>Maintenance</MenuItem>
                  </Select>
                </FormControl>
              </Grid>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label="Floor"
                  name="floor"
                  value={room.floor}
                  onChange={handleChange}
                  required
                  className="form-input"
                  size="small"
                  type="number"
                />
              </Grid>
              <Grid item xs={12}>
                <TextField
                  fullWidth
                  label="Description"
                  name="description"
                  value={room.description}
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
                {loading ? 'Adding...' : 'Add Room'}
              </Button>
            </Box>
          </form>
        )}
      </Paper>
    </Box>
  );
};

export default RoomAdd;
