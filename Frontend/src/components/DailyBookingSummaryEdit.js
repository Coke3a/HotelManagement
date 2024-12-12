import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import dayjs from 'dayjs';
import { message } from 'antd';
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
  Typography,
  Card,
  CardContent,
  Grid,
  Chip,
  ToggleButtonGroup,
  ToggleButton,
  FormControl,
  Select,
  MenuItem,
} from '@mui/material';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import { BookingStatus, getBookingStatusMessage } from '../utils/bookingStatusEnums';

const DailyBookingSummaryEdit = () => {
  const { date } = useParams();
  const navigate = useNavigate();
  const [summary, setSummary] = useState(null);
  const [createdBookings, setCreatedBookings] = useState([]);
  const [completedBookings, setCompletedBookings] = useState([]);
  const [canceledBookings, setCanceledBookings] = useState([]);
  const [loading, setLoading] = useState(false);
  const [statusLoading, setStatusLoading] = useState(false);
  const token = localStorage.getItem('token');

  const fetchBookings = async (filterParams, skip = 0, limit = 100) => {
    const queryParams = new URLSearchParams({
      ...filterParams,
      skip: skip.toString(),
      limit: limit.toString(),
    });

    try {
      const response = await fetch(`http://localhost:8080/v1/booking/?${queryParams.toString()}`, {
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
      });

      if (response.status === 401) {
        handleTokenExpiration(new Error("access token has expired"), navigate);
        return [];
      }

      const data = await response.json();
      return Array.isArray(data.data?.booking_customer_payments) ? data.data.booking_customer_payments : [];
    } catch (error) {
      console.error('Error fetching bookings:', error);
      return [];
    }
  };

  const fetchSummary = async () => {
    setLoading(true);
    try {
      const response = await fetch(`http://localhost:8080/v1/daily-summary/?date=${date}`, {
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
      });

      if (response.status === 401) {
        handleTokenExpiration(new Error("access token has expired"), navigate);
        return;
      }

      const result = await response.json();
      setSummary(result);

      // Fetch created bookings
      const createdBookings = await fetchBookings({ created_at: date }, 0, 100);
      setCreatedBookings(createdBookings);

      // Fetch completed bookings
      const completedBookings = await fetchBookings({ updated_at: date, booking_status: BookingStatus.COMPLETED }, 0, 100);
      setCompletedBookings(completedBookings);

      // Fetch canceled bookings
      const canceledBookings = await fetchBookings({ updated_at: date, booking_status: BookingStatus.CANCELED }, 0, 100);
      setCanceledBookings(canceledBookings);

    } catch (error) {
      console.error('Error fetching summary:', error);
      message.error('Failed to fetch summary');
    }
    setLoading(false);
  };

  const renderBookingStatus = (status) => {
    const statusMessage = getBookingStatusMessage(status);
    let backgroundColor;
    let textColor;

    switch (status) {
      case BookingStatus.PENDING:
        backgroundColor = '#FFF4E5';
        textColor = '#663C00';
        break;
      case BookingStatus.CONFIRMED:
        backgroundColor = '#E8F4FD';
        textColor = '#0A4FA6';
        break;
      case BookingStatus.CHECKED_IN:
        backgroundColor = '#EDF7ED';
        textColor = '#1E4620';
        break;
      case BookingStatus.CHECKED_OUT:
        backgroundColor = '#F5F5F5';
        textColor = '#666666';
        break;
      case BookingStatus.CANCELED:
        backgroundColor = '#FEEBEE';
        textColor = '#7F1D1D';
        break;
      case BookingStatus.COMPLETED:
        backgroundColor = '#E5D7FD';
        textColor = '#4A1D96';
        break;
      default:
        backgroundColor = '#F5F5F5';
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

  const handleStatusChange = async (event) => {
    const newStatus = event.target.value;
    setStatusLoading(true);
    try {
      const response = await fetch(`http://localhost:8080/v1/daily-summary/status?date=${date}&status=${newStatus}`, {
        method: 'PUT',
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
        throw new Error('Failed to update status');
      }

      setSummary(prev => ({ ...prev, Status: parseInt(newStatus) }));
      message.success('Status updated successfully');
    } catch (error) {
      console.error('Error updating status:', error);
      message.error('Failed to update status');
    } finally {
      setStatusLoading(false);
    }
  };

  useEffect(() => {
    if (date) {
      fetchSummary();
    }
  }, [date]);

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="200px">
        <CircularProgress />
      </Box>
    );
  }

  if (!summary) {
    return (
      <Box p={3}>
        <Typography variant="h6">No summary found for this date.</Typography>
      </Box>
    );
  }

  const renderBookingsTable = (bookings, title) => (
    <Grid item xs={12}>
      <Card>
        <CardContent>
          <Typography variant="h6" gutterBottom>{title}</Typography>
          <TableContainer>
            <Table size="small">
              <TableHead>
                <TableRow sx={{ backgroundColor: '#f5f5f5' }}>
                  <TableCell>Booking ID</TableCell>
                  <TableCell>Customer</TableCell>
                  <TableCell>Room</TableCell>
                  <TableCell>Check In</TableCell>
                  <TableCell>Check Out</TableCell>
                  <TableCell>Amount</TableCell>
                  <TableCell>Status</TableCell>
                  <TableCell>Action</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {bookings.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={8} align="center">
                      <Typography variant="body2">No bookings found</Typography>
                    </TableCell>
                  </TableRow>
                ) : (
                  bookings.map((booking, index) => (
                    <TableRow key={booking.booking_id} className={index % 2 === 0 ? "bg-gray-50" : "bg-white"}>
                      <TableCell>{booking.booking_id}</TableCell>
                      <TableCell>
                        <Button
                          size="small"
                          color="primary"
                          onClick={() => navigate(`/guest/edit/${booking.customer_id}`)}
                        >
                          {booking.customer_firstname ? 
                            `${booking.customer_firstname}. ${booking.customer_surname.charAt(0)}` :
                            `Guest ${booking.customer_id}`
                          }
                        </Button>
                      </TableCell>
                      <TableCell>{booking.room_number}</TableCell>
                      <TableCell>{dayjs(booking.check_in_date).format('YYYY-MM-DD')}</TableCell>
                      <TableCell>{dayjs(booking.check_out_date).format('YYYY-MM-DD')}</TableCell>
                      <TableCell>{booking.booking_price}</TableCell>
                      <TableCell>{renderBookingStatus(booking.booking_status)}</TableCell>
                      <TableCell>
                        <Button
                          variant="outlined"
                          size="small"
                          onClick={() => navigate(`/booking/edit/${booking.booking_id}`)}
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
        </CardContent>
      </Card>
    </Grid>
  );

  return (
    <Box p={3}>
      <Button
        startIcon={<ArrowBackIcon />}
        onClick={() => navigate('/daily-summary')}
        sx={{ mb: 3 }}
        variant="outlined"
      >
        Back
      </Button>

      <Grid container spacing={3}>
        <Grid item xs={12}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>Summary for {date}</Typography>
              <Grid container spacing={2}>
                <Grid item xs={12} sm={6} md={4}>
                  <Typography variant="subtitle2" color="textSecondary">Created Bookings</Typography>
                  <Typography variant="body1">
                    {summary.CreatedBookings ? summary.CreatedBookings.split(';').filter(id => id).length : 0}
                  </Typography>
                </Grid>
                <Grid item xs={12} sm={6} md={4}>
                  <Typography variant="subtitle2" color="textSecondary">Completed Bookings</Typography>
                  <Typography variant="body1">
                    {summary.CompletedBookings ? summary.CompletedBookings.split(';').filter(id => id).length : 0}
                  </Typography>
                </Grid>
                <Grid item xs={12} sm={6} md={4}>
                  <Typography variant="subtitle2" color="textSecondary">Canceled Bookings</Typography>
                  <Typography variant="body1">
                    {summary.CanceledBookings ? summary.CanceledBookings.split(';').filter(id => id).length : 0}
                  </Typography>
                </Grid>
                <Grid item xs={12} sm={6} md={4}>
                  <Typography variant="subtitle2" color="textSecondary">Total Amount</Typography>
                  <Typography variant="body1">
                    {summary.TotalAmount}
                  </Typography>
                </Grid>
                <Grid item xs={12} sm={6} md={4}>
                  <Typography variant="subtitle2" color="textSecondary">Created At</Typography>
                  <Typography variant="body1">
                    {summary.CreatedAt ? dayjs(summary.CreatedAt).format('YYYY-MM-DD HH:mm:ss') : '-'}
                  </Typography>
                </Grid>
                <Grid item xs={12} sm={6} md={4}>
                  <Typography variant="subtitle2" color="textSecondary">Updated At</Typography>
                  <Typography variant="body1">
                    {summary.UpdatedAt ? dayjs(summary.UpdatedAt).format('YYYY-MM-DD HH:mm:ss') : '-'}
                  </Typography>
                </Grid>
                <Grid item xs={12} sm={6} md={4}>
                  <Typography variant="subtitle2" color="textSecondary">Status</Typography>
                  <FormControl fullWidth size="small" sx={{ mt: 1 }}>
                    <Select
                      value={summary.Status}
                      onChange={handleStatusChange}
                      disabled={statusLoading}
                    >
                      <MenuItem value={0}>Unchecked</MenuItem>
                      <MenuItem value={1}>Checked</MenuItem>
                      <MenuItem value={2}>Confirmed</MenuItem>
                    </Select>
                  </FormControl>
                </Grid>
              </Grid>
            </CardContent>
          </Card>
        </Grid>

        {renderBookingsTable(createdBookings, "Created Bookings")}
        {renderBookingsTable(completedBookings, "Completed Bookings")}
        {renderBookingsTable(canceledBookings, "Canceled Bookings")}
      </Grid>
    </Box>
  );
};

export default DailyBookingSummaryEdit;