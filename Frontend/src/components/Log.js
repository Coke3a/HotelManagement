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
  Stack,
  Pagination,
} from '@mui/material';
import { handleTokenExpiration } from '../utils/api';

const Log = () => {
  const navigate = useNavigate();
  const [logs, setLogs] = useState([]);
  const [loading, setLoading] = useState(true);
  const [users, setUsers] = useState({});
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [totalCount, setTotalCount] = useState(0);
  const [error, setError] = useState(null);

  const fetchLogs = async (page, rowsPerPage) => {
    setLoading(true);
    try {
      const skip = page * rowsPerPage;
      const queryParams = new URLSearchParams({
        skip: skip.toString(),
        limit: rowsPerPage.toString(),
      });
      const token = localStorage.getItem('token');
      const response = await fetch(`http://localhost:8080/v1/logs/?${queryParams.toString()}`, {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
      });

      if (response.status === 401) {
        handleTokenExpiration(new Error("access token has expired"), navigate);
        return;
      }

      const data = await response.json();
      if (data && data.data) {
        setLogs(data.data.logs || []);
        if (data.data.meta && typeof data.data.meta.total === 'number') {
          setTotalCount(data.data.meta.total);
        } else {
          setTotalCount(data.data?.logs?.length || 0);
        }
      } else {
        setLogs([]);
        setTotalCount(0);
      }
    } catch (error) {
      console.error('Error fetching logs:', error);
      setError(error.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchLogs(page, rowsPerPage);
  }, [page, rowsPerPage, navigate]);

  const handlePageChange = (event, newPage) => {
    setPage(newPage - 1);
  };

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
      case 'payments':
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
      <Box sx={{ 
        display: 'flex', 
        justifyContent: 'space-between', 
        alignItems: 'center',
        mb: 2,
        mt: 0,
        pt: 0
      }}>
        {/* Right side - Pagination controls */}
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

      {error && (
        <Typography color="error" align="center" style={{ padding: '1rem' }}>
          {error}
        </Typography>
      )}
    </div>
  );
};

export default Log;
