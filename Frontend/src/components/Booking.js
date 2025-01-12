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
  Stack,
  Pagination,
  Chip,
  Grid,
} from '@mui/material';
import { LocalizationProvider, MobileDatePicker } from '@mui/x-date-pickers';
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns';
import { getPaymentStatusMessage, getPaymentStatusColor } from '../utils/paymentEnums';
import SearchIcon from '@mui/icons-material/Search';
import { getBookingStatusMessage, getBookingStatusColor } from '../utils/bookingStatusEnums';
import { PaymentStatus } from '../utils/paymentEnums';

const Booking = () => {
  const token = localStorage.getItem('token');
  const navigate = useNavigate();
  const [bookings, setBookings] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [statusFilter, setStatusFilter] = useState('');
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(25);
  const [totalCount, setTotalCount] = useState(0);
  const [filters, setFilters] = useState({
    booking_id: '',
    booking_price: '',
    booking_status: '',
    check_in_date: null,
    check_out_date: null,
    room_number: '',
    room_type_name: '',
    customer_firstname: '',
    customer_surname: '',
    payment_status: '',
    created_at: null,
  });

  const fetchBookings = async (page, rowsPerPage) => {
    setLoading(true);
    setError(null);
    try {
      const skip = page * rowsPerPage;
      const queryParams = new URLSearchParams({
        skip: skip.toString(),
        limit: rowsPerPage.toString(),
        ...(filters.booking_id && { booking_id: filters.booking_id }),
        ...(filters.booking_price && { booking_price: filters.booking_price }),
        ...(filters.booking_status && { booking_status: filters.booking_status }),
        ...(filters.check_in_date && { check_in_date: filters.check_in_date.toISOString().split('T')[0] }),
        ...(filters.check_out_date && { check_out_date: filters.check_out_date.toISOString().split('T')[0] }),
        ...(filters.created_at && { created_at: filters.created_at.toISOString().split('T')[0] }),
        ...(filters.room_number && { room_number: filters.room_number }),
        ...(filters.room_type_name && { room_type_name: filters.room_type_name }),
        ...(filters.customer_firstname && { customer_firstname: filters.customer_firstname }),
        ...(filters.customer_surname && { customer_surname: filters.customer_surname }),
        ...(filters.payment_status && { payment_status: filters.payment_status }),
      });

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
      if (data && data.data && (data.data.booking_customer_payments || data.data.bookings)) {
        setBookings(data.data.booking_customer_payments || data.data.bookings);
        if (data.data.meta && typeof data.data.meta.total === 'number') {
          setTotalCount(data.data.meta.total);
        } else {
          setTotalCount((data.data.booking_customer_payments || data.data.bookings).length);
        }
      } else {
        setBookings([]);
        setTotalCount(0);
      }
    } catch (error) {
      console.error('Error fetching bookings:', error);
      setError(error.message);
      setBookings([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchBookings(page, rowsPerPage);
  }, [page, rowsPerPage]);

  const handlePageChange = (event, newPage) => {
    setPage(newPage - 1);
  };

  const handleGuestClick = (guestId) => {
    navigate(`/guest/edit/${guestId}`);
  };

  const handleRoomClick = (roomId) => {
    navigate(`/room/edit/${roomId}`);
  };

  const handleEditBooking = (id) => {
    navigate(`/booking/edit/${id}`);
  };

  const handleAddBooking = () => {
    navigate('/booking/add');
  };

  const getStatusString = (status, type = 'booking') => {
    if (type === 'booking') {
      switch (status) {
        case 1: return 'Uncheck-in';
        case 2: return 'Check-in';
        case 3: return 'Check-out';
        case 4: return 'Cancelled';
        case 5: return 'Completed';
        default: return 'Unknown';
      }
    } else if (type === 'payment') {
      return getPaymentStatusMessage(status);
    }
    return 'Unknown';
  };

  const renderBookingStatus = (status) => {
    const statusMessage = getBookingStatusMessage(status);
    const { chipColor } = getBookingStatusColor(status);

    return (
      <Chip 
        label={statusMessage} 
        color={chipColor}
        size="small"
      />
    );
  };

  const renderPaymentStatus = (status, bookingId) => {
    const statusMessage = getPaymentStatusMessage(status);
    const { backgroundColor, textColor } = getPaymentStatusColor(status);

    return (
      <Chip 
        label={statusMessage}
        sx={{
          backgroundColor: backgroundColor,
          color: textColor,
          fontWeight: 'medium',
          '&:hover': {
            backgroundColor: backgroundColor,
            opacity: 0.9
          }
        }}
        size="small"
      />
    );
  };

  return (
    <div className="table-container">
      <Box sx={{ 
        display: 'flex', 
        justifyContent: 'space-between', 
        alignItems: 'center',
        mb: 2
      }}>
        <Button
          variant="contained"
          className="bg-blue-600 text-white hover:bg-blue-700"
          onClick={handleAddBooking}
        >
          Add Booking
        </Button>
      </Box>

      <Box sx={{ mb: 2 }}>
        <Paper elevation={1} sx={{ p: 1.5, backgroundColor: '#f8f9fa' }}>
          <Grid container spacing={1} alignItems="center">
            <Grid item xs={12} md={10}>
              <Stack direction={{ xs: 'column', md: 'row' }} spacing={1}>
                {/* ID and Price Fields */}
                <Stack direction="row" spacing={1} sx={{ width: { xs: '100%', md: '30%' } }}>
                  <TextField
                    size="small"
                    placeholder="Booking ID"
                    value={filters.booking_id}
                    onChange={(e) => setFilters({ ...filters, booking_id: e.target.value })}
                    sx={{ width: '45%' }}
                  />
                  <TextField
                    size="small"
                    placeholder="Price"
                    value={filters.booking_price}
                    onChange={(e) => setFilters({ ...filters, booking_price: e.target.value })}
                    sx={{ width: '55%' }}
                  />
                </Stack>

                {/* Date Fields */}
                <Stack direction="row" spacing={1} sx={{ width: { xs: '100%', md: '40%' } }}>
                  <TextField
                    type="date"
                    size="small"
                    value={filters.check_in_date ? filters.check_in_date.toISOString().split('T')[0] : ''}
                    onChange={(e) => setFilters({ ...filters, check_in_date: new Date(e.target.value) })}
                    sx={{ width: '50%' }}
                    InputLabelProps={{ shrink: true }}
                    placeholder="Check-in"
                  />
                  <TextField
                    type="date"
                    size="small"
                    value={filters.check_out_date ? filters.check_out_date.toISOString().split('T')[0] : ''}
                    onChange={(e) => setFilters({ ...filters, check_out_date: new Date(e.target.value) })}
                    sx={{ width: '50%' }}
                    InputLabelProps={{ shrink: true }}
                    placeholder="Check-out"
                  />
                </Stack>

                {/* Status Fields */}
                <Stack direction="row" spacing={1} sx={{ width: { xs: '100%', md: '30%' } }}>
                  <FormControl size="small" sx={{ width: '50%' }}>
                    <Select
                      value={filters.booking_status}
                      onChange={(e) => setFilters({ ...filters, booking_status: e.target.value })}
                      displayEmpty
                      sx={{ height: '40px' }}
                    >
                      <MenuItem value="">All Status</MenuItem>
                      <MenuItem value="1">Uncheck-in</MenuItem>
                      <MenuItem value="2">Check-in</MenuItem>
                      <MenuItem value="3">Check-out</MenuItem>
                      <MenuItem value="4">Cancelled</MenuItem>
                      <MenuItem value="5">Completed</MenuItem>
                    </Select>
                  </FormControl>
                  <FormControl size="small" sx={{ width: '50%' }}>
                    <Select
                      value={filters.payment_status}
                      onChange={(e) => setFilters({ ...filters, payment_status: e.target.value })}
                      displayEmpty
                      sx={{ height: '40px' }}
                    >
                      <MenuItem value="">All Payments</MenuItem>
                      <MenuItem value={PaymentStatus.UNPAID}>Unpaid</MenuItem>
                      <MenuItem value={PaymentStatus.PAID}>Paid</MenuItem>
                      <MenuItem value={PaymentStatus.FAILED}>Failed</MenuItem>
                      <MenuItem value={PaymentStatus.REFUNDED}>Refunded</MenuItem>
                    </Select>
                  </FormControl>
                </Stack>
              </Stack>

              {/* Second Row */}
              <Stack direction={{ xs: 'column', md: 'row' }} spacing={1} sx={{ mt: 1 }}>
                <TextField
                  size="small"
                  placeholder="Room Number"
                  value={filters.room_number}
                  onChange={(e) => setFilters({ ...filters, room_number: e.target.value })}
                  sx={{ width: { xs: '100%', md: '25%' } }}
                />
                <TextField
                  size="small"
                  placeholder="Room Type"
                  value={filters.room_type_name}
                  onChange={(e) => setFilters({ ...filters, room_type_name: e.target.value })}
                  sx={{ width: { xs: '100%', md: '25%' } }}
                />
                <TextField
                  size="small"
                  placeholder="Guest First Name"
                  value={filters.customer_firstname}
                  onChange={(e) => setFilters({ ...filters, customer_firstname: e.target.value })}
                  sx={{ width: { xs: '100%', md: '25%' } }}
                />
                <TextField
                  size="small"
                  placeholder="Guest Last Name"
                  value={filters.customer_surname}
                  onChange={(e) => setFilters({ ...filters, customer_surname: e.target.value })}
                  sx={{ width: { xs: '100%', md: '25%' } }}
                />
              </Stack>
            </Grid>

            {/* Buttons */}
            <Grid item xs={12} md={2}>
              <Stack direction="row" spacing={1}>
                <Button
                  variant="contained"
                  onClick={() => {
                    setPage(0);
                    fetchBookings(0, rowsPerPage);
                  }}
                  startIcon={<SearchIcon />}
                  size="small"
                  fullWidth
                >
                  Search
                </Button>
                <Button
                  variant="outlined"
                  onClick={() => {
                    setFilters({
                      booking_id: '',
                      booking_price: '',
                      booking_status: '',
                      check_in_date: null,
                      check_out_date: null,
                      room_number: '',
                      room_type_name: '',
                      customer_firstname: '',
                      customer_surname: '',
                      payment_status: '',
                      created_at: null,
                    });
                    setStatusFilter('');
                    setPage(0);
                    fetchBookings(0, rowsPerPage);
                  }}
                  size="small"
                  sx={{ minWidth: 'auto', width: '40px', p: 0 }}
                >
                  ×
                </Button>
              </Stack>
            </Grid>
          </Grid>
        </Paper>
      </Box>

      <Stack 
        direction="row" 
        spacing={2} 
        alignItems="center"
        sx={{ marginLeft: 'auto' }}
      >
        <Box sx={{ mr: 2 }}>
          <select
            value={rowsPerPage}
            onChange={(e) => {
              setRowsPerPage(Number(e.target.value));
              setPage(0);
            }}
            style={{
              padding: '5px',
              borderRadius: '4px',
              border: '1px solid #ccc'
            }}
          >
            <option value={25}>25 per page</option>
            <option value={50}>50 per page</option>
            <option value={75}>75 per page</option>
          </select>
        </Box>
        
        <Pagination
          count={Math.ceil(totalCount / rowsPerPage)}
          page={page + 1}
          onChange={handlePageChange}
          color="primary"
          showFirstButton
          showLastButton
          siblingCount={1}
          boundaryCount={1}
          sx={{
            '& .MuiPagination-ul': {
              flexWrap: 'nowrap'
            }
          }}
        />
      </Stack>
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
                <TableRow key={booking.id || booking.booking_id} className={index % 2 === 0 ? "bg-gray-50" : "bg-white"}>
                  <TableCell>{booking.id || booking.booking_id}</TableCell>
                  <TableCell style={{ paddingRight: '8px' }}>
                    <Button
                      size="small"
                      color="primary"
                      onClick={() => handleGuestClick(booking.customer_id)}
                    >
                      {booking.customer_firstname ? 
                        `${booking.customer_firstname} ${booking.customer_surname}` :
                        `Guest ${booking.customer_id}`
                      }
                    </Button>
                  </TableCell>
                  <TableCell align="center" style={{ paddingLeft: '8px' }}>
                    <Box display="flex" flexDirection="row" alignItems="center" justifyContent="">
                      <Button
                        size="small"
                        color="primary"
                        onClick={() => handleRoomClick(booking.room_id)}
                      >
                        {booking.room_number || `Room ${booking.room_id}`}
                        {booking.room_type_name && (
                          <Typography variant="caption" color="textSecondary" style={{ marginLeft: '4px' }}>
                            {booking.room_type_name}
                          </Typography>
                        )}
                      </Button>
                    </Box>
                  </TableCell>
                  <TableCell>{new Date(booking.check_in_date).toLocaleDateString('en-GB')}</TableCell>
                  <TableCell>{new Date(booking.check_out_date).toLocaleDateString('en-GB')}</TableCell>
                  <TableCell>{booking.total_amount || booking.booking_price}</TableCell>
                  <TableCell>
                    {booking.payment_status ? 
                      renderPaymentStatus(booking.payment_status, booking.id || booking.booking_id) 
                      : '-'
                    }
                  </TableCell>
                  <TableCell>{renderBookingStatus(booking.status || booking.booking_status)}</TableCell>
                  <TableCell>
                    <Button
                      variant="outlined"
                      color="primary"
                      size="small"
                      onClick={() => handleEditBooking(booking.id || booking.booking_id)}
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
