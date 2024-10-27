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
  const token = localStorage.getItem('token');
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
            'Authorization': `Bearer ${token}`,
          },
        });

        if (response.status === 401) {
          handleTokenExpiration(new Error("access token has expired"), navigate);
          return;
        }

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
    <div className="table-container">
      <div className="flex justify-end mb-4">
        <Button
          variant="contained"
          className="bg-blue-600 text-white hover:bg-blue-700"
          onClick={handleAddGuestType}
        >
          Add Guest Type
        </Button>
      </div>
      <TableContainer component={Paper}>
        <Table size="small" className="compact-table">
          <TableHead>
            <TableRow>
              <TableCell>ID</TableCell>
              <TableCell>Name</TableCell>
              <TableCell>Description</TableCell>
              <TableCell>Actions</TableCell>
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
