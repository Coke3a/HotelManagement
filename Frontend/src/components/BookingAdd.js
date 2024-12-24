import React, { useState, useEffect } from 'react';
import { TextField, Button, Box, Typography, CircularProgress, Select, MenuItem, FormControl, InputLabel, Grid } from '@mui/material';
import { useNavigate } from 'react-router-dom';
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns';
import { LocalizationProvider, DatePicker } from '@mui/x-date-pickers';
import Dialog from '@mui/material/Dialog';
import DialogTitle from '@mui/material/DialogTitle';
import DialogContent from '@mui/material/DialogContent';
import DialogActions from '@mui/material/DialogActions';
import GuestAdd from './GuestAdd';

const BookingAdd = () => {
  const token = localStorage.getItem('token');
  const navigate = useNavigate();
  const [booking, setBooking] = useState({
    customer_id: '',
    rate_prices_id: '',
    room_id: '',
    room_type_id: '',
    check_in_date: null,
    check_out_date: null,
    status: 1,
    total_amount: '',
  });
  const [loading, setLoading] = useState(false);
  const [guests, setGuests] = useState([]);
  const [availableRooms, setAvailableRooms] = useState([]);
  const [ratePrices, setRatePrices] = useState([]);
  const [stayDuration, setStayDuration] = useState(1);
  const [selectedRoom, setSelectedRoom] = useState(null);
  const [roomError, setRoomError] = useState('');
  const [openGuestModal, setOpenGuestModal] = useState(false);
  
  const fetchAvailableRooms = async () => {
    if (!booking.check_in_date || !booking.check_out_date) return;
    setLoading(true);
    try {
      const checkInDate = booking.check_in_date.toLocaleDateString('en-CA', { 
        timeZone: 'Asia/Bangkok'
      });

      const checkOutDate = booking.check_out_date.toLocaleDateString('en-CA', {
        timeZone: 'Asia/Bangkok'
      });

      const response = await fetch(`http://localhost:8080/v1/rooms/available?check_in_date=${checkInDate}&check_out_date=${checkOutDate}`, {
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
      });

      if (response.status === 401) {
        handleTokenExpiration(new Error("access token has expired"), navigate);
        return;
      }

      if (!response.ok) throw new Error('Failed to fetch available rooms');

      const data = await response.json();
      if (data.success && data.data && data.data.length > 0) {
        setAvailableRooms(data.data);
        setRoomError('');
      } else {
        setAvailableRooms([]);
        setRoomError('No rooms available for the selected dates');
      }
    } catch (error) {
      console.error('Error fetching available rooms:', error);
      setRoomError('Error fetching available rooms');
    } finally {
      setLoading(false);
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

      const data = await response.json();
      if (!data.success) {
        if (data.messages !== "data not found") {
          alert('No rate prices found for this room type');
        } else {
          alert(data.messages?.[0] || 'No rate prices found for this room type');
        }

        setRatePrices([]);
        setBooking(prev => ({ ...prev, rate_prices_id: '', total_amount: '' }));
        return;
      }

      if (data && data.data && data.data.ratePrices && Array.isArray(data.data.ratePrices)) {
        setRatePrices(data.data.ratePrices);
      } else {
        setRatePrices([]);
        setBooking(prev => ({ ...prev, rate_prices_id: '', total_amount: '' }));
        alert('No rate prices available for this room type');
      }
    } catch (error) {
      console.error('Error fetching rate prices:', error);
      setRatePrices([]);
      setBooking(prev => ({ ...prev, rate_prices_id: '', total_amount: '' }));
      alert('Error fetching rate prices. Please try again.');
    }
  };

  const calculateStayDuration = () => {
    if (booking.check_in_date && booking.check_out_date) {
      const diffTime = Math.abs(booking.check_out_date - booking.check_in_date);
      const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
      setStayDuration(diffDays);
    }
  };

  useEffect(() => {
    calculateStayDuration();
    if (booking.check_in_date && booking.check_out_date) {
      fetchAvailableRooms();
    } else {
      setAvailableRooms([]);
    }
  }, [booking.check_in_date, booking.check_out_date]);

  const handleDateChange = (name) => (date) => {
    if (name === 'check_in_date') {
      const checkOutDate = new Date(date);
      checkOutDate.setDate(checkOutDate.getDate() + stayDuration);
      setBooking({ ...booking, [name]: date, check_out_date: checkOutDate });
    } else if (name === 'check_out_date') {
      setBooking({ ...booking, [name]: date });
      if (booking.check_in_date) {
        const diffTime = Math.abs(date - booking.check_in_date);
        const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
        setStayDuration(diffDays);
      }
    }
    // Reset selected room and rate prices when dates change
    setSelectedRoom(null);
    setRatePrices([]);
    setBooking(prev => ({ ...prev, rate_prices_id: '', total_amount: '' }));
  };

  const handleStayDurationChange = (event) => {
    const duration = parseInt(event.target.value);
    setStayDuration(duration);
    
    // Recalculate total amount if rate price is selected
    if (booking.rate_prices_id) {
      const selectedRatePrice = ratePrices.find(rp => rp.id === parseInt(booking.rate_prices_id));
      if (selectedRatePrice) {
        const totalAmount = (selectedRatePrice.price_per_night * duration).toFixed(2);
        setBooking(prev => ({
          ...prev,
          total_amount: totalAmount
        }));
      }
    }

    if (booking.check_in_date) {
      const checkOutDate = new Date(booking.check_in_date);
      checkOutDate.setDate(checkOutDate.getDate() + duration);
      setBooking(prev => ({ ...prev, check_out_date: checkOutDate }));
    }
  };

  const handleRoomChange = async (e) => {
    const roomId = e.target.value;
    const selected = availableRooms.find(room => room.id === roomId);
    setSelectedRoom(selected);
    
    if (selected && selected.room_type_id) {
      await fetchRatePrices(selected.room_type_id);
    } else {
      console.error('Selected room or room_type_id is undefined');
    }
    
    setBooking(prev => ({ ...prev, rate_prices_id: '', total_amount: '' }));
  };

  const handleRatePriceChange = (e) => {
    const ratePriceId = e.target.value;
    const selectedRatePrice = ratePrices.find(rp => rp.id === parseInt(ratePriceId));
    if (selectedRatePrice) {
      const totalAmount = (selectedRatePrice.price_per_night * stayDuration).toFixed(2);
      setBooking({
        ...booking,
        rate_prices_id: ratePriceId,
        total_amount: totalAmount
      });
    }
  };

  const handleGuestChange = (e) => {
    const { name, value } = e.target;
    setBooking(prev => {
      const newBooking = { ...prev, [name]: value };
      if (name === 'rate_prices_id') {
        const selectedRatePrice = ratePrices.find(rp => rp.id === parseInt(value));
        if (selectedRatePrice) {
          newBooking.total_amount = parseFloat(selectedRatePrice.price_per_night).toFixed(2);
        }
      }
      return newBooking;
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    try {
      // Convert dates to ISO format with Bangkok timezone offset
      const checkInDate = new Date(booking.check_in_date);
      checkInDate.setHours(17, 0, 0, 0); // Set to 17:00:00.000
      
      const checkOutDate = new Date(booking.check_out_date);
      checkOutDate.setHours(17, 0, 0, 0); // Set to 17:00:00.000

      const response = await fetch('http://localhost:8080/v1/booking/', {
        method: 'POST',
        headers: {  
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify({
          customer_id: parseInt(booking.customer_id),
          rate_prices_id: parseInt(booking.rate_prices_id),
          room_id: parseInt(selectedRoom.id),
          room_type_id: parseInt(selectedRoom.room_type_id),
          check_in_date: checkInDate.toISOString(), // Will output format like 2024-12-21T17:00:00.000Z
          check_out_date: checkOutDate.toISOString(), // Will output format like 2024-12-21T17:00:00.000Z
          status: parseInt(booking.status),
          total_amount: parseFloat(booking.total_amount),
        }),
      });

      if (response.status === 401) {
        handleTokenExpiration(new Error("access token has expired"), navigate);
        return;
      }

      if (!response.ok) {
        throw new Error('Failed to create booking');
      }

      const data = await response.json();
      console.log('Booking created:', data);
      navigate('/booking');
    } catch (error) {
      console.error('Error adding booking:', error);
      alert(error.message);
    } finally {
      setLoading(false);
    }
  };

  const today = new Date();
  today.setHours(0, 0, 0, 0);

  useEffect(() => {
    const fetchData = async () => {
      try {
        await fetchGuests();
        const ratePricesResponse = await fetch('http://localhost:8080/v1/rate_prices/?skip=0&limit=100', {
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`,
          },
        });

        if (ratePricesResponse.status === 401) {
          handleTokenExpiration(new Error("access token has expired"), navigate);
          return;
        }

        const ratePricesData = await ratePricesResponse.json();
        if (ratePricesData && ratePricesData.data && ratePricesData.data.ratePrices) {
          setRatePrices(ratePricesData.data.ratePrices);
        } else {
          setRatePrices([]);
        }
      } catch (error) {
        console.error('Error fetching data:', error);
        setRatePrices([]);
      }
    };

    fetchData();
  }, []);

  const handleAddNewGuest = () => {
    setOpenGuestModal(true);
  };

  const handleCloseGuestModal = () => {
    setOpenGuestModal(false);
  };

  const handleGuestAdded = async (newGuest) => {
    setOpenGuestModal(false);
    // Fetch the updated list of guests
    await fetchGuests();
    // Select the newly added guest
    setBooking(prev => ({ ...prev, customer_id: newGuest.id }));
  };

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

  return (
    <Box className="form-container">
      <Typography variant="h4" gutterBottom className="form-title">
        Add New Booking
      </Typography>
      {loading ? (
        <Box display="flex" justifyContent="center" alignItems="center" height="300px">
          <CircularProgress />
        </Box>
      ) : (
        <form onSubmit={handleSubmit} className="form">
          <Grid container spacing={2}>
            <Grid item xs={12} sm={6}>
              <FormControl fullWidth margin="normal" required className="form-input">
                <LocalizationProvider dateAdapter={AdapterDateFns}>
                  <DatePicker
                    label="Check-in Date"
                    value={booking.check_in_date}
                    onChange={(date) => handleDateChange('check_in_date')(date)}
                    renderInput={(params) => <TextField {...params} fullWidth />}
                    inputFormat="yyyy-MM-dd"
                    minDate={today}
                    timezone="Asia/Bangkok"
                  />
                </LocalizationProvider>
              </FormControl>
            </Grid>
            <Grid item xs={12} sm={6}>
              <FormControl fullWidth margin="normal" required className="form-input">
                <LocalizationProvider dateAdapter={AdapterDateFns}>
                  <DatePicker
                    label="Check-out Date"
                    value={booking.check_out_date}
                    onChange={(date) => handleDateChange('check_out_date')(date)}
                    renderInput={(params) => <TextField {...params} fullWidth />}
                    inputFormat="yyyy-MM-dd"
                    minDate={booking.check_in_date ? new Date(booking.check_in_date.getTime() + 86400000) : new Date(today.getTime() + 86400000)}
                    timezone="Asia/Bangkok"
                  />
                </LocalizationProvider>
              </FormControl>
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Stay Duration (nights)"
                type="number"
                value={stayDuration}
                onChange={handleStayDurationChange}
                inputProps={{ min: 1 }}
                margin="normal"
                className="form-input"
                disabled={!booking.check_in_date}
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <FormControl 
                fullWidth 
                margin="normal" 
                required 
                className="form-input" 
                disabled={availableRooms.length === 0}
                error={!!roomError}
              >
                <InputLabel>Available Rooms</InputLabel>
                <Select
                  value={selectedRoom ? selectedRoom.id : ''}
                  onChange={handleRoomChange}
                >
                  {availableRooms && availableRooms.length > 0 ? (
                    availableRooms.map((room) => (
                      <MenuItem key={room.id} value={room.id}>
                        {`${room.room_number} - ชั้น ${room.floor} - ${room.room_type_name}`}
                      </MenuItem>
                    ))
                  ) : (
                    <MenuItem disabled value="">
                      No rooms available
                    </MenuItem>
                  )}
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} sm={6}>
              <FormControl fullWidth margin="normal" required className="form-input" disabled={!selectedRoom || !ratePrices || ratePrices.length === 0}>
                <InputLabel>Rate Price</InputLabel>
                <Select
                  name="rate_prices_id"
                  value={booking.rate_prices_id}
                  onChange={handleRatePriceChange}
                >
                  {ratePrices && ratePrices.map((ratePrice) => (
                    <MenuItem key={ratePrice.id} value={ratePrice.id}>
                      {ratePrice.name} - {ratePrice.price_per_night} per night
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} sm={6}>
              {booking.total_amount && (
                <TextField
                  fullWidth
                  label="Total Amount"
                  name="total_amount"
                  value={parseInt(booking.total_amount)}
                  margin="normal"
                  InputProps={{
                    readOnly: true,
                  }}
                  className="form-input"
                />
              )}
            </Grid>
            <Grid item xs={12} sm={6}>
              <FormControl fullWidth margin="normal" required className="form-input">
                <InputLabel>Guest</InputLabel>
                <Select
                  name="customer_id"
                  value={booking.customer_id}
                  onChange={handleGuestChange}
                >
                  {guests && guests.map((guest) => (
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
              <Button
                type="submit"
                variant="contained"
                color="primary"
                className="form-submit"
                // fullWidth
                disabled={loading}
              >
                {loading ? 'Creating Booking...' : 'Create Booking'}
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

export default BookingAdd;
