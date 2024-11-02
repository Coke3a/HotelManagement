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
  Box,
  CircularProgress,
  Typography,
  Button,
} from '@mui/material';
import { handleTokenExpiration } from '../utils/auth';

const Log = () => {
  const token = localStorage.getItem('token');
  const navigate = useNavigate();
  const [logs, setLogs] = useState([]);
  const [users, setUsers] = useState({});
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchUsers = async () => {
      try {
        const response = await fetch('http://localhost:8080/v1/users/?skip=0&limit=100', {
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
          throw new Error('Failed to fetch users');
        }

        const data = await response.json();
        const usersMap = {};
        data.data.users.forEach(user => {
          usersMap[user.id] = user.username;
        });
        setUsers(usersMap);
      } catch (error) {
        console.error('Error fetching users:', error);
      }
    };

    const fetchLogs = async () => {
      try {
        const queryParams = new URLSearchParams({
          skip: '0',
          limit: '100',
        });

        const response = await fetch(`http://localhost:8080/v1/logs/?${queryParams.toString()}`, {
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
          throw new Error('Failed to fetch logs');
        }

        const data = await response.json();
        setLogs(data.data.logs || []);
      } catch (error) {
        console.error('Error fetching logs:', error);
        setLogs([]);
      } finally {
        setLoading(false);
      }
    };

    fetchUsers();
    fetchLogs();
  }, []);

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleString();
  };

  const handleRecordClick = (tableName, recordId) => {
    switch (tableName) {
      case 'users':
        navigate(`/user/edit/${recordId}`);
        break;
      case 'rooms':
        navigate(`/room/edit/${recordId}`);
        break;
      case 'customers':
        navigate(`/guest/edit/${recordId}`);
        break;
      case 'bookings':
        navigate(`/booking/edit/${recordId}`);
        break;
      default:
        console.log('No redirect configured for this table type');
    }
  };

  const handleUserClick = (userId) => {
    navigate(`/user/edit/${userId}`);
  };

  return (
    <div className="table-container">
      <Typography variant="h5" gutterBottom>
        System Logs
      </Typography>
      <TableContainer component={Paper}>
        <Table size="small" className="compact-table">
          <TableHead>
            <TableRow>
              <TableCell>ID</TableCell>
              <TableCell>Table Name</TableCell>
              <TableCell>Record</TableCell>
              <TableCell>Action</TableCell>
              <TableCell>User</TableCell>
              <TableCell>Created At</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {loading ? (
              <TableRow>
                <TableCell colSpan={6} align="center">
                  <Box display="flex" justifyContent="center" alignItems="center" height="50px">
                    <CircularProgress />
                  </Box>
                </TableCell>
              </TableRow>
            ) : !logs || logs.length === 0 ? (
              <TableRow>
                <TableCell colSpan={6} align="center">
                  <Typography variant="body1">No logs found</Typography>
                </TableCell>
              </TableRow>
            ) : (
              logs.map((log) => (
                <TableRow key={log.id}>
                  <TableCell>{log.id}</TableCell>
                  <TableCell>{log.table_name}</TableCell>
                  <TableCell>
                    <Button
                      color="primary"
                      onClick={() => handleRecordClick(log.table_name, log.record_id)}
                    >
                      {log.record_id}
                    </Button>
                  </TableCell>
                  <TableCell>{log.action}</TableCell>
                  <TableCell>
                    <Button
                      color="primary"
                      onClick={() => handleUserClick(log.user_id)}
                    >
                      {users[log.user_id] || log.user_id}
                    </Button>
                  </TableCell>
                  <TableCell>{formatDate(log.created_at)}</TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </TableContainer>
    </div>
  );
};

export default Log;
