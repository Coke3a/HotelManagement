import React, { useState, useEffect } from 'react';
import { Card, Space, Button, Tag, Table, Radio, message } from 'antd';
import { useParams, useNavigate } from 'react-router-dom';
import dayjs from 'dayjs';
import { handleTokenExpiration } from '../utils/api';

const EditDailySummary = () => {
  const { date } = useParams();
  const navigate = useNavigate();
  const [summary, setSummary] = useState(null);
  const [bookings, setBookings] = useState([]);
  const [loading, setLoading] = useState(false);
  const token = localStorage.getItem('token');

  const handleStatusChange = async (status) => {
    try {
      const response = await fetch(`http://localhost:8080/v1/daily-summary/status?date=${date}&status=${status}`, {
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

      message.success('Status updated successfully');
      fetchSummary();
    } catch (error) {
      console.error('Error updating status:', error);
      message.error('Failed to update status');
    }
  };

  const fetchSummary = async () => {
    setLoading(true);
    try {
      const response = await fetch(`http://localhost:8080/v1/daily-summary?date=${date}`, {
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
      setSummary(data.data);
      if (data.data.bookingIDs) {
        const bookingsResponse = await fetch(`http://localhost:8080/v1/bookings/ids?ids=${data.data.bookingIDs}`, {
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`,
          },
        });
        const bookingsData = await bookingsResponse.json();
        setBookings(bookingsData.data || []);
      }
    } catch (error) {
      console.error('Error fetching summary:', error);
      message.error('Failed to fetch summary');
    }
    setLoading(false);
  };

  const bookingColumns = [
    {
      title: 'Booking ID',
      dataIndex: 'id',
      key: 'id',
    },
    {
      title: 'Customer',
      dataIndex: ['customer', 'name'],
      key: 'customer',
    },
    {
      title: 'Room',
      dataIndex: ['room', 'number'],
      key: 'room',
    },
    {
      title: 'Amount',
      dataIndex: 'totalAmount',
      key: 'amount',
      render: (amount) => `$${amount.toFixed(2)}`,
    },
    {
      title: 'Status',
      dataIndex: 'status',
      key: 'status',
      render: (status) => <Tag color="blue">{status}</Tag>,
    },
  ];

  useEffect(() => {
    if (date) {
      fetchSummary();
    }
  }, [date]);

  if (!summary) {
    return null;
  }

  return (
    <div className="p-6">
      <div className="mb-4">
        <Button onClick={() => navigate('/daily-summary')}>Back</Button>
      </div>
      <Card title={`Summary for ${date}`} loading={loading}>
        <div className="grid grid-cols-2 gap-4 mb-6">
          <div>
            <p>Total Bookings: {summary.totalBookings}</p>
            <p>Total Amount: ${summary.totalAmount.toFixed(2)}</p>
          </div>
          <div>
            <p>Pending: {summary.pendingBookings}</p>
            <p>Confirmed: {summary.confirmedBookings}</p>
            <p>Checked In: {summary.checkedInBookings}</p>
            <p>Checked Out: {summary.checkedOutBookings}</p>
            <p>Canceled: {summary.canceledBookings}</p>
            <p>Completed: {summary.completedBookings}</p>
          </div>
        </div>
        <div className="mb-6">
          <h3 className="mb-2">Status</h3>
          <Radio.Group value={summary.status} onChange={(e) => handleStatusChange(e.target.value)}>
            <Radio.Button value={0}>Unchecked</Radio.Button>
            <Radio.Button value={1}>Checked</Radio.Button>
            <Radio.Button value={2}>Confirmed</Radio.Button>
          </Radio.Group>
        </div>
        <div>
          <h3 className="mb-2">Bookings</h3>
          <Table
            columns={bookingColumns}
            dataSource={bookings}
            rowKey="id"
            pagination={false}
          />
        </div>
      </Card>
    </div>
  );
};

export default EditDailySummary;