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
} from '@mui/material';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import { BookingStatus, getBookingStatusMessage } from '../utils/bookingStatusEnums';

const DailyBookingSummaryEdit = () => {
  const { date } = useParams();
  const navigate = useNavigate();
  const [summary, setSummary] = useState(null);
  const [bookings, setBookings] = useState([]);
  const [loading, setLoading] = useState(false);
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
      return Array.isArray(data.data) ? data.data : [];
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

      // Fetch bookings based on date and status
      const allBookings = [];

      // Fetch created bookings
      const createdBookings = await fetchBookings({ created_at: date }, 0, 100);
      allBookings.push(...createdBookings);

      // Fetch completed bookings
      const completedBookings = await fetchBookings({ updated_at: date, booking_status: BookingStatus.COMPLETED }, 0, 100);
      allBookings.push(...completedBookings);

      // Fetch canceled bookings
      const canceledBookings = await fetchBookings({ updated_at: date, booking_status: BookingStatus.CANCELED }, 0, 100);
      allBookings.push(...canceledBookings);

      setBookings(allBookings);
    } catch (error) {
      console.error('Error fetching summary:', error);
      message.error('Failed to fetch summary');
    }
    setLoading(false);
  };

  const handleStatusChange = async (event, newStatus) => {
    if (newStatus !== null) {
      try {
        const response = await fetch(`http://localhost:8080/v1/daily-summary/status?date=${date}&status=${newStatus}`, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`,
          },
        });

        if (!response.ok) {
          throw new Error('Failed to update status');
        }

        message.success('Status updated successfully');
        fetchSummary();
      } catch (error) {
        console.error('Error updating status:', error);
        message.error('Failed to update status');
      }
    }
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
                    {summary.TotalAmount.toFixed(2)}
                  </Typography>
                </Grid>
                <Grid item xs={12} sm={6} md={4}>
                  <Typography variant="subtitle2" color="textSecondary">Created At</Typography>
                  <Typography variant="body1">
                    {summary.created_at ? dayjs(summary.created_at).format('YYYY-MM-DD HH:mm:ss') : '-'}
                  </Typography>
                </Grid>
                <Grid item xs={12} sm={6} md={4}>
                  <Typography variant="subtitle2" color="textSecondary">Updated At</Typography>
                  <Typography variant="body1">
                    {summary.updated_at ? dayjs(summary.updated_at).format('YYYY-MM-DD HH:mm:ss') : '-'}
                  </Typography>
                </Grid>
              </Grid>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>Bookings</Typography>
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
                    {bookings.map((booking, index) => (
                      <TableRow key={booking.id} className={index % 2 === 0 ? "bg-gray-50" : "bg-white"}>
                        <TableCell>{booking.id}</TableCell>
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
                        <TableCell>{booking.total_amount || booking.booking_price}</TableCell>
                        <TableCell>{renderBookingStatus(booking.booking_status)}</TableCell>
                        <TableCell>
                          <Button
                            variant="outlined"
                            size="small"
                            onClick={() => navigate(`/booking/edit/${booking.id}`)}
                          >
                            Detail
                          </Button>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </TableContainer>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Box>
  );
};

export default DailyBookingSummaryEdit;