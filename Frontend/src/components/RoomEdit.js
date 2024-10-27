import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import {
  Box,
  Typography,
  TextField,
  Button,
  CircularProgress,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Grid,
} from '@mui/material';

const RoomEdit = () => {
  const navigate = useNavigate();
  const { id } = useParams();
  const [room, setRoom] = useState({
    id: '',
    room_number: '',
    type_id: '',
    description: '',
    status: '',
    floor: '',
  });
  const [loading, setLoading] = useState(true);
  const [roomTypes, setRoomTypes] = useState([]);

  useEffect(() => {
    const fetchRoom = async () => {
      setLoading(true);
      try {
        const response = await fetch(`http://localhost:8080/v1/rooms/${id}`, {
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include'
        });

        if (!response.ok) {
          throw new Error('Failed to fetch room');
        }

        const data = await response.json();
        setRoom(data.data);
      } catch (error) {
        console.error('Error fetching room:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchRoom();
  }, [id]);

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
        setRoomTypes(data.data || []);
      } catch (error) {
        console.error('Error fetching room types:', error);
      }
    };

    fetchRoomTypes();
  }, []);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setRoom({ ...room, [name]: value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    try {
      const response = await fetch(`http://localhost:8080/v1/rooms/`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          ...room,
          status: parseInt(room.status),
          floor: parseInt(room.floor),
        }),
        credentials: 'include'
      });

      if (!response.ok) {
        throw new Error('Failed to update room');
      }

      const data = await response.json();
      console.log('Room updated:', data);
      navigate('/room');
    } catch (error) {
      console.error('Error updating room:', error);
      alert(error.message);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async () => {
    if (window.confirm('Are you sure you want to delete this room?')) {
      setLoading(true);
      try {
        const response = await fetch(`http://localhost:8080/v1/rooms/${id}`, {
          method: 'DELETE',
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include'
        });

        if (!response.ok) {
          throw new Error('Failed to delete room');
        }

        console.log('Room deleted');
        navigate('/room');
      } catch (error) {
        console.error('Error deleting room:', error);
        alert(error.message);
      } finally {
        setLoading(false);
      }
    }
  };

  return (
    <Box className="form-container">
      <Typography variant="h4" gutterBottom className="form-title">
        Edit Room
      </Typography>
      {loading ? (
        <Box display="flex" justifyContent="center" alignItems="center" height="300px">
          <CircularProgress />
        </Box>
      ) : (
        <form onSubmit={handleSubmit} className="form">
          <Grid container spacing={2}>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Room Number"
                name="room_number"
                value={room.room_number}
                onChange={handleChange}
                margin="normal"
                required
                className="form-input"
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <FormControl fullWidth margin="normal" required className="form-input">
                <InputLabel>Room Type</InputLabel>
                <Select
                  name="type_id"
                  value={room.room_type_id}
                  onChange={handleChange}
                  label="Room Type"
                >
                  {roomTypes.map((type) => (
                    <MenuItem key={type.id} value={type.id}>
                      {type.name}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} sm={6}>
              <FormControl fullWidth margin="normal" required className="form-input">
                <InputLabel>Status</InputLabel>
                <Select
                  name="status"
                  value={room.status}
                  onChange={handleChange}
                >
                  <MenuItem value={1}>Available</MenuItem>
                  <MenuItem value={2}>Occupied</MenuItem>
                  <MenuItem value={3}>Maintenance</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Floor"
                name="floor"
                value={room.floor}
                onChange={(e) => {
                  const value = e.target.value;
                  if (value === '' || /^-?\d+$/.test(value)) {
                    handleChange(e);
                  }
                }}
                margin="normal"
                required
                className="form-input"
                inputProps={{ inputMode: 'numeric', pattern: '-?[0-9]*' }}
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Description"
                name="description"
                value={room.description}
                onChange={handleChange}
                margin="normal"
                multiline
                rows={3}
                className="form-input"
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Created At"
                name="created_at"
                value={new Date(room.created_at).toLocaleString()}
                margin="normal"
                disabled
                className="form-input"
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Updated At"
                name="updated_at"
                value={new Date(room.updated_at).toLocaleString()}
                margin="normal"
                disabled
                className="form-input"
              />
            </Grid>
            <Grid item xs={12}>
              <Box display="flex" justifyContent="space-between" mt={2}>
                <Button
                  variant="contained"
                  color="error"
                  onClick={handleDelete}
                  disabled={loading}
                >
                  Delete Room
                </Button>
                <Button type="submit" variant="contained" color="primary" disabled={loading}>
                  {loading ? 'Updating...' : 'Update Room'}
                </Button>
              </Box>
            </Grid>
          </Grid>
        </form>
      )}
    </Box>
  );
};

export default RoomEdit;