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
} from '@mui/material';

const GuestType = () => {
  const navigate = useNavigate();
  const [guestTypes, setGuestTypes] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchGuestTypes = async () => {
      setLoading(true);
      try {
        const queryParams = new URLSearchParams({
          skip: '0',
          limit: '10',
        });

        const response = await fetch(`http://localhost:8080/v1/customer-types/?${queryParams.toString()}`, {
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include'
        });

        if (!response.ok) {
          throw new Error('Failed to fetch guest types');
        }

        const data = await response.json();
        setGuestTypes(data.data);
      } catch (error) {
        console.error('Error fetching guest types:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchGuestTypes();
  }, []);

  const handleEditGuestType = (id) => {
    navigate(`/guest-type/edit/${id}`);
  };

  const handleAddGuestType = () => {
    navigate('/guest-type/add');
  };

  return (
    <div className="mt-5 p-5 rounded-lg bg-gray-100">
      <div className="flex justify-end mb-4">
        <Button
          variant="contained"
          className="bg-green-500 text-white hover:bg-green-600"
          onClick={handleAddGuestType}
        >
          Add Guest Type
        </Button>
      </div>
      <TableContainer component={Paper} style={{ fontSize: '0.9rem' }}>
        <Table size="small">
          <TableHead className="bg-blue-600">
            <TableRow>
              <TableCell className="text-white font-bold">ID</TableCell>
              <TableCell className="text-white font-bold">Name</TableCell>
              <TableCell className="text-white font-bold">Description</TableCell>
              <TableCell className="text-white font-bold">Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {loading ? (
              <TableRow>
                <TableCell colSpan={4} align="center">
                  <Box display="flex" justifyContent="center" alignItems="center" height="50px">
                    <CircularProgress />
                  </Box>
                </TableCell>
              </TableRow>
            ) : !guestTypes || guestTypes.length === 0 ? (
              <TableRow>
                <TableCell colSpan={4} align="center">
                  <Typography variant="body1">No guest types found</Typography>
                </TableCell>
              </TableRow>
            ) : (
              guestTypes.map((guestType, index) => (
                <TableRow key={guestType.id} className={index % 2 === 0 ? "bg-gray-50" : "bg-white"}>
                  <TableCell>{guestType.id}</TableCell>
                  <TableCell>{guestType.name}</TableCell>
                  <TableCell>{guestType.description}</TableCell>
                  <TableCell>
                    <Button
                      variant="outlined"
                      color="primary"
                      size="small"
                      onClick={() => handleEditGuestType(guestType.id)}
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

export default GuestType;
