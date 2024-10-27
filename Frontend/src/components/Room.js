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

const Room = () => {
  const navigate = useNavigate();
  const [rooms, setRooms] = useState([]);
  const [loading, setLoading] = useState(true);
  const [roomTypes, setRoomTypes] = useState({});
  const [userRole, setUserRole] = useState('');
  const [roomTypeMap, setRoomTypeMap] = useState({});

  useEffect(() => {
    setUserRole(getUserRole());
  }, []);

  useEffect(() => {
    const fetchRooms = async () => {
      setLoading(true);
      try {
        const queryParams = new URLSearchParams({
          skip: '0',
          limit: '10',
        });

        const response = await fetch(`http://localhost:8080/v1/rooms/?${queryParams.toString()}`, {
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include'
        });

        if (!response.ok) {
          throw new Error('Failed to fetch rooms');
        }

        const data = await response.json();
        setRooms(data.data.rooms);
      } catch (error) {
        console.error('Error fetching rooms:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchRooms();
  }, []);

  useEffect(() => {
    const fetchRoomTypes = async () => {
      try {
        const response = await fetch('http://localhost:8080/v1/room-types/', {
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include'
        });

        if (!response.ok) {
          throw new Error('Failed to fetch room types');
        }

        const data = await response.json();
        const roomTypesMap = {};
        data.data.forEach(type => {
          roomTypesMap[type.id] = type.name;
        });
        setRoomTypes(roomTypesMap);

        // Create roomTypeMap here after fetching room types
        const newRoomTypeMap = data.data.reduce((acc, roomType) => {
          acc[roomType.id] = roomType.name;
          return acc;
        }, {});
        setRoomTypeMap(newRoomTypeMap);
      } catch (error) {
        console.error('Error fetching room types:', error);
      }
    };

    fetchRoomTypes();
  }, []);

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
    <div className="mt-5 p-5 rounded-lg bg-gray-100">
      {userRole === UserRoleEnum.ADMIN && (
        <div className="flex justify-end mb-4">
          <Button
            variant="contained"
            className="bg-green-500 text-white hover:bg-green-600"
            onClick={handleAddRoom}
          >
            Add Room
          </Button>
        </div>
      )}
      <TableContainer component={Paper} style={{ fontSize: '0.9rem' }}>
        <Table size="small">
          <TableHead className="bg-blue-600">
            <TableRow>
              <TableCell className="text-white font-bold">ID</TableCell>
              <TableCell className="text-white font-bold">Room Number</TableCell>
              <TableCell className="text-white font-bold">Room Type</TableCell>
              <TableCell className="text-white font-bold">Status</TableCell>
              <TableCell className="text-white font-bold">Floor</TableCell>
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
            ) : !rooms || rooms.length === 0 ? (
              <TableRow>
                <TableCell colSpan={6} align="center">
                  <Typography variant="body1">No rooms found</Typography>
                </TableCell>
              </TableRow>
            ) : (
              rooms.map((room, index) => (
                <TableRow key={room.id} className={index % 2 === 0 ? "bg-gray-50" : "bg-white"}>
                  <TableCell>{room.id}</TableCell>
                  <TableCell>{room.room_number}</TableCell>
                  <TableCell>{roomTypeMap[String(room.room_type_id)] || 'Unknown'}</TableCell>
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
