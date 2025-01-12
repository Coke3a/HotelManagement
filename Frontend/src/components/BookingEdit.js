import React, { useState, useEffect } from 'react';
import { TextField, Button, Box, Typography, CircularProgress, Select, MenuItem, FormControl, InputLabel, Grid, Dialog, DialogTitle, DialogContent, DialogActions } from '@mui/material';
import { useNavigate, useParams } from 'react-router-dom';
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns';
import { LocalizationProvider, DatePicker } from '@mui/x-date-pickers';
import GuestAdd from './GuestAdd';
import ReceiptLongIcon from '@mui/icons-material/ReceiptLong';
import GuestSearch from './GuestSearch';
import { PaymentStatus, PaymentMethod } from '../utils/paymentEnums';

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
  const [customerIdentityNumber, setCustomerIdentityNumber] = useState('');
  const [payment, setPayment] = useState({
    id: '',
    amount: '',
    payment_method: '',
    status: '',
    payment_date: null
  });
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchData = async () => {
      try {
        await fetchBookingData();
        await fetchPaymentData(id);
      } catch (error) {
        console.error('Error fetching data:', error);
        setError('Failed to fetch data');
      }
    };

    fetchData();
  }, [id]);

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
      setCustomerIdentityNumber(data.data.customer_identity_number);
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

      // Ensure customer_id and rate_price_id are set
      if (data.data.customer_id && data.data.rate_price_id) {
        setBooking(prev => ({
          ...prev,
          customer_id: data.data.customer_id,
          rate_price_id: data.data.rate_price_id,
        }));
      }

      // Calculate stay duration
      const diffTime = Math.abs(new Date(data.data.check_out_date) - new Date(data.data.check_in_date));
      const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
      setStayDuration(diffDays);

      await Promise.all([
        fetchRatePrices(data.data.room_type_id)
      ]);

      const initialGuest = {
        id: data.data.customer_id,
        firstname: data.data.customer_firstname,
        surname: data.data.customer_surname,
        identity_number: data.data.customer_identity_number,
        phone: data.data.customer_phone || '',
        email: data.data.customer_email || ''
      };
      setGuests([initialGuest]);

      fetchPaymentData(data.data.booking_id);
    } catch (error) {
      console.error('Error fetching booking data:', error);
      alert(error.message);
    } finally {
      setLoading(false);
    }
  };

  const fetchGuests = async () => {
    try {
      const queryParams = new URLSearchParams({
        skip: '0',
        limit: '100',
        identity_number: customerIdentityNumber
      });

      const response = await fetch(`http://localhost:8080/v1/customers/?${queryParams.toString()}`, {
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
      if (data.success && data.data && data.data.customers) {
        setGuests(data.data.customers);
        // If we found the guest, update the booking's customer_id
        if (data.data.customers.length > 0) {
          const foundGuest = data.data.customers[0];
          setBooking(prev => ({ ...prev, customer_id: foundGuest.id }));
        }
      }
    } catch (error) {
      console.error('Error fetching guests:', error);
    }
  };

  // Add useEffect to trigger guest search when customerIdentityNumber changes
  useEffect(() => {
    if (customerIdentityNumber) {
      fetchGuests();
    }
  }, [customerIdentityNumber]);

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

      const data = await response.json();
      if (!data.success) {
        setRatePrices([]);
        setBooking(prev => ({ ...prev, rate_price_id: '', total_amount: prev.total_amount }));
        alert(data.messages?.[0] || 'No rate prices found for this room type');
        return;
      }

      if (data && data.data && data.data.ratePrices && Array.isArray(data.data.ratePrices)) {
        setRatePrices(data.data.ratePrices);
      } else {
        setRatePrices([]);
        setBooking(prev => ({ ...prev, rate_price_id: '', total_amount: prev.total_amount }));
        alert('No rate prices available for this room type');
      }
    } catch (error) {
      console.error('Error fetching rate prices:', error);
      setRatePrices([]);
      setBooking(prev => ({ ...prev, rate_price_id: '', total_amount: prev.total_amount }));
      alert('Error fetching rate prices. Please try again.');
    }
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    if (name === 'rate_price_id') {
      const selectedRatePrice = ratePrices.find(rp => rp.id === parseInt(value));
      if (selectedRatePrice) {
        const totalAmount = (selectedRatePrice.price_per_night * stayDuration).toFixed(2);
        setBooking(prev => ({
          ...prev,
          [name]: value,
          total_amount: totalAmount
        }));
      } else {
        setBooking(prev => ({ ...prev, [name]: value }));
      }
    } else {
      setBooking(prev => ({ ...prev, [name]: value }));
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    try {
      // Update booking
      const bookingResponse = await fetch(`http://localhost:8080/v1/booking/`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify({
          id: parseInt(booking.booking_id),
          customer_id: parseInt(booking.customer_id),
          rate_prices_id: booking.rate_price_id ? parseInt(booking.rate_price_id) : null,
          room_id: parseInt(booking.room_id),
          room_type_id: parseInt(booking.room_type_id),
          check_in_date: booking.check_in_date.toISOString(),
          check_out_date: booking.check_out_date.toISOString(),
          status: parseInt(booking.status),
          total_amount: parseFloat(booking.total_amount),
        }),
      });

      if (!bookingResponse.ok) {
        throw new Error('Failed to update booking');
      }

      // Update payment if it exists
      if (payment.id) {
        const paymentResponse = await fetch(`http://localhost:8080/v1/payments/`, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`,
          },
          body: JSON.stringify({
            id: parseInt(payment.id),
            amount: parseFloat(payment.amount),
            payment_method: parseInt(payment.payment_method),
            status: parseInt(payment.status),
          }),
        });

        if (!paymentResponse.ok) {
          throw new Error('Failed to update payment');
        }
      }

      navigate('/booking');
    } catch (error) {
      console.error('Error updating booking/payment:', error);
      setError(error.message);
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
    // Update guests list with the new guest
    setGuests(prevGuests => [...prevGuests, newGuest]);
    // Select the newly added guest
    setBooking(prev => ({ ...prev, customer_id: newGuest.id }));
    // Set the selected guest in GuestSearch component
    const guestSearchComponent = document.querySelector('GuestSearch');
    if (guestSearchComponent) {
      guestSearchComponent.setSelectedGuest(newGuest);
    }
  };

  const handleStayDurationChange = (event) => {
    const duration = parseInt(event.target.value);
    setStayDuration(duration);
    
    // Recalculate total amount if rate price is selected
    if (booking.rate_price_id) {
      const selectedRatePrice = ratePrices.find(rp => rp.id === parseInt(booking.rate_price_id));
      if (selectedRatePrice) {
        const totalAmount = (selectedRatePrice.price_per_night * duration).toFixed(2);
        setBooking(prev => ({
          ...prev,
          total_amount: totalAmount
        }));
      }
    }

    // Update check-out date based on new duration
    if (booking.check_in_date) {
      const checkOutDate = new Date(booking.check_in_date);
      checkOutDate.setDate(checkOutDate.getDate() + duration);
      setBooking(prev => ({ ...prev, check_out_date: checkOutDate }));
    }
  };

  const isFormValid = () => {
    return booking.customer_id && 
           booking.check_in_date && 
           booking.check_out_date && 
           booking.room_id && 
           booking.rate_price_id;
  };

  const fetchPaymentData = async (bookingId) => {
    try {
      const response = await fetch(`http://localhost:8080/v1/payments/${bookingId}`, {
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
      if (data.success && data.data) {
        setPayment({
          id: data.data.id,
          amount: data.data.amount,
          payment_method: data.data.payment_method,
          status: data.data.status,
          payment_date: data.data.payment_date
        });
      }
    } catch (error) {
      console.error('Error fetching payment:', error);
    }
  };

  return (
    <Box className="form-container">
      <Typography variant="h4" gutterBottom className="form-title">
        Edit Booking
      </Typography>
      <Box display="flex" justifyContent="space-between" alignItems="center" marginBottom={2}>
        <Button
          variant="contained"
          color="primary"
          onClick={() => navigate(`/receipt/${id}`)}
          startIcon={<ReceiptLongIcon />}
          style={{ textTransform: 'none', fontWeight: 'bold' }}
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
              <FormControl fullWidth margin="normal" required className="form-input">
                <LocalizationProvider dateAdapter={AdapterDateFns}>
                  <DatePicker
                    label="Check-in Date"
                    value={booking.check_in_date}
                    onChange={() => {}}
                    renderInput={(params) => <TextField {...params} fullWidth />}
                    inputFormat="dd/MM/yyyy"
                    disabled
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
                    onChange={() => {}}
                    renderInput={(params) => <TextField {...params} fullWidth />}
                    inputFormat="dd/MM/yyyy"
                    disabled
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
              {booking.total_amount && (
                <TextField
                  fullWidth
                  label="Total Amount"
                  name="total_amount"
                  value={parseInt(booking.total_amount)}
                  onChange={handleChange}
                  margin="normal"
                  className="form-input"
                  InputProps={{
                    readOnly: true,
                  }}
                />
              )}
            </Grid>
            <Grid item xs={12}>
              <GuestSearch
                onGuestSelected={(guestId) => setBooking(prev => ({ ...prev, customer_id: guestId }))}
                currentGuestId={booking.customer_id}
                guests={guests}
                setGuests={setGuests}
                onAddNewGuest={handleAddNewGuest}
                initialSelectedGuest={guests[0]}
              />
            </Grid>
            <Grid item xs={12}>
              <FormControl fullWidth margin="normal" required className="form-input">
                <InputLabel>Status</InputLabel>
                <Select
                  name="status"
                  value={booking.status}
                  onChange={handleChange}
                >
                  <MenuItem value={1}>Uncheck-in</MenuItem>
                  <MenuItem value={2}>Check-in</MenuItem>
                  <MenuItem value={3}>Check-out</MenuItem>
                  <MenuItem value={4}>Cancelled</MenuItem>
                  <MenuItem value={5}>Completed</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} sm={6}>
              <FormControl fullWidth margin="normal" required className="form-input">
                <InputLabel>Payment Status</InputLabel>
                <Select
                  name="payment_status"
                  value={payment.status}
                  onChange={(e) => setPayment({ ...payment, status: e.target.value })}
                >
                  <MenuItem value={PaymentStatus.UNPAID}>Unpaid</MenuItem>
                  <MenuItem value={PaymentStatus.PAID}>Paid</MenuItem>
                  <MenuItem value={PaymentStatus.FAILED}>Failed</MenuItem>
                  <MenuItem value={PaymentStatus.REFUNDED}>Refunded</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} sm={6}>
              <FormControl fullWidth margin="normal" required className="form-input">
                <InputLabel>Payment Method</InputLabel>
                <Select
                  name="payment_method"
                  value={payment.payment_method}
                  onChange={(e) => setPayment({ ...payment, payment_method: e.target.value })}
                >
                  <MenuItem value={PaymentMethod.NOT_SPECIFIED}>Not specified</MenuItem>
                  <MenuItem value={PaymentMethod.CREDIT_CARD}>Credit Card</MenuItem>
                  <MenuItem value={PaymentMethod.DEBIT_CARD}>Debit Card</MenuItem>
                  <MenuItem value={PaymentMethod.CASH}>Cash</MenuItem>
                  <MenuItem value={PaymentMethod.BANK_TRANSFER}>Bank Transfer</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12}>
              <Button
                type="submit"
                variant="contained"
                color="primary"
                className="form-submit"
                disabled={!isFormValid() || loading}
                sx={{ mt: 3, mb: 2 }}
              >
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
