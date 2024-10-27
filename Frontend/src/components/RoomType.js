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
import { getUserRole } from '../utils/auth';
import { UserRoleEnum } from '../utils/userRoleEnum';

const RoomType = () => {
  const navigate = useNavigate();
  const [roomTypes, setRoomTypes] = useState([]);
  const [loading, setLoading] = useState(true);
  const [userRole, setUserRole] = useState('');

  useEffect(() => {
    setUserRole(getUserRole());
  }, []);

  useEffect(() => {
    const fetchRoomTypes = async () => {
      setLoading(true);
      try {
        const queryParams = new URLSearchParams({
          skip: '0',
          limit: '10',
        });

        const response = await fetch(`http://localhost:8080/v1/room-types/?${queryParams.toString()}`, {
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include'
        });

        if (!response.ok) {
          throw new Error('Failed to fetch room types');
        }

        const data = await response.json();
        if (data && data.data) {
          setRoomTypes(Array.isArray(data.data) ? data.data : []);
        } else {
          console.error('Unexpected data structure:', data);
          setRoomTypes([]);
        }
      } catch (error) {
        console.error('Error fetching room types:', error);
        setRoomTypes([]);
      } finally {
        setLoading(false);
      }
    };

    fetchRoomTypes();
  }, []);

  const handleEditRoomType = (id) => {
    navigate(`/room-type/edit/${id}`);
  };

  const handleAddRoomType = () => {
    navigate('/room-type/add');
  };

  return (
    <div className="mt-5 p-5 rounded-lg bg-gray-100">
      {userRole === UserRoleEnum.ADMIN && (
        <div className="flex justify-end mb-4">
          <Button
            variant="contained"
            className="bg-green-500 text-white hover:bg-green-600"
            onClick={handleAddRoomType}
          >
            Add Room Type
          </Button>
        </div>
      )}
      <TableContainer component={Paper} style={{ fontSize: '0.9rem' }}>
        <Table size="small">
          <TableHead className="bg-blue-600">
            <TableRow>
              <TableCell className="text-white font-bold">ID</TableCell>
              <TableCell className="text-white font-bold">Name</TableCell>
              <TableCell className="text-white font-bold">Description</TableCell>
              <TableCell className="text-white font-bold">Capacity</TableCell>
              <TableCell className="text-white font-bold">Default Price</TableCell>
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
            ) : !roomTypes || roomTypes.length === 0 ? (
              <TableRow>
                <TableCell colSpan={6} align="center">
                  <Typography variant="body1">No room types found</Typography>
                </TableCell>
              </TableRow>
            ) : (
              roomTypes.map((roomType, index) => (
                <TableRow key={roomType.id} className={index % 2 === 0 ? "bg-gray-50" : "bg-white"}>
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
    </div>
  );
};

export default RoomType;
