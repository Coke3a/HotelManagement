import React, { useState, useEffect } from 'react';
import { TextField, Button, Box, Typography, CircularProgress, Select, MenuItem, FormControl, InputLabel, Grid } from '@mui/material';
import { useNavigate, useParams } from 'react-router-dom';

const RatePriceEdit = () => {
  const navigate = useNavigate();
  const { id } = useParams();
  const [ratePrice, setRatePrice] = useState({
    id: '',
    name: '',
    description: '',
    price_per_night: '',
    room_type_id: '',
  });
  const [roomTypes, setRoomTypes] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      try {
        const [ratePriceResponse, roomTypesResponse] = await Promise.all([
          fetch(`http://localhost:8080/v1/rate_prices/${id}`, {
            headers: {
              'Content-Type': 'application/json',
            },
            credentials: 'include',
          }),
          fetch('http://localhost:8080/v1/room-types/?skip=0&limit=100', {
            headers: {
              'Content-Type': 'application/json',
            },
            credentials: 'include',
          })
        ]);

        if (!ratePriceResponse.ok) {
          throw new Error('Failed to fetch rate price');
        }
        if (!roomTypesResponse.ok) {
          throw new Error('Failed to fetch room types');
        }

        const ratePriceData = await ratePriceResponse.json();
        const roomTypesData = await roomTypesResponse.json();

        setRatePrice(ratePriceData.data);
        setRoomTypes(roomTypesData.data || []);
      } catch (error) {
        console.error('Error fetching data:', error);
        alert(error.message);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [id]);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setRatePrice({ ...ratePrice, [name]: value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    try {
      const response = await fetch(`http://localhost:8080/v1/rate_prices/`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          id: parseInt(ratePrice.id),
          name: ratePrice.name,
          description: ratePrice.description,
          price_per_night: parseFloat(ratePrice.price_per_night),
          room_type_id: parseInt(ratePrice.room_type_id),
        }),
        credentials: 'include',
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.messages[0] || 'Failed to update rate price');
      }

      const data = await response.json();
      console.log('Updated rate price:', data);
      navigate('/rate_price');
    } catch (error) {
      console.error('Error updating rate price:', error);
      alert(error.message);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async () => {
    if (window.confirm('Are you sure you want to delete this rate price?')) {
      setLoading(true);
      try {
        const response = await fetch(`http://localhost:8080/v1/rate_prices/${id}`, {
          method: 'DELETE',
          credentials: 'include',
        });

        if (!response.ok) {
          throw new Error('Failed to delete rate price');
        }

        console.log('Deleted rate price:', id);
        navigate('/rate_price');
      } catch (error) {
        console.error('Error deleting rate price:', error);
        alert(error.message);
      } finally {
        setLoading(false);
      }
    }
  };

  return (
    <Box className="form-container">
      <Typography variant="h4" gutterBottom className="form-title">
        Edit Rate Price
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
                label="Rate Price ID"
                name="id"
                value={ratePrice.id}
                onChange={handleChange}
                margin="normal"
                required
                type="number"
                className="form-input"
                disabled
              />
            </Grid>
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
                      {roomType.id} - {roomType.name} - Capacity: {roomType.capacity} - Default Price: {roomType.default_price}
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
                step="0.01"
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
              <Box display="flex" justifyContent="space-between" mt={2}>
                <Button
                  variant="contained"
                  color="error"
                  onClick={handleDelete}
                  disabled={loading}
                >
                  Delete Rate Price
                </Button>
                <Button type="submit" variant="contained" color="primary" disabled={loading}>
                  {loading ? 'Updating...' : 'Update Rate Price'}
                </Button>
              </Box>
            </Grid>
          </Grid>
        </form>
      )}
    </Box>
  );
};

export default RatePriceEdit;