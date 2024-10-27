import React, { useState, useEffect } from 'react';
import { TextField, Button, Box, Typography, CircularProgress, Select, MenuItem, FormControl, InputLabel, Grid, Dialog, DialogTitle, DialogContent, DialogActions } from '@mui/material';
import { useNavigate, useParams } from 'react-router-dom';
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns';
import { LocalizationProvider, DatePicker } from '@mui/x-date-pickers';
import GuestAdd from './GuestAdd';

const BookingEdit = () => {
  const token = localStorage.getItem('token');
  const navigate = useNavigate();
  const { id } = useParams();
  const [booking, setBooking] = useState({
    booking_id: '',
    customer_id: '',
    rate_price_id: '',
    room_id: '',
    room_type_id: '',
    room_number: '',
    room_type_name: '',
    check_in_date: null,
    check_out_date: null,
    status: 1,
    total_amount: '',
  });
  const [loading, setLoading] = useState(true);
  const [guests, setGuests] = useState([]);
  const [ratePrices, setRatePrices] = useState([]);
  const [stayDuration, setStayDuration] = useState(1);
  const [openGuestModal, setOpenGuestModal] = useState(false);

  useEffect(() => {
    const fetchBookingData = async () => {
      try {
        const response = await fetch(`http://localhost:8080/v1/booking/${id}/details`, {
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
          throw new Error('Failed to fetch booking data');
        }

        const data = await response.json();
        setBooking({
          booking_id: data.data.booking_id,
          customer_id: data.data.customer_id,
          rate_price_id: data.data.rate_price_id || '',
          room_id: data.data.room_id,
          room_type_id: data.data.room_type_id,
          room_number: data.data.room_number,
          room_type_name: data.data.room_type_name,
          check_in_date: new Date(data.data.check_in_date),
          check_out_date: new Date(data.data.check_out_date),
          status: data.data.booking_status,
          total_amount: data.data.booking_price,
        });

        // Calculate stay duration
        const diffTime = Math.abs(new Date(data.data.check_out_date) - new Date(data.data.check_in_date));
        const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
        setStayDuration(diffDays);

        // Fetch guests, rooms, and rate prices
        await Promise.all([
          fetchGuests(),
          fetchRatePrices(data.data.room_type_id)
        ]);
      } catch (error) {
        console.error('Error fetching booking data:', error);
        alert(error.message);
      } finally {
        setLoading(false);
      }
    };

    fetchBookingData();
  }, [id]);

  const fetchGuests = async () => {
    try {
      const response = await fetch('http://localhost:8080/v1/customers/?skip=0&limit=100', {
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
      });

      if (response.status === 401) {
        handleTokenExpiration(new Error("access token has expired"), navigate);
        return;
      }

      const data = await response.json();
      setGuests(data.data.customers || []);
    } catch (error) {
      console.error('Error fetching guests:', error);
    }
  };

  const fetchRooms = async () => {
    try {
      const response = await fetch('http://localhost:8080/v1/rooms', {
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
      });

      if (response.status === 401) {
        handleTokenExpiration(new Error("access token has expired"), navigate);
        return;
      }

      if (!response.ok) throw new Error('Failed to fetch rooms');

      const data = await response.json();
      setRooms(data.data || []);
    } catch (error) {
      console.error('Error fetching rooms:', error);
      alert(error.message);
    }
  };

  const fetchRatePrices = async (roomTypeId) => {
    try {
      const response = await fetch(`http://localhost:8080/v1/rate_prices/by-room-type/${roomTypeId}`, {
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
      });

      if (response.status === 401) {
        handleTokenExpiration(new Error("access token has expired"), navigate);
        return;
      }

      if (!response.ok) throw new Error('Failed to fetch rate prices');

      const data = await response.json();
      setRatePrices(data.data || []);
    } catch (error) {
      console.error('Error fetching rate prices:', error);
      alert(error.message);
    }
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setBooking({ ...booking, [name]: value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    try {
      const response = await fetch(`http://localhost:8080/v1/booking/`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify({
          id: booking.booking_id,
          customer_id: booking.customer_id,
          rate_prices_id: booking.rate_price_id ? parseInt(booking.rate_price_id) : null,
          room_id: booking.room_id,
          room_type_id: booking.room_type_id,
          check_in_date: booking.check_in_date.toISOString(),
          check_out_date: booking.check_out_date.toISOString(),
          status: parseInt(booking.status),
          total_amount: booking.total_amount,
        }),
      });

      if (response.status === 401) {
        handleTokenExpiration(new Error("access token has expired"), navigate);
        return;
      }

      if (!response.ok) {
        throw new Error('Failed to update booking');
      }

      const data = await response.json();
      console.log('Updated booking data:', data);
      navigate('/booking');
    } catch (error) {
      console.error('Error updating booking:', error);
      alert(error.message);
    } finally {
      setLoading(false);
    }
  };

  const handleAddNewGuest = () => {
    setOpenGuestModal(true);
  };

  const handleCloseGuestModal = () => {
    setOpenGuestModal(false);
  };

  const handleGuestAdded = async (newGuest) => {
    setOpenGuestModal(false);
    await fetchGuests();
    setBooking(prev => ({ ...prev, customer_id: newGuest.id }));
  };

  return (
    <Box className="form-container">
      <Typography variant="h4" gutterBottom className="form-title">
        Edit Booking
      </Typography>
      <Box display="flex" justifyContent="space-between" alignItems="center" marginBottom={2}>
        <Button
          variant="contained"
          color="secondary"
          onClick={() => navigate(`/receipt/${id}`)}
        >
          Print Receipt
        </Button>
      </Box>
      {loading ? (
        <Box display="flex" justifyContent="center" alignItems="center" height="300px">
          <CircularProgress />
        </Box>
      ) : (
        <form onSubmit={handleSubmit} className="form">
          <Grid container spacing={2}>
            <Grid item xs={12} sm={6}>
              <LocalizationProvider dateAdapter={AdapterDateFns}>
                <DatePicker
                  label="Check-in Date"
                  value={booking.check_in_date}
                  onChange={() => {}}
                  renderInput={(params) => <TextField {...params} fullWidth margin="normal" required className="form-input" />}
                  inputFormat="yyyy-MM-dd"
                  disabled
                />
              </LocalizationProvider>
            </Grid>
            <Grid item xs={12} sm={6}>
              <LocalizationProvider dateAdapter={AdapterDateFns}>
                <DatePicker
                  label="Check-out Date"
                  value={booking.check_out_date}
                  onChange={() => {}}
                  renderInput={(params) => <TextField {...params} fullWidth margin="normal" required className="form-input" />}
                  inputFormat="yyyy-MM-dd"
                  disabled
                />
              </LocalizationProvider>
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Stay Duration (nights)"
                type="number"
                value={stayDuration}
                onChange={() => {}}
                inputProps={{ min: 1 }}
                margin="normal"
                className="form-input"
                disabled
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Room"
                value={`${booking.room_number} - ${booking.room_type_name}`}
                margin="normal"
                className="form-input"
                disabled
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <FormControl fullWidth margin="normal" required className="form-input">
                <InputLabel>Rate Price</InputLabel>
                <Select
                  name="rate_price_id"
                  value={booking.rate_price_id}
                  onChange={handleChange}
                >
                  {ratePrices.map((ratePrice) => (
                    <MenuItem key={ratePrice.id} value={ratePrice.id}>
                      {ratePrice.name} - {ratePrice.price_per_night} per night
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Total Amount"
                name="total_amount"
                value={booking.total_amount}
                onChange={handleChange}
                margin="normal"
                className="form-input"
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <FormControl fullWidth margin="normal" required className="form-input">
                <InputLabel>Guest</InputLabel>
                <Select
                  name="customer_id"
                  value={booking.customer_id}
                  onChange={handleChange}
                >
                  {guests.map((guest) => (
                    <MenuItem key={guest.id} value={guest.id}>
                      {guest.firstname} {guest.surname}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} sm={6}>
              <Button
                variant="contained"
                color="secondary"
                onClick={handleAddNewGuest}
                style={{ marginTop: '16px', height: '56px' }}
              >
                Add New Guest
              </Button>
            </Grid>
            <Grid item xs={12}>
              <FormControl fullWidth margin="normal" required className="form-input">
                <InputLabel>Status</InputLabel>
                <Select
                  name="status"
                  value={booking.status}
                  onChange={handleChange}
                >
                  <MenuItem value={1}>Pending</MenuItem>
                  <MenuItem value={2}>Confirmed</MenuItem>
                  <MenuItem value={3}>Checked In</MenuItem>
                  <MenuItem value={4}>Checked Out</MenuItem>
                  <MenuItem value={5}>Canceled</MenuItem>
                  <MenuItem value={6}>Completed</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12}>
              <Button type="submit" variant="contained" color="primary" fullWidth disabled={loading}>
                {loading ? 'Updating Booking...' : 'Update Booking'}
              </Button>
            </Grid>
          </Grid>
        </form>
      )}
      <Dialog open={openGuestModal} onClose={handleCloseGuestModal} maxWidth="md" fullWidth>
        <DialogTitle>Add New Guest</DialogTitle>
        <DialogContent>
          <GuestAdd onGuestAdded={handleGuestAdded} isFromBooking={true} />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseGuestModal} color="primary">
            Cancel
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default BookingEdit;
