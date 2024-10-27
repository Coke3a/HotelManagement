import React, { useState, useEffect } from 'react';
import { TextField, Button, Box, Typography, CircularProgress, Select, MenuItem, FormControl, InputLabel, Grid, Paper } from '@mui/material';
import { useNavigate } from 'react-router-dom';

const RatePriceAdd = () => {
  const token = localStorage.getItem('token');
  const navigate = useNavigate();
  const [ratePrice, setRatePrice] = useState({
    name: '',
    description: '',
    price_per_night: '',
    room_type_id: '',
  });
  const [roomTypes, setRoomTypes] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchRoomTypes = async () => {
      try {
        const response = await fetch('http://localhost:8080/v1/room-types/?skip=0&limit=100', {
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
          throw new Error('Failed to fetch room types');
        }

        const data = await response.json();
        setRoomTypes(data.data || []);
      } catch (error) {
        console.error('Error fetching room types:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchRoomTypes();
  }, []);

  const handleChange = (e) => {
    const { name, value } = e.target;
    if (name === 'room_type_id') {
      const selectedRoomType = roomTypes.find(roomType => roomType.id === parseInt(value));
      setRatePrice({
        ...ratePrice,
        [name]: value,
        price_per_night: selectedRoomType ? selectedRoomType.default_price.toString() : ''
      });
    } else {
      setRatePrice({ ...ratePrice, [name]: value });
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    try {
      const response = await fetch('http://localhost:8080/v1/rate_prices/', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          ...ratePrice,
          price_per_night: parseFloat(ratePrice.price_per_night),
        }),
      });

      if (response.status === 401) {
        handleTokenExpiration(new Error("access token has expired"), navigate);
        return;
      }

      if (!response.ok) {
        throw new Error('Failed to create rate price');
      }

      const data = await response.json();
      console.log('Rate price created:', data);
      navigate('/rate_price');
    } catch (error) {
      console.error('Error adding rate price:', error);
      alert(error.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Box sx={{ maxWidth: '800px', margin: '0 auto', padding: '20px' }}>
      <Paper elevation={3} sx={{ padding: '24px', backgroundColor: '#f8f9fa' }}>
        <Typography variant="h5" gutterBottom className="form-title">
          Add New Rate Price
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
                  value={ratePrice.name}
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
                    name="room_type_id"
                    value={ratePrice.room_type_id}
                    onChange={handleChange}
                    label="Room Type"
                  >
                    {roomTypes.map((roomType) => (
                      <MenuItem key={roomType.id} value={roomType.id}>
                        {roomType.name} - Default Price: {roomType.default_price}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
              </Grid>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label="Price Per Night"
                  name="price_per_night"
                  value={ratePrice.price_per_night}
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
                  value={ratePrice.description}
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
                {loading ? 'Adding...' : 'Add Rate Price'}
              </Button>
            </Box>
          </form>
        )}
      </Paper>
    </Box>
  );
};

export default RatePriceAdd;
