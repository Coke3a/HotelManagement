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
            'Access-Control-Expose-Headers': 'X-My-Custom-Header, X-Another-Custom-Header'
          },
          credentials: 'include'
        });

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
    <div className="mt-5 p-5 rounded-lg bg-gray-100">
      <div className="flex justify-end mb-4">
        <Button
          variant="contained"
          className="bg-green-500 text-white hover:bg-green-600"
          onClick={handleAddBooking}
        >
          Add Booking
        </Button>
      </div>
      <div className="flex flex-wrap gap-4 mb-4">
        <TextField
          label="Booking ID"
          value={filters.id}
          onChange={(e) => setFilters({ ...filters, id: e.target.value })}
          className="w-40"
        />
        <TextField
          label="Customer ID"
          value={filters.customer_id}
          onChange={(e) => setFilters({ ...filters, customer_id: e.target.value })}
          className="w-40"
        />
        <LocalizationProvider dateAdapter={AdapterDateFns}>
          <DatePicker
            label="Check-in Date"
            value={filters.check_in_date}
            onChange={(date) => setFilters({ ...filters, check_in_date: date })}
            renderInput={(params) => <TextField {...params} className="w-40" />}
          />
          <DatePicker
            label="Check-out Date"
            value={filters.check_out_date}
            onChange={(date) => setFilters({ ...filters, check_out_date: date })}
            renderInput={(params) => <TextField {...params} className="w-40" />}
          />
        </LocalizationProvider>
        <TextField
          label="Total Amount"
          value={filters.total_amount}
          onChange={(e) => setFilters({ ...filters, total_amount: e.target.value })}
          className="w-40"
        />
        <FormControl variant="outlined" className="w-40">
          <InputLabel>Filter by Status</InputLabel>
          <Select
            value={statusFilter}
            onChange={(e) => setStatusFilter(e.target.value)}
            label="Filter by Status"
          >
            <MenuItem value="">All</MenuItem>
            <MenuItem value="1">Pending</MenuItem>
            <MenuItem value="2">Confirmed</MenuItem>
            <MenuItem value="3">Checked In</MenuItem>
            <MenuItem value="4">Checked Out</MenuItem>
            <MenuItem value="5">Canceled</MenuItem>
            <MenuItem value="6">Completed</MenuItem>
          </Select>
        </FormControl>
      </div>
      <TableContainer component={Paper} style={{ fontSize: '0.9rem' }}>
        <Table size="small">
          <TableHead className="bg-blue-600">
            <TableRow>
              <TableCell className="text-white font-bold">ID</TableCell>
              <TableCell className="text-white font-bold" style={{ paddingRight: '8px' }}>Guest</TableCell>
              <TableCell className="text-white font-bold" align="center" style={{ paddingLeft: '8px' }}>Room</TableCell>
              <TableCell className="text-white font-bold">Check-in</TableCell>
              <TableCell className="text-white font-bold">Check-out</TableCell>
              <TableCell className="text-white font-bold">Amount</TableCell>
              <TableCell className="text-white font-bold">Payment</TableCell>
              <TableCell className="text-white font-bold">Status</TableCell>
              <TableCell className="text-white font-bold">Actions</TableCell>
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
