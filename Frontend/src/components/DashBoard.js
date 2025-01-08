import React, { useState, useEffect } from 'react';
import { Typography, CircularProgress, Button, Card, CardContent, Grid, Chip } from '@mui/material';
import { PaymentStatus, getPaymentStatusMessage, getPaymentStatusColor } from '../utils/paymentEnums';
import { BookingStatus, getBookingStatusMessage, getBookingStatusColor } from '../utils/bookingStatusEnums';
import { useNavigate } from 'react-router-dom';
import { handleTokenExpiration } from '../utils/api';

const DashBoard = () => {
  const token = localStorage.getItem('token');
  const [rooms, setRooms] = useState([]);
  const [bookings, setBookings] = useState([]);
  const [todayCheckIns, setTodayCheckIns] = useState([]);
  const [todayCheckOuts, setTodayCheckOuts] = useState([]);
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
        const today = new Date().toISOString().slice(0, 10);

        const [bookingsResponse, roomsResponse, checkInsResponse, checkOutsResponse] = await Promise.all([
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
          }),
          fetch(`http://localhost:8080/v1/booking/?${new URLSearchParams({
            skip: '0',
            limit: '100',
            check_in_date: today,
          })}`, {
            headers: {
              'Content-Type': 'application/json',
              'Authorization': `Bearer ${token}`,
            },
          }),
          fetch(`http://localhost:8080/v1/booking/?${new URLSearchParams({
            skip: '0',
            limit: '100',
            check_out_date: today,
          })}`, {
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

        const [bookingsData, roomsData, checkInsData, checkOutsData] = await Promise.all([
          bookingsResponse.json(),
          roomsResponse.json(),
          checkInsResponse.json(),
          checkOutsResponse.json()
        ]);

        const bookingPayments = bookingsData.data.booking_customer_payments || [];
        setBookings(bookingPayments);

        const checkInPayments = checkInsData.data.booking_customer_payments || [];
        setTodayCheckIns(checkInPayments);

        const checkOutPayments = checkOutsData.data.booking_customer_payments || [];
        setTodayCheckOuts(checkOutPayments);

        const roomsWithRoomType = roomsData.data?.rooms || [];
        setRooms(roomsWithRoomType);
      } catch (error) {
        console.error('Error fetching data:', error);
        setError(error.message);
        setRooms([]);
        setBookings([]);
        setTodayCheckIns([]);
        setTodayCheckOuts([]);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  const renderBookingRow = (room) => {
    const datesList = generateDates();
    let skipCells = 0;

    return datesList.map((dateObj, index) => {
      if (skipCells > 0) {
        skipCells--;
        return null;
      }

      const booking = bookings.find(
        (booking) => {
          const currentDate = new Date(dateObj.date);
          const checkInDate = new Date(booking.check_in_date);
          const checkOutDate = new Date(booking.check_out_date);
          
          // Reset hours to compare dates properly
          currentDate.setHours(0, 0, 0, 0);
          checkInDate.setHours(0, 0, 0, 0);
          checkOutDate.setHours(0, 0, 0, 0);
          
          return booking.room_id === room.id &&
                 checkInDate <= currentDate &&
                 checkOutDate > currentDate;
        }
      );

      if (booking) {
        const startDate = new Date(booking.check_in_date);
        const endDate = new Date(booking.check_out_date);
        startDate.setHours(0, 0, 0, 0);
        endDate.setHours(0, 0, 0, 0);

        const startIndex = datesList.findIndex(d => {
          const date = new Date(d.date);
          date.setHours(0, 0, 0, 0);
          return date.getTime() === startDate.getTime();
        });

        const dayBeforeCheckout = new Date(endDate);
        dayBeforeCheckout.setDate(endDate.getDate() - 1);

        const endIndex = datesList.findIndex(d => {
          const date = new Date(d.date);
          date.setHours(0, 0, 0, 0);
          return date.getTime() === dayBeforeCheckout.getTime();
        });

        if (startIndex !== -1 && endIndex !== -1) {
          const colspan = endIndex - startIndex + 1;
          skipCells = colspan - 1;

          return (
            <td 
              key={`${room.id}-${dateObj.date}`}
              colSpan={colspan}
              className="p-0.5 h-10 border-r border-gray-200"
            >
              <div
                className="h-full w-full transition-all duration-200 hover:shadow-md cursor-pointer flex flex-col items-center justify-center p-1 space-y-1 rounded booking-hover-effect"
                style={{ 
                  backgroundColor: getStatusColor(booking),
                  transition: 'all 0.2s ease-in-out',
                  minHeight: '80px',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center'
                }}
                onClick={() => navigate(`/booking/edit/${booking.booking_id}`)}
                title={`${booking.customer_firstname} ${booking.customer_surname} (${booking.check_in_date} - ${booking.check_out_date})`}
                onMouseEnter={(e) => {
                  e.currentTarget.style.transform = 'translateY(-2px)';
                  e.currentTarget.style.boxShadow = '0 4px 6px -1px rgba(0, 0, 0, 0.1)';
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.transform = 'translateY(0)';
                  e.currentTarget.style.boxShadow = 'none';
                }}
              >
                <div className="flex flex-col items-center justify-center space-y-1 w-full">
                  <div className="font-semibold text-blue-900 text-[11px] text-center">
                    {`${booking.customer_firstname}. ${booking.customer_surname.charAt(0)}`}
                  </div>
                  {renderBookingStatusBadge(booking.booking_status)}
                  {renderPaymentStatusBadge(booking.payment_status)}
                </div>
              </div>
            </td>
          );
        }
      }

      return (
        <td 
          key={`${room.id}-${dateObj.date}`}
          className="p-1 h-16 border-r border-gray-200"
        />
      );
    }).filter(cell => cell !== null);
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

  const getStatusColor = (booking) => {
    const currentDate = new Date().toISOString().slice(0, 10);
    if (booking.check_out_date <= currentDate) {
      return '#EDF7ED'; // Past stays - light green
    } else if (booking.check_in_date <= currentDate) {
      return '#FFF4E5'; // Current stays - light orange
    } else {
      return '#E8F4FD'; // Future stays - light blue
    }
  };

  const renderBookingStatusBadge = (status) => {
    const statusMessage = getBookingStatusMessage(status);
    const { chipColor } = getBookingStatusColor(status);

    return (
      <div className="flex items-center space-x-0.5 text-[10px] bg-white/80 px-1.5 py-0.5 rounded">
        <Chip 
          label={statusMessage} 
          color={chipColor}
          size="small"
          sx={{ height: '20px', '& .MuiChip-label': { fontSize: '10px', px: 1 } }}
        />
      </div>
    );
  };

  const renderPaymentStatusBadge = (status) => {
    const statusMessage = getPaymentStatusMessage(status);
    const { backgroundColor, textColor } = getPaymentStatusColor(status);

    return (
      <div className="flex items-center space-x-0.5 text-[10px] bg-white/80 px-1.5 py-0.5 rounded">
        <Chip 
          label={statusMessage}
          sx={{
            backgroundColor: backgroundColor,
            color: textColor,
            height: '20px',
            '& .MuiChip-label': { 
              fontSize: '10px',
              px: 1,
              fontWeight: 'medium'
            }
          }}
          size="small"
        />
      </div>
    );
  };

  // Calculate today's check-ins, check-outs, and remaining bookings
  const today = new Date().toISOString().slice(0, 10);
  const checkInTodayCount = todayCheckIns.length;
  const checkOutTodayCount = todayCheckOuts.length;
  const remainingCheckInCount = todayCheckIns.filter(booking => 
    booking.booking_status !== BookingStatus.CHECKED_IN
  ).length;
  const remainingCheckOutCount = todayCheckOuts.filter(booking => 
    booking.booking_status !== BookingStatus.CHECKED_OUT
  ).length;

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <Typography variant="h5" className="text-indigo-900 font-semibold">
          Dashboard
        </Typography>
      </div>

      <Grid container spacing={3} className="mb-4">
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography variant="subtitle2" color="textSecondary">
                Check-ins Today
              </Typography>
              <Typography variant="h6" className="text-indigo-900">
                {checkInTodayCount}
              </Typography>
              <Typography variant="subtitle2" color="textSecondary" className="mt-2">
                Remaining Check-ins
              </Typography>
              <Typography variant="h6" className="text-indigo-900">
                {remainingCheckInCount}
              </Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography variant="subtitle2" color="textSecondary">
                Check-outs Today
              </Typography>
              <Typography variant="h6" className="text-indigo-900">
                {checkOutTodayCount}
              </Typography>
              <Typography variant="subtitle2" color="textSecondary" className="mt-2">
                Remaining Check-outs
              </Typography>
              <Typography variant="h6" className="text-indigo-900">
                {remainingCheckOutCount}
              </Typography>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      <div className="overflow-auto rounded-lg shadow">
        <table className="w-full border-collapse table-auto dashboard-table">
          <thead>
            <tr>
              <th className="sticky left-0 z-10 bg-indigo-900 p-3 text-sm font-semibold tracking-wide text-left text-white border-b border-r border-indigo-800">
                ROOM
              </th>
              {generateDates().map((dateObj) => {
                const isToday = dateObj.date === new Date().toLocaleString('en-US', { timeZone: 'Asia/Bangkok' }).split(',')[0];
                return (
                  <th 
                    key={dateObj.date}
                    className={`bg-indigo-900 p-1.5 text-[11px] font-semibold tracking-wide text-center text-white border-b border-r border-indigo-800 min-w-[80px] ${isToday ? 'relative bg-indigo-700' : ''}`}
                  >
                    {isToday && (
                      <div className="absolute top-0 left-0 right-0 bg-yellow-500 text-indigo-900 text-[10px] font-bold px-1 py-0.5">
                        TODAY
                      </div>
                    )}
                    <div className="font-bold mt-3">{dateObj.dayOfWeek}</div>
                    <div className={`text-[10px] mt-0.5 ${isToday ? 'text-yellow-300 font-bold' : 'text-indigo-200'}`}>
                      {dateObj.date}
                    </div>
                  </th>
                );
              })}
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {loading ? (
              <tr>
                <td colSpan={generateDates().length + 1} className="text-center p-4">
                  <CircularProgress size={24} />
                </td>
              </tr>
            ) : error ? (
              <tr>
                <td colSpan={generateDates().length + 1} className="text-center p-4 text-red-500">
                  {error}
                </td>
              </tr>
            ) : rooms.length === 0 ? (
              <tr>
                <td colSpan={generateDates().length + 1} className="text-center p-4">
                  No rooms available
                </td>
              </tr>
            ) : (
              rooms.map((room) => (
                <tr key={room.id} className="hover:bg-gray-50">
                  <td className="room-cell p-1.5 text-[11px] border-r border-gray-200">
                    <div className="font-semibold text-indigo-900">{room.room_number}</div>
                    <div className="text-gray-600 text-[10px]">{room.room_type_name}</div>
                  </td>
                  {renderBookingRow(room)}
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default DashBoard;
