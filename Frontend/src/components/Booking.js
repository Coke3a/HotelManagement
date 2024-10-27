import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { handleTokenExpiration } from '../utils/api';
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
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  TextField,
  Typography,
} from '@mui/material';
import { LocalizationProvider, DatePicker } from '@mui/x-date-pickers';
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns';
import { getPaymentStatusMessage } from '../utils/paymentEnums';

const Booking = () => {
  const token = localStorage.getItem('token');
  const navigate = useNavigate();
  const [bookings, setBookings] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [statusFilter, setStatusFilter] = useState('');
  const [filters, setFilters] = useState({
    id: '',
    customer_id: '',
    check_in_date: null,
    check_out_date: null,
    total_amount: '',
  });

  const handleGuestClick = (guestId) => {
    navigate(`/guest/edit/${guestId}`);
  };

  const handleRoomClick = (roomId) => {
    navigate(`/room/edit/${roomId}`);
  };

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      setError(null);
      try {
        // Construct query parameters
        const queryParams = new URLSearchParams({
          skip: '0',
          limit: '10',
          ...(filters.id && { id: filters.id }),
          ...(filters.customer_id && { customer_id: filters.customer_id }),
          ...(statusFilter && { status: statusFilter }),
          ...(filters.check_in_date && { check_in_date: filters.check_in_date.toISOString().split('T')[0] }),
          ...(filters.check_out_date && { check_out_date: filters.check_out_date.toISOString().split('T')[0] }),
          ...(filters.total_amount && { total_amount: filters.total_amount }),
        });

        // Fetch bookings with filters
        const response = await fetch(`http://localhost:8080/v1/booking/?${queryParams.toString()}`, {
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
          throw new Error(`Failed to fetch bookings: ${response.status} ${response.statusText}`);
        }

        const data = await response.json();
        setBookings(data.data.booking_customer_payments || []);
      } catch (error) {
        console.error('Error fetching bookings:', error);
        setError(error.message);
        setBookings([]);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [filters, statusFilter]);

  const handleEditBooking = (id) => {
    navigate(`/booking/edit/${id}`);
  };

  const handleAddBooking = () => {
    navigate('/booking/add');
  };

  const getStatusString = (status, type = 'booking') => {
    if (type === 'booking') {
      switch (status) {
        case 1: return 'Pending';
        case 2: return 'Confirmed';
        case 3: return 'Checked In';
        case 4: return 'Checked Out';
        case 5: return 'Canceled';
        case 6: return 'Completed';
        default: return 'Unknown';
      }
    } else if (type === 'payment') {
      return getPaymentStatusMessage(status);
    }
    return 'Unknown';
  };

  return (
    <div className="table-container">
      <div className="flex justify-end mb-4">
        <Button
          variant="contained"
          className="bg-blue-600 text-white hover:bg-blue-700"
          onClick={handleAddBooking}
        >
          Add Booking
        </Button>
      </div>
      <TableContainer component={Paper}>
        <Table size="small" className="compact-table">
          <TableHead>
            <TableRow>
              <TableCell>ID</TableCell>
              <TableCell>Guest</TableCell>
              <TableCell>Room</TableCell>
              <TableCell>Check-in</TableCell>
              <TableCell>Check-out</TableCell>
              <TableCell>Amount</TableCell>
              <TableCell>Payment</TableCell>
              <TableCell>Status</TableCell>
              <TableCell>Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {loading ? (
              <TableRow>
                <TableCell colSpan={9} align="center">
                  <CircularProgress size={24} />
                </TableCell>
              </TableRow>
            ) : error ? (
              <TableRow>
                <TableCell colSpan={9} align="center">
                  <Typography color="error" variant="body2">{error}</Typography>
                </TableCell>
              </TableRow>
            ) : bookings.length === 0 ? (
              <TableRow>
                <TableCell colSpan={9} align="center">
                  <Typography variant="body2">No bookings found</Typography>
                </TableCell>
              </TableRow>
            ) : (
              bookings.map((booking, index) => (
                <TableRow key={booking.booking_id} className={index % 2 === 0 ? "bg-gray-50" : "bg-white"}>
                  <TableCell>{booking.booking_id}</TableCell>
                  <TableCell style={{ paddingRight: '8px' }}>
                    <Button
                      size="small"
                      color="primary"
                      onClick={() => handleGuestClick(booking.customer_id)}
                    >
                      {`${booking.customer_firstname}. ${booking.customer_surname.charAt(0)}`}
                    </Button>
                  </TableCell>
                  <TableCell align="center" style={{ paddingLeft: '8px' }}>
                    <Box display="flex" flexDirection="row" alignItems="center" justifyContent="center">
                      <Button
                        size="small"
                        color="primary"
                        onClick={() => handleRoomClick(booking.room_id)}
                      >
                        {booking.room_number}
                        <Typography variant="caption" color="textSecondary" style={{ marginLeft: '4px' }}>
                        {booking.room_type_name}
                      </Typography>

                      </Button>
                    </Box>
                  </TableCell>
                  <TableCell>{new Date(booking.check_in_date).toLocaleDateString()}</TableCell>
                  <TableCell>{new Date(booking.check_out_date).toLocaleDateString()}</TableCell>
                  <TableCell>{booking.booking_price}</TableCell>
                  <TableCell>{getPaymentStatusMessage(parseInt(booking.payment_status))}</TableCell>
                  <TableCell>{getStatusString(booking.booking_status)}</TableCell>
                  <TableCell>
                    <Button
                      variant="outlined"
                      color="primary"
                      size="small"
                      onClick={() => handleEditBooking(booking.booking_id)}
                    >
                      Detail
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
}

export default Booking;
