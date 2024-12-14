import React, { useState } from 'react';
import { TextField, Button, Box, Typography, CircularProgress, Select, MenuItem, FormControl, InputLabel } from '@mui/material';
import { useNavigate } from 'react-router-dom';

const PaymentAdd = () => {
  const token = localStorage.getItem('token');
  const navigate = useNavigate();
  const [payment, setPayment] = useState({
    booking_id: '',
    amount: '',
    payment_method: '',
    status: '',
  });
  const [loading, setLoading] = useState(false);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setPayment({ ...payment, [name]: value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    try {
      const response = await fetch('http://localhost:8080/v1/payments/', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify(payment),
      });

      if (response.status === 401) {
        handleTokenExpiration(new Error("access token has expired"), navigate);
        return;
      }

      if (!response.ok) {
        throw new Error('Failed to create payment');
      }

      const data = await response.json();
      console.log('Payment created:', data);
      navigate('/payment');
    } catch (error) {
      console.error('Error adding payment:', error);
      alert(error.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Box className="form-container">
      <Typography variant="h4" gutterBottom className="form-title">
        Add New Payment
      </Typography>
      {loading ? (
        <Box display="flex" justifyContent="center" alignItems="center" height="300px">
          <CircularProgress />
        </Box>
      ) : (
        <form onSubmit={handleSubmit} className="form">
          <TextField
            fullWidth
            label="Booking ID"
            name="booking_id"
            value={payment.booking_id}
            onChange={handleChange}
            margin="normal"
            type="number"
            required
            className="form-input"
          />
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
          <FormControl fullWidth margin="normal" required className="form-input">
            <InputLabel>Payment Method</InputLabel>
            <Select
              name="payment_method"
              value={payment.payment_method}
              onChange={handleChange}
            >
              <MenuItem value={PaymentMethod.CREDIT_CARD}>Credit Card</MenuItem>
              <MenuItem value={PaymentMethod.DEBIT_CARD}>Debit Card</MenuItem>
              <MenuItem value={PaymentMethod.CASH}>Cash</MenuItem>
              <MenuItem value={PaymentMethod.BANK_TRANSFER}>Bank Transfer</MenuItem>
            </Select>
          </FormControl>
          <FormControl fullWidth margin="normal" required className="form-input">
            <InputLabel>Status</InputLabel>
            <Select
              name="status"
              value={payment.status}
              onChange={handleChange}
            >
              <MenuItem value="unpaid">Unpaid</MenuItem>
              <MenuItem value="paid">Paid</MenuItem>
            </Select>
          </FormControl>
          <Button type="submit" variant="contained" color="primary" className="form-submit" disabled={loading}>
            {loading ? 'Adding...' : 'Add Payment'}
          </Button>
        </form>
      )}
    </Box>
  );
};

export default PaymentAdd;
