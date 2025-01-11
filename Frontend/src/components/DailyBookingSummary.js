import React, { useState, useEffect } from 'react';
import { Space, Button, DatePicker, message } from 'antd';
import { useNavigate } from 'react-router-dom';
import { EditOutlined } from '@ant-design/icons';
import dayjs from 'dayjs';
import { handleTokenExpiration } from '../utils/api';
import { countBookings } from '../utils/bookingUtils';
import { 
  Table, 
  TableBody, 
  TableCell, 
  TableContainer, 
  TableHead, 
  TableRow, 
  Paper,
  Typography,
  CircularProgress,
  Chip
} from '@mui/material';

const DailyBookingSummary = () => {
  const [summaries, setSummaries] = useState([]);
  const [loading, setLoading] = useState(false);
  const [selectedDate, setSelectedDate] = useState(dayjs());
  const navigate = useNavigate();
  const token = localStorage.getItem('token');

  const fetchSummaries = async (page = 1, limit = 10) => {
    setLoading(true);
    try {
      const response = await fetch(`http://localhost:8080/v1/daily-summary/list?page=${page}&limit=${limit}`, {
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
      });

      if (response.status === 401) {
        handleTokenExpiration(new Error("access token has expired"), navigate);
        return;
      }

      const data = await response.json();
      setSummaries(data.data || []);
    } catch (error) {
      console.error('Error fetching summaries:', error);
      message.error('Failed to fetch summaries');
      setSummaries([]);
    }
    setLoading(false);
  };

  const handleGenerateSummary = async (date) => {
    try {
      const response = await fetch(`http://localhost:8080/v1/daily-summary/generate?date=${date.format('YYYY-MM-DD')}`, {
        method: 'POST',
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
        throw new Error('Failed to generate summary');
      }

      message.success('Summary generated successfully');
      fetchSummaries();
    } catch (error) {
      console.error('Error generating summary:', error);
      message.error('Failed to generate summary');
    }
  };

  const getStatusColor = (status) => {
    switch (status) {
      case 0: // Unchecked
        return {
          backgroundColor: '#FFF4E5',
          textColor: '#663C00'
        };
      case 1: // Checked
        return {
          backgroundColor: '#E8F4FD',
          textColor: '#0A4FA6'
        };
      case 2: // Confirmed
        return {
          backgroundColor: '#EDF7ED',
          textColor: '#1E4620'
        };
      default:
        return {
          backgroundColor: '#F5F5F5',
          textColor: '#666666'
        };
    }
  };

  const renderStatus = (status) => {
    const statusText = {
      0: 'Unchecked',
      1: 'Checked',
      2: 'Confirmed'
    }[status] || 'Unknown';

    const { backgroundColor, textColor } = getStatusColor(status);

    return (
      <Chip 
        label={statusText}
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

  const handleDateChange = (date) => {
    if (date && dayjs(date).isValid()) {
      setSelectedDate(dayjs(date));
    } else {
      console.error('Invalid date selected');
    }
  };

  useEffect(() => {
    fetchSummaries();
  }, []);

  return (
    <div className="p-6">
      <div className="flex justify-between mb-4">
        <Space>
          <DatePicker
            value={selectedDate}
            onChange={handleDateChange}
            format="DD/MM/YYYY"
            disabledDate={(current) => current && current > dayjs().endOf('day')}
          />
          <Button type="primary" onClick={() => handleGenerateSummary(selectedDate)}>
            Generate Summary for Selected Date
          </Button>
        </Space>
      </div>
      
      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow sx={{ backgroundColor: '#f5f5f5' }}>
              <TableCell>Date</TableCell>
              <TableCell>Created Bookings</TableCell>
              <TableCell>Completed Bookings</TableCell>
              <TableCell>Canceled Bookings</TableCell>
              <TableCell>Total Amount</TableCell>
              <TableCell>Status</TableCell>
              <TableCell>Created At</TableCell>
              <TableCell>Action</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {loading ? (
              <TableRow>
                <TableCell colSpan={8} align="center">
                  <CircularProgress size={24} />
                </TableCell>
              </TableRow>
            ) : (!summaries || summaries.length === 0) ? (
              <TableRow>
                <TableCell colSpan={8} align="center">
                  <Typography variant="body2">No summaries found</Typography>
                </TableCell>
              </TableRow>
            ) : (
              summaries.map((summary, index) => (
                <TableRow key={summary.SummaryDate} className={index % 2 === 0 ? "bg-gray-50" : "bg-white"}>
                  <TableCell>{dayjs(summary.SummaryDate).format('DD/MM/YYYY')}</TableCell>
                  <TableCell>{countBookings(summary.CreatedBookings)}</TableCell>
                  <TableCell>{countBookings(summary.CompletedBookings)}</TableCell>
                  <TableCell>{countBookings(summary.CanceledBookings)}</TableCell>
                  <TableCell>{summary.TotalAmount}</TableCell>
                  <TableCell>{renderStatus(summary.Status)}</TableCell>
                  <TableCell>{dayjs(summary.CreatedAt).format('DD/MM/YYYY HH:mm:ss')}</TableCell>
                  <TableCell>
                    <Button
                      variant="outlined"
                      size="small"
                      onClick={() => navigate(`/daily-summary/edit/${dayjs(summary.SummaryDate).format('YYYY-MM-DD')}`)}
                      icon={<EditOutlined />}
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
};

export default DailyBookingSummary;