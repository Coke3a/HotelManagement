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
import { handleTokenExpiration } from '../utils/api';

const RoomType = () => {
  const navigate = useNavigate();
  const [roomTypes, setRoomTypes] = useState([]);
  const [loading, setLoading] = useState(true);
  const [userRole, setUserRole] = useState('');
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [totalCount, setTotalCount] = useState(0);
  const [error, setError] = useState(null);

  useEffect(() => {
    setUserRole(getUserRole());
  }, []);

  const fetchRoomTypes = async (page, rowsPerPage) => {
    setLoading(true);
    try {
      const skip = page * rowsPerPage;
      const queryParams = new URLSearchParams({
        skip: skip.toString(),
        limit: rowsPerPage.toString(),
      });
      const token = localStorage.getItem('token');
      const response = await fetch(`http://localhost:8080/v1/room-types/?${queryParams.toString()}`, {
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
      if (data && data.data && data.data.roomTypes) {
        setRoomTypes(data.data.roomTypes || []);
      } else {
        setRoomTypes([]);
      }
      if (data && data.data && data.data.meta && typeof data.data.meta.total === 'number') {
        setTotalCount(data.data.meta.total);
      } else {
        setTotalCount(data.data?.room_types?.length || 0);
      }
    } catch (error) {
      console.error('Error fetching room types:', error);
      setError(error.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchRoomTypes(page, rowsPerPage);
  }, [page, rowsPerPage, navigate]);

  const handlePageChange = (event, newPage) => {
    setPage(newPage - 1);
  };

  const handleEditRoomType = (id) => {
    navigate(`/room-type/edit/${id}`);
  };

  const handleAddRoomType = () => {
    navigate('/room-type/add');
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
            sx={{ mb: 1 }}
            onClick={handleAddRoomType}
          >
            Add Room Type
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
              <TableCell>Name</TableCell>
              <TableCell>Description</TableCell>
              <TableCell>Capacity</TableCell>
              <TableCell>Default Price</TableCell>
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
            ) : !roomTypes || roomTypes.length === 0 ? (
              <TableRow>
                <TableCell colSpan={6} align="center">
                  <Typography variant="body2">No room types found</Typography>
                </TableCell>
              </TableRow>
            ) : (
              roomTypes.map((roomType) => (
                <TableRow key={roomType.id}>
                  <TableCell>{roomType.id}</TableCell>
                  <TableCell>{roomType.name}</TableCell>
                  <TableCell>{roomType.description}</TableCell>
                  <TableCell>{roomType.capacity}</TableCell>
                  <TableCell>{roomType.default_price}</TableCell>
                  <TableCell>
                    <Button
                      variant="outlined"
                      color="primary"
                      size="small"
                      onClick={() => handleEditRoomType(roomType.id)}
                    >
                      Edit
                    </Button>
                  </TableCell>
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

export default RoomType;
