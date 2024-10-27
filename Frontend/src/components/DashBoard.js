import React, { useState, useEffect } from 'react';
import { Grid, Paper, Typography, Table, TableBody, TableCell, TableContainer, TableHead, TableRow } from '@mui/material';
import CircularProgress from '@mui/material/CircularProgress';
import { PaymentStatus } from '../utils/paymentEnums';
import { BookingStatus, getBookingStatusMessage } from '../utils/bookingStatusEnums'; // Import the enum and function
import { getPaymentStatusMessage } from '../utils/paymentEnums'; // Import the function
import { Button } from '@mui/material';
import { useNavigate } from 'react-router-dom';

const DashBoard = () => {
  const token = localStorage.getItem('token');
  const [rooms, setRooms] = useState([]);
  const [bookings, setBookings] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);
      setError(null);
      try {
        const dates = generateDates();
        const startDate = dates[0].date;
        const endDate = dates[dates.length - 1].date;

        const [bookingsResponse, roomsResponse] = await Promise.all([
          fetch(`http://localhost:8080/v1/booking/?${new URLSearchParams({
            skip: '0',
            limit: '100',
            check_in_date: startDate,
            check_out_date: endDate,
          })}`, {
            headers: {
              'Content-Type': 'application/json',
              'Authorization': `Bearer ${token}`,
            },
          }),
          fetch('http://localhost:8080/v1/rooms/with-room-type?skip=0&limit=100', {
            headers: {
              'Content-Type': 'application/json',
              'Authorization': `Bearer ${token}`,
            },
          })
        ]);

        if (bookingsResponse.status === 401) {
          handleTokenExpiration(new Error("access token has expired"), navigate);
          return;
        }

        if (!bookingsResponse.ok || !roomsResponse.ok) {
          throw new Error(`Failed to fetch data: ${bookingsResponse.status} ${bookingsResponse.statusText}, ${roomsResponse.status} ${roomsResponse.statusText}`);
        }

        const [bookingsData, roomsData] = await Promise.all([
          bookingsResponse.json(),
          roomsResponse.json()
        ]);

        const bookingPayments = bookingsData.data.booking_customer_payments || [];
        setBookings(bookingPayments);

        const roomsWithRoomType = roomsData.data || [];
        setRooms(roomsWithRoomType);
      } catch (error) {
        console.error('Error fetching data:', error);
        setError(error.message);
        setRooms([]);
        setBookings([]);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  const renderGuestName = (roomId, date) => {
    const booking = bookings.find(
      (booking) =>
        booking.room_id === roomId &&
        new Date(booking.check_in_date) <= new Date(date) &&
        new Date(booking.check_out_date) > new Date(date)
    );

    if (booking) {
      const checkInDate = new Date(booking.check_in_date);
      const checkOutDate = new Date(booking.check_out_date);
      const currentDate = new Date(date);
      const stayDuration = (checkOutDate - checkInDate) / (1000 * 60 * 60 * 24);
      const middleDayIndex = Math.floor(stayDuration / 2);
      const middleDay = new Date(checkInDate);
      middleDay.setDate(checkInDate.getDate() + middleDayIndex);

      const isMiddleDay = middleDay.toDateString() === currentDate.toDateString();

      const getStatusColor = (booking) => {
        const currentDate = new Date().toISOString().slice(0, 10);
        if (booking.check_out_date <= currentDate) {
          return '#c8e6c9'; // checked-out
        } else if (booking.check_in_date <= currentDate) {
          return '#ffecb3'; // in-progress
        } else {
          return '#e3f2fd'; // upcoming
        }
      };

      const getStatusText = (booking) => {
        return getBookingStatusMessage(booking.booking_status);
      };

      const getPaymentStatusText = (status) => {
        return getPaymentStatusMessage(status);
      };

      return (
        <div
          style={{
            backgroundColor: getStatusColor(booking),
            borderRadius: '4px',
            cursor: 'pointer',
            position: 'relative',
            height: '100%',
            width: '100%',
            display: 'flex',
            flexDirection: 'column',
            justifyContent: 'center',
            alignItems: 'center',
            padding: '4px',
            boxShadow: '0 2px 4px rgba(0, 0, 0, 0.1)',
          }}
          onClick={() => window.location.href = `/booking/edit/${booking.booking_id}`}
          title={`${booking.customer_firstname} ${booking.customer_surname} (${booking.check_in_date} - ${booking.check_out_date})`}
        >
          {isMiddleDay && (
            <div style={{ textAlign: 'center' }}>
              <div style={{ fontSize: '14px', fontWeight: 'bold', marginBottom: '2px' }}>
                {`${booking.customer_firstname}. ${booking.customer_surname.charAt(0)}`}
              </div>
              <div style={{ fontSize: '12px', color: '#555' }}>
                {getStatusText(booking)}
              </div>
              <div style={{ fontSize: '12px', color: '#777' }}>
                {getPaymentStatusText(booking.payment_status)}
              </div>
            </div>
          )}
        </div>
      );
    }

    return null;
  };

  const generateDates = () => {
    const today = new Date();
    const dates = [];

    const startDate = new Date(today);
    startDate.setDate(today.getDate() - 2);

    for (let i = 0; i < 13; i++) {
      const date = new Date(startDate);
      date.setDate(startDate.getDate() + i);
      dates.push({
        date: date.toLocaleString('en-US', { timeZone: 'Asia/Bangkok' }).split(',')[0],
        dayOfWeek: date.toLocaleString('en-US', { timeZone: 'Asia/Bangkok', weekday: 'long' }),
      });
    }

    return dates;
  };

  const dates = generateDates();

  return (
    <Grid container spacing={2} style={{ padding: '10px' }}>
      <Grid item xs={12} style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Typography variant="h5" gutterBottom>
          Dashboard
        </Typography>
        <Button
          variant="contained"
          color="secondary"
          onClick={() => navigate('/receipt/latest')}
        >
          Print Latest Receipt
        </Button>
      </Grid>
      <Grid item xs={12}>
        <TableContainer component={Paper} style={{ maxHeight: '70vh', overflow: 'auto' }}>
          <Table size="small" stickyHeader>
            <TableHead>
              <TableRow>
                <TableCell style={{ padding: '4px', fontSize: '0.8rem' }}>ROOM</TableCell>
                {dates.map((dateObj) => (
                  <TableCell key={dateObj.date} align="center" style={{ padding: '4px', fontSize: '0.8rem' }}>
                    {dateObj.date}<br/>({dateObj.dayOfWeek})
                  </TableCell>
                ))}
              </TableRow>
            </TableHead>
            <TableBody>
              {loading ? (
                <TableRow>
                  <TableCell colSpan={dates.length + 1} align="center">
                    <CircularProgress size={20} />
                  </TableCell>
                </TableRow>
              ) : error ? (
                <TableRow>
                  <TableCell colSpan={dates.length + 1} align="center">
                    <Typography color="error" style={{ fontSize: '0.8rem' }}>{error}</Typography>
                  </TableCell>
                </TableRow>
              ) : rooms.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={dates.length + 1} align="center">
                    <Typography style={{ fontSize: '0.8rem' }}>No rooms available</Typography>
                  </TableCell>
                </TableRow>
              ) : (
                rooms.map((room) => (
                  <TableRow key={room.id}>
                    <TableCell style={{ padding: '4px', fontSize: '0.8rem' }}>
                      {room.room_number}
                      <br />
                      <Typography variant="caption" color="textSecondary">
                        {room.room_type_name}
                      </Typography>
                    </TableCell>
                    {dates.map((dateObj) => (
                      <TableCell 
                        key={`${room.id}-${dateObj.date}`} 
                        align="center" 
                        style={{ 
                          padding: '2px', 
                          height: '50px',
                          position: 'relative',
                          fontSize: '0.8rem'
                        }}
                      >
                        {renderGuestName(room.id, dateObj.date)}
                      </TableCell>
                    ))}
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </TableContainer>
      </Grid>
    </Grid>
  );
};

export default DashBoard;
