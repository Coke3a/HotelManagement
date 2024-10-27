import React, { useState, useEffect } from 'react';
import { TextField, Button, Box, Typography, CircularProgress, Select, MenuItem, FormControl, InputLabel, Grid, Paper } from '@mui/material';
import { useNavigate, useParams } from 'react-router-dom';

const RatePriceEdit = () => {
  const token = localStorage.getItem('token');
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
              'Authorization': `Bearer ${token}`,
              'Content-Type': 'application/json',
            },
          }),
          fetch('http://localhost:8080/v1/room-types/?skip=0&limit=100', {
            headers: {
              'Content-Type': 'application/json',
              'Authorization': `Bearer ${token}`,
            },
          })
        ]);

        if (ratePriceResponse.status === 401) {
          handleTokenExpiration(new Error("access token has expired"), navigate);
          return;
        }

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
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          id: parseInt(ratePrice.id),
          name: ratePrice.name,
          description: ratePrice.description,
          price_per_night: parseFloat(ratePrice.price_per_night),
          room_type_id: parseInt(ratePrice.room_type_id),
        }),
      });

      if (response.status === 401) {
        handleTokenExpiration(new Error("access token has expired"), navigate);
        return;
      }

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
          headers: {
            'Authorization': `Bearer ${token}`,
          },
        });

        if (response.status === 401) {
          handleTokenExpiration(new Error("access token has expired"), navigate);
          return;
        }

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
    <Box sx={{ maxWidth: '800px', margin: '0 auto', padding: '20px' }}>
      <Paper elevation={3} sx={{ padding: '24px', backgroundColor: '#f8f9fa' }}>
        <Typography variant="h5" gutterBottom className="form-title">
          Edit Rate Price
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
                variant="outlined"
                color="error"
                onClick={handleDelete}
                disabled={loading}
                size="medium"
              >
                Delete Rate Price
              </Button>
              <Button
                type="submit"
                variant="contained"
                color="primary"
                disabled={loading}
                size="medium"
              >
                {loading ? 'Updating...' : 'Update Rate Price'}
              </Button>
            </Box>
          </form>
        )}
      </Paper>
    </Box>
  );
};

export default RatePriceEdit;
