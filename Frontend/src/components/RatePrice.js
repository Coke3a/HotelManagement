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
} from '@mui/material';
import { getUserRole } from '../utils/auth';

const RatePrice = () => {
  const navigate = useNavigate();
  const [ratePrices, setRatePrices] = useState([]);
  const [loading, setLoading] = useState(true);
  const [roomTypeData, setRoomTypeData] = useState({});
  const [userRole, setUserRole] = useState('');

  const fetchRoomTypeData = async (roomTypeId) => {
    try {
      const response = await fetch(`http://localhost:8080/v1/room-types/${roomTypeId}`, {
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include'
      });

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
        const queryParams = new URLSearchParams({
          skip: '0',
          limit: '10',
        });

        const response = await fetch(`http://localhost:8080/v1/rate_prices/?${queryParams.toString()}`, {
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include'
        });

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
        } else {
          setRatePrices([]);
        }
      } catch (error) {
        console.error('Error fetching rate prices:', error);
        setRatePrices([]);
      } finally {
        setLoading(false);
      }
    };

    fetchRatePrices();
  }, []);

  useEffect(() => {
    setUserRole(getUserRole());
  }, []);

  const handleEditRatePrice = (id) => {
    navigate(`/rate_price/edit/${id}`);
  };

  const handleAddRatePrice = () => {
    navigate('/rate_price/add');
  };

  return (
    <div className="mt-5 p-5 rounded-lg bg-gray-100">
      <div className="flex justify-end mb-4">
        <Button
          variant="contained"
          className="bg-green-500 text-white hover:bg-green-600"
          onClick={handleAddRatePrice}
        >
          Add Rate Price
        </Button>
      </div>
      <TableContainer component={Paper} style={{ fontSize: '0.9rem' }}>
        <Table size="small">
          <TableHead className="bg-blue-600">
            <TableRow>
              <TableCell className="text-white font-bold">ID</TableCell>
              <TableCell className="text-white font-bold">Name</TableCell>
              <TableCell className="text-white font-bold">Description</TableCell>
              <TableCell className="text-white font-bold">Price per Night</TableCell>
              <TableCell className="text-white font-bold">Room Type</TableCell>
              <TableCell className="text-white font-bold">Actions</TableCell>
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
