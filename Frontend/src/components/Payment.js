import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Button,
  Box,
  CircularProgress,
  Typography,
} from '@mui/material';
import { getPaymentStatusMessage, getPaymentMethodMessage } from '../utils/paymentEnums';

const Payment = () => {
  const navigate = useNavigate();
  const [payments, setPayments] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchPayments = async () => {
      setLoading(true);
      try {
        const response = await fetch('http://localhost:8080/v1/payments/?skip=0&limit=10', {
          headers: {
            'Content-Type': 'application/json',
            'Access-Control-Expose-Headers': 'X-My-Custom-Header, X-Another-Custom-Header',
          },
          credentials: 'include',
        });

        if (!response.ok) {
          throw new Error('Failed to fetch payments');
        }

        const data = await response.json();
        setPayments(data.data.payments || []);
      } catch (error) {
        console.error('Error fetching payments:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchPayments();
  }, []);

  const handleEditPayment = (id) => {
    navigate(`/payment/edit/${id}`);
  };

  const handleBookingClick = (bookingId) => {
    navigate(`/booking/edit/${bookingId}`);
  };

  return (
    <div>
      <TableContainer component={Paper} style={{ fontSize: '0.9rem' }}>
        <Table size="small">
          <TableHead className="bg-blue-600">
            <TableRow>
              <TableCell className="text-white font-bold">Booking ID</TableCell>
              <TableCell className="text-white font-bold">Payment ID</TableCell>
              <TableCell className="text-white font-bold">Amount</TableCell>
              <TableCell className="text-white font-bold">Method</TableCell>
              <TableCell className="text-white font-bold">Date</TableCell>
              <TableCell className="text-white font-bold">Status</TableCell>
              <TableCell className="text-white font-bold">Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {loading ? (
              <TableRow>
                <TableCell colSpan={7} align="center">
                  <Box display="flex" justifyContent="center" alignItems="center" height="50px">
                    <CircularProgress />
                  </Box>
                </TableCell>
              </TableRow>
            ) : payments.length === 0 ? (
              <TableRow>
                <TableCell colSpan={7} align="center">
                  <Typography variant="body1">No payments found</Typography>
                </TableCell>
              </TableRow>
            ) : (
              payments.map((payment) => (
                <TableRow key={payment.id}>
                  <TableCell>
                    <Button
                      color="primary"
                      onClick={() => handleBookingClick(payment.booking_id)}
                    >
                      {payment.booking_id}
                    </Button>
                  </TableCell>
                  <TableCell>{payment.id}</TableCell>
                  <TableCell>{payment.amount}</TableCell>
                  <TableCell>{getPaymentMethodMessage(payment.payment_method)}</TableCell>
                  <TableCell>{new Date(payment.payment_date).toLocaleDateString()}</TableCell>
                  <TableCell>
                    {getPaymentStatusMessage(payment.status)}
                  </TableCell>
                  <TableCell>
                    <Button
                      variant="outlined"
                      color="primary"
                      size="small"
                      onClick={() => handleEditPayment(payment.id)}
                    >
                      Edit
                    </Button>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </TableContainer>
    </div>
  );
};

export default Payment;
