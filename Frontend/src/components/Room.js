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
  Button,
  Box,
  CircularProgress,
  Typography,
  Stack,
  Pagination,
} from '@mui/material';
import { getUserRole } from '../utils/auth';
import { UserRoleEnum } from '../utils/userRoleEnum';

const Room = () => {
  const navigate = useNavigate();
  const [rooms, setRooms] = useState([]);
  const [loading, setLoading] = useState(true);
  const [userRole, setUserRole] = useState('');
  const [roomTypeMap, setRoomTypeMap] = useState({});
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [totalCount, setTotalCount] = useState(0);
  const [error, setError] = useState(null);

  useEffect(() => {
    setUserRole(getUserRole());
  }, []);

  const fetchRooms = async (page, rowsPerPage) => {
    setLoading(true);
    try {
      const token = localStorage.getItem('token');
      const skip = page * rowsPerPage;
      const queryParams = new URLSearchParams({
        skip: skip.toString(),
        limit: rowsPerPage.toString(),
      });
      const response = await fetch(`http://localhost:8080/v1/rooms/with-room-type?${queryParams.toString()}`, {
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
      if (data && data.data && data.data.rooms) {
        setRooms(data.data.rooms);
      } else {
        setRooms([]);
      }
      if (data && data.data && data.data.meta && typeof data.data.meta.total === 'number') {
        setTotalCount(data.data.meta.total);
      } else {
        setTotalCount(data.data?.rooms?.length || 0);
      }
    } catch (error) {
      console.error('Error fetching rooms:', error);
      setError(error.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchRooms(page, rowsPerPage);
  }, [page, rowsPerPage, navigate]);

  const handlePageChange = (event, newPage) => {
    setPage(newPage - 1);
  };

  const handleEditRoom = (id) => {
    navigate(`/room/edit/${id}`);
  };

  const handleAddRoom = () => {
    navigate('/room/add');
  };

  const getStatusMessage = (status) => {
    switch (status) {
      case 1:
        return 'Available';
      case 2:
        return 'Maintenance';
      default:
        return 'Unknown';
    }
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
        {userRole === UserRoleEnum.ADMIN && (
          <Button
            variant="contained"
            className="bg-blue-600 text-white hover:bg-blue-700"
            onClick={handleAddRoom}
          >
            Add Room
          </Button>
        )}
        
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
              <TableCell>Room Number</TableCell>
              <TableCell>Room Type</TableCell>
              <TableCell>Status</TableCell>
              <TableCell>Floor</TableCell>
              <TableCell>Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {loading ? (
              <TableRow>
                <TableCell colSpan={6} align="center">
                  <CircularProgress size={24} />
                </TableCell>
              </TableRow>
            ) : !rooms || rooms.length === 0 ? (
              <TableRow>
                <TableCell colSpan={6} align="center">
                  <Typography variant="body2">No rooms found</Typography>
                </TableCell>
              </TableRow>
            ) : (
              rooms.map((room) => (
                <TableRow key={room.id}>
                  <TableCell>{room.id}</TableCell>
                  <TableCell>{room.room_number}</TableCell>
                  <TableCell>{room.room_type_name || 'Unknown'}</TableCell>
                  <TableCell>{getStatusMessage(room.status)}</TableCell>
                  <TableCell>{room.floor}</TableCell>
                  <TableCell>
                    {userRole === UserRoleEnum.ADMIN && (
                      <Button
                        variant="outlined"
                        color="primary"
                        size="small"
                        onClick={() => handleEditRoom(room.id)}
                      >
                        Detail
                      </Button>
                    )}
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

export default Room;
