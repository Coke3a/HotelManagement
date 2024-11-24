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
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [totalCount, setTotalCount] = useState(0);
  const [filters, setFilters] = useState({
    id: '',
    customer_id: '',
    check_in_date: null,
    check_out_date: null,
    total_amount: '',
  });

  const fetchBookings = async (page, rowsPerPage) => {
    setLoading(true);
    setError(null);
    try {
      const skip = page * rowsPerPage;
      const queryParams = new URLSearchParams({
        skip: skip.toString(),
        limit: rowsPerPage.toString(),
        ...(filters.id && { id: filters.id }),
        ...(filters.customer_id && { customer_id: filters.customer_id }),
        ...(statusFilter && { status: statusFilter }),
        ...(filters.check_in_date && { check_in_date: filters.check_in_date.toISOString().split('T')[0] }),
        ...(filters.check_out_date && { check_out_date: filters.check_out_date.toISOString().split('T')[0] }),
        ...(filters.total_amount && { total_amount: filters.total_amount }),
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
      if (data && data.data && data.data.booking_customer_payments) {
        setBookings(data.data.booking_customer_payments);
        if (data.data.meta && typeof data.data.meta.total === 'number') {
          setTotalCount(data.data.meta.total);
        } else {
          setTotalCount(data.data.booking_customer_payments.length);
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
  }, [page, rowsPerPage, filters, statusFilter]);

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

  const renderBookingStatus = (status) => {
    const statusMessage = getStatusString(status);
    let statusColor;
    let backgroundColor;
    let textColor;

    switch (status) {
      case 1: // Pending
        backgroundColor = '#FFF4E5';  // Light orange
        textColor = '#663C00';
        break;
      case 2: // Confirmed
        backgroundColor = '#E8F4FD';  // Light blue
        textColor = '#0A4FA6';
        break;
      case 3: // Checked In
        backgroundColor = '#EDF7ED';  // Light green
        textColor = '#1E4620';
        break;
      case 4: // Checked Out
        backgroundColor = '#F5F5F5';  // Light grey
        textColor = '#666666';
        break;
      case 5: // Canceled
        backgroundColor = '#FEEBEE';  // Light red
        textColor = '#7F1D1D';
        break;
      case 6: // Completed
        backgroundColor = '#E5D7FD';  // Light purple
        textColor = '#4A1D96';
        break;
      default:
        backgroundColor = '#F5F5F5';  // Light grey
        textColor = '#666666';
    }

    return (
      <Chip 
        label={statusMessage}
        sx={{
          backgroundColor: backgroundColor,
          color: textColor,
          fontWeight: 'medium',
          '&:hover': {
            backgroundColor: backgroundColor,
          }
        }}
        size="small"
      />
    );
  };

  const renderPaymentStatus = (status) => {
    const statusMessage = getPaymentStatusMessage(parseInt(status));
    let statusColor;

    switch (parseInt(status)) {
      case 1: // Pending
        statusColor = 'warning';
        break;
      case 2: // Completed
        statusColor = 'success';
        break;
      case 3: // Failed
        statusColor = 'error';
        break;
      case 4: // Refunded
        statusColor = 'info';
        break;
      default:
        statusColor = 'default';
    }

    return (
      <Chip 
        label={statusMessage} 
        color={statusColor} 
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
              <option value={5}>5 per page</option>
              <option value={10}>10 per page</option>
              <option value={25}>25 per page</option>
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
      </Box>
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
                  <TableCell>{renderPaymentStatus(booking.payment_status)}</TableCell>
                  <TableCell>{renderBookingStatus(booking.booking_status)}</TableCell>
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
