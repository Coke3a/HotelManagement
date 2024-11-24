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
  TextField,
  Stack,
  Pagination,
} from '@mui/material';
import { getUserRole } from '../utils/auth';

const RatePrice = () => {
  const token = localStorage.getItem('token');
  const navigate = useNavigate();
  const [ratePrices, setRatePrices] = useState([]);
  const [loading, setLoading] = useState(true);
  const [roomTypeData, setRoomTypeData] = useState({});
  const [userRole, setUserRole] = useState('');
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [totalCount, setTotalCount] = useState(0);
  const [error, setError] = useState(null);

  const fetchRoomTypeData = async (roomTypeId) => {
    try {
      const response = await fetch(`http://localhost:8080/v1/room-types/${roomTypeId}`, {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
      });

      if (response.status === 401) {
        handleTokenExpiration(new Error("access token has expired"), navigate);
        return;
      }

      if (!response) {
        throw new Error('No response received from the server');
      }

      if (!response.ok) {
        throw new Error('Failed to fetch room type data');
      }

      const data = await response.json();
      return data.data;
    } catch (error) {
      console.error('Error fetching room type data:', error);
      return null;
    }
  };

  useEffect(() => {
    const fetchRatePrices = async () => {
      setLoading(true);
      try {
        const skip = page * rowsPerPage;
        const queryParams = new URLSearchParams({
          skip: skip.toString(),
          limit: rowsPerPage.toString(),
        });

        const response = await fetch(`http://localhost:8080/v1/rate_prices/?${queryParams.toString()}`, {
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json',
          },
        });

        if (response.status === 401) {
          handleTokenExpiration(new Error("access token has expired"), navigate);
          return;
        }

        if (!response.ok) {
          throw new Error('Failed to fetch rate prices');
        }

        const data = await response.json();
        if (data.success && data.data && data.data.ratePrices) {
          const ratePricesWithRoomTypeData = await Promise.all(
            data.data.ratePrices.map(async (ratePrice) => {
              const roomTypeData = await fetchRoomTypeData(ratePrice.room_type_id);
              return { ...ratePrice, roomTypeData };
            })
          );
          setRatePrices(ratePricesWithRoomTypeData);
          if (data.data.meta && typeof data.data.meta.total === 'number') {
            setTotalCount(data.data.meta.total);
          } else {
            setTotalCount(data.data.ratePrices.length);
          }
        } else {
          setRatePrices([]);
          setTotalCount(0);
        }
      } catch (error) {
        console.error('Error fetching rate prices:', error);
        setRatePrices([]);
        setTotalCount(0);
      } finally {
        setLoading(false);
      }
    };

    fetchRatePrices();
  }, [page, rowsPerPage, navigate]);

  useEffect(() => {
    setUserRole(getUserRole());
  }, []);

  const handleEditRatePrice = (id) => {
    navigate(`/rate_price/edit/${id}`);
  };

  const handleAddRatePrice = () => {
    navigate('/rate_price/add');
  };

  const handlePageChange = (event, newPage) => {
    setPage(newPage - 1);
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
        <Button
          variant="contained"
          className="bg-blue-600 text-white hover:bg-blue-700"
          onClick={handleAddRatePrice}
        >
          Add Rate Price
        </Button>
        
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
              <TableCell>Price per Night</TableCell>
              <TableCell>Room Type</TableCell>
              <TableCell>Actions</TableCell>
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
            ) : ratePrices.length === 0 ? (
              <TableRow>
                <TableCell colSpan={6} align="center">
                  <Typography variant="body1">No rate prices found</Typography>
                </TableCell>
              </TableRow>
            ) : (
              ratePrices.map((ratePrice, index) => (
                <TableRow key={ratePrice.id} className={index % 2 === 0 ? "bg-gray-50" : "bg-white"}>
                  <TableCell>{ratePrice.id}</TableCell>
                  <TableCell>{ratePrice.name}</TableCell>
                  <TableCell>{ratePrice.description}</TableCell>
                  <TableCell>{ratePrice.price_per_night}</TableCell>
                  <TableCell>{ratePrice.roomTypeData ? ratePrice.roomTypeData.name : 'N/A'}</TableCell>
                  <TableCell>
                    <Button
                      variant="outlined"
                      color="primary"
                      size="small"
                      onClick={() => handleEditRatePrice(ratePrice.id)}
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
    </div>
  );
};

export default RatePrice;
