import React, { useState, useEffect } from 'react';
import { TextField, Button, Box, Typography, CircularProgress, Select, MenuItem, FormControl, InputLabel, Grid } from '@mui/material';
import { useNavigate, useParams } from 'react-router-dom';
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns';
import { LocalizationProvider, DatePicker } from '@mui/x-date-pickers';
import { PaymentStatus, PaymentMethod } from '../utils/paymentEnums';

const PaymentEdit = () => {
  const token = localStorage.getItem('token');
  const navigate = useNavigate();
  const { id } = useParams();
  const [payment, setPayment] = useState({
    id: '',
    booking_id: '',
    amount: '',
    payment_method: '',
    payment_date: null,
    status: '',
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchPayment = async () => {
      setLoading(true);
      try {
        const response = await fetch(`http://localhost:8080/v1/payments/${id}`, {
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
          throw new Error('Failed to fetch payment');
        }

        const data = await response.json();
        setPayment({
          ...data.data,
          payment_date: new Date(data.data.payment_date),
          payment_method: data.data.payment_method === "Not specified" ? PaymentMethod.NOT_SPECIFIED : data.data.payment_method,
        });
      } catch (error) {
        console.error('Error fetching payment:', error);
        setError('Failed to fetch payment data. Please try again.');
      } finally {
        setLoading(false);
      }
    };

    fetchPayment();
  }, [id]);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setPayment({ ...payment, [name]: value });
  };

  const handleDateChange = (date) => {
    setPayment({ ...payment, payment_date: date });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    try {
      const response = await fetch(`http://localhost:8080/v1/payments/`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify({
          ...payment,
          payment_date: payment.payment_date.toISOString(),
          amount: parseFloat(payment.amount),
          status: parseInt(payment.status),
          payment_method: parseInt(payment.payment_method), // Ensure this is sent as an integer
        }),
      });

      if (response.status === 401) {
        handleTokenExpiration(new Error("access token has expired"), navigate);
        return;
      }

      if (!response.ok) {
        throw new Error('Failed to update payment');
      }

      const data = await response.json();
      console.log('Updated payment data:', data);
      navigate('/payment');
    } catch (error) {
      console.error('Error updating payment:', error);
      setError('Failed to update payment. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Box className="form-container">
      <Typography variant="h4" gutterBottom className="form-title">
        Edit Payment
      </Typography>
      {loading ? (
        <Box display="flex" justifyContent="center" alignItems="center" height="300px">
          <CircularProgress />
        </Box>
      ) : (
        <form onSubmit={handleSubmit} className="form">
          {error && (
            <Typography color="error" gutterBottom>
              {error}
            </Typography>
          )}
          <Grid container spacing={2}>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Payment ID"
                name="id"
                value={payment.id}
                margin="normal"
                disabled
                className="form-input"
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Booking ID"
                name="booking_id"
                value={payment.booking_id}
                margin="normal"
                disabled
                className="form-input"
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Amount"
                name="amount"
                value={payment.amount}
                onChange={handleChange}
                margin="normal"
                type="number"
                step="0.01"
                required
                className="form-input"
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <FormControl fullWidth margin="normal" required className="form-input">
                <InputLabel>Payment Method</InputLabel>
                <Select
                  name="payment_method"
                  value={payment.payment_method}
                  onChange={handleChange}
                >
                  <MenuItem value={PaymentMethod.NOT_SPECIFIED}>Not specified</MenuItem>
                  <MenuItem value={PaymentMethod.CREDIT_CARD}>Credit Card</MenuItem>
                  <MenuItem value={PaymentMethod.DEBIT_CARD}>Debit Card</MenuItem>
                  <MenuItem value={PaymentMethod.CASH}>Cash</MenuItem>
                  <MenuItem value={PaymentMethod.BANK_TRANSFER}>Bank Transfer</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} sm={6}>
              <LocalizationProvider dateAdapter={AdapterDateFns}>
                <DatePicker
                  label="Payment Date"
                  value={payment.payment_date}
                  onChange={handleDateChange}
                  renderInput={(params) => <TextField {...params} fullWidth margin="normal" required className="form-input" />}
                  inputFormat="yyyy-MM-dd"
                />
              </LocalizationProvider>
            </Grid>
            <Grid item xs={12} sm={6}>
              <FormControl fullWidth margin="normal" required className="form-input">
                <InputLabel>Status</InputLabel>
                <Select
                  name="status"
                  value={payment.status}
                  onChange={handleChange}
                >
                  <MenuItem value={PaymentStatus.PENDING}>Pending</MenuItem>
                  <MenuItem value={PaymentStatus.COMPLETED}>Completed</MenuItem>
                  <MenuItem value={PaymentStatus.FAILED}>Failed</MenuItem>
                  <MenuItem value={PaymentStatus.REFUNDED}>Refunded</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12}>
              <Button type="submit" variant="contained" color="primary" fullWidth disabled={loading}>
                {loading ? 'Updating Payment...' : 'Update Payment'}
              </Button>
            </Grid>
          </Grid>
        </form>
      )}
    </Box>
  );
};

export default PaymentEdit;
