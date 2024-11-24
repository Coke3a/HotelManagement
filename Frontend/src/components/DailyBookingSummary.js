import React, { useState, useEffect } from 'react';
import { Table, Space, Button, Tag, DatePicker, message } from 'antd';
import { useNavigate } from 'react-router-dom';
import { EditOutlined } from '@ant-design/icons';
import dayjs from 'dayjs';
import { handleTokenExpiration } from '../utils/api';

const DailyBookingSummary = () => {
  const [summaries, setSummaries] = useState([]);
  const [loading, setLoading] = useState(false);
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  });
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
      setSummaries(data.data);
      setPagination({
        ...pagination,
        total: data.total,
        current: page,
        pageSize: limit,
      });
    } catch (error) {
      console.error('Error fetching summaries:', error);
      message.error('Failed to fetch summaries');
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

  const getStatusTag = (status) => {
    const statusMap = {
      0: { color: 'default', text: 'Unchecked' },
      1: { color: 'processing', text: 'Checked' },
      2: { color: 'success', text: 'Confirmed' },
    };
    
    const statusInfo = statusMap[status] || { color: 'default', text: 'Unknown' };
    return <Tag color={statusInfo.color}>{statusInfo.text}</Tag>;
  };

  const columns = [
    {
      title: 'Date',
      dataIndex: 'summaryDate',
      key: 'summaryDate',
      render: (date) => dayjs(date).format('YYYY-MM-DD'),
    },
    {
      title: 'Total Bookings',
      dataIndex: 'totalBookings',
      key: 'totalBookings',
    },
    {
      title: 'Total Amount',
      dataIndex: 'totalAmount',
      key: 'totalAmount',
      render: (amount) => amount ? `${amount.toFixed(2)}` : '0.00',
    },
    {
      title: 'Status',
      dataIndex: 'status',
      key: 'status',
      render: getStatusTag,
    },
    {
      title: 'Action',
      key: 'action',
      render: (_, record) => (
        <Space>
          <Button
            type="primary"
            icon={<EditOutlined />}
            onClick={() => navigate(`/daily-summary/${dayjs(record.summaryDate).format('YYYY-MM-DD')}`)}
          >
            Edit
          </Button>
        </Space>
      ),
    },
  ];

  useEffect(() => {
    fetchSummaries();
  }, []);

  return (
    <div className="p-6">
      <div className="flex justify-between mb-4">
        <h1 className="text-2xl font-bold">Daily Booking Summaries</h1>
        <Space>
          <DatePicker onChange={handleGenerateSummary} />
          <Button type="primary" onClick={() => handleGenerateSummary(dayjs())}>
            Generate Today's Summary
          </Button>
        </Space>
      </div>
      <Table
        columns={columns}
        dataSource={summaries}
        rowKey="summaryDate"
        pagination={pagination}
        loading={loading}
        onChange={(pagination) => fetchSummaries(pagination.current, pagination.pageSize)}
      />
    </div>
  );
};

export default DailyBookingSummary;