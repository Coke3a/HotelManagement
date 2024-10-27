import React, { useState, useEffect } from 'react';
import { TextField, Button, Box, Typography, CircularProgress, Select, MenuItem, FormControl, InputLabel, Grid } from '@mui/material';
import { useNavigate } from 'react-router-dom';

const RatePriceAdd = () => {
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
            'Content-Type': 'application/json',
          },
          credentials: 'include'
        });

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
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          ...ratePrice,
          price_per_night: parseFloat(ratePrice.price_per_night),
        }),
        credentials: 'include'
      });

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
    <Box className="form-container">
      <Typography variant="h4" gutterBottom className="form-title">
        Add New Rate Price
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
                value={ratePrice.name}
                onChange={handleChange}
                margin="normal"
                required
                className="form-input"
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <FormControl fullWidth margin="normal" required className="form-input">
                <InputLabel>Room Type</InputLabel>
                <Select
                  name="room_type_id"
                  value={ratePrice.room_type_id}
                  onChange={handleChange}
                >
                  {roomTypes.map((roomType) => (
                    <MenuItem key={roomType.id} value={roomType.id}>
                      {roomType.id} - {roomType.name}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Price Per Night"
                name="price_per_night"
                value={ratePrice.price_per_night}
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
                value={ratePrice.description}
                onChange={handleChange}
                margin="normal"
                multiline
                rows={3}
                className="form-input"
              />
            </Grid>
            <Grid item xs={12}>
              <Button type="submit" variant="contained" color="primary" className="form-submit" disabled={loading}>
                {loading ? 'Adding...' : 'Add Rate Price'}
              </Button>
            </Grid>
          </Grid>
        </form>
      )}
    </Box>
  );
};

export default RatePriceAdd;