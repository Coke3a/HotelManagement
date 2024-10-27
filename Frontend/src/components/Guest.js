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
  FormControl,
  InputLabel,
  Select,
  MenuItem,
} from '@mui/material';

const Guest = () => {
  const navigate = useNavigate();
  const [guests, setGuests] = useState([]);
  const [loading, setLoading] = useState(true);
  const [filters, setFilters] = useState({
    id: '',
    name: '',
    email: '',
    membership_status: '',
  });
  const [guestTypes, setGuestTypes] = useState([]);

  useEffect(() => {
    const fetchGuests = async () => {
      setLoading(true);
      try {
        const queryParams = new URLSearchParams({
          skip: '0',
          limit: '10',
          ...(filters.id && { id: filters.id }),
          ...(filters.name && { name: filters.name }),
          ...(filters.email && { email: filters.email }),
          ...(filters.membership_status && { membership_status: filters.membership_status }),
        });

        const response = await fetch(`http://localhost:8080/v1/customers/?${queryParams.toString()}`, {
          headers: {
            'Content-Type': 'application/json',
            'Access-Control-Expose-Headers': 'X-My-Custom-Header, X-Another-Custom-Header'
          },
          credentials: 'include'
        });

        if (!response.ok) {
          throw new Error('Failed to fetch guests');
        }

        const data = await response.json();
        setGuests(data.data.customers || []); // Update this line
      } catch (error) {
        console.error('Error fetching guests:', error);
        setGuests([]);
      } finally {
        setLoading(false);
      }
    };

    fetchGuests();
  }, [filters]);

  useEffect(() => {
    const fetchGuestTypes = async () => {
      try {
        const response = await fetch('http://localhost:8080/v1/customer-types/', {
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include'
        });

        if (!response.ok) {
          throw new Error('Failed to fetch guest types');
        }

        const data = await response.json();
        setGuestTypes(data.data || []); // Ensure it's always an array
      } catch (error) {
        console.error('Error fetching guest types:', error);
        setGuestTypes([]); // Set to empty array on error
      }
    };

    fetchGuestTypes();
  }, []);

  const handleEditGuest = (id) => {
    navigate(`/guest/edit/${id}`);
  };

  const handleAddGuest = () => {
    navigate('/guest/add', { state: { isFromBooking: false } });
  };

  return (
    <div className="mt-5 p-5 rounded-lg bg-gray-100">
      <div className="flex justify-end mb-4">
        <Button
          variant="contained"
          className="bg-primary text-white hover:bg-primary-dark"
          onClick={handleAddGuest}
        >
          Add Guest
        </Button>
      </div>
      <div className="flex flex-wrap gap-4 mb-4">
        <TextField
          label="Guest ID"
          value={filters.id}
          onChange={(e) => setFilters({ ...filters, id: e.target.value })}
          className="w-40"
        />
        <TextField
          label="Firstname"
          value={filters.firstname}
          onChange={(e) => setFilters({ ...filters, firstname: e.target.value })}
          className="w-40"
        />
        <TextField
          label="Surname"
          value={filters.surname}
          onChange={(e) => setFilters({ ...filters, surname: e.target.value })}
          className="w-40"
        />
        <TextField
          label="Email"
          value={filters.email}
          onChange={(e) => setFilters({ ...filters, email: e.target.value })}
          className="w-40"
        />
        <FormControl className="w-40">
          <InputLabel>Guest Type</InputLabel>
          <Select
            value={filters.customer_type_id || ''}
            onChange={(e) => setFilters({ ...filters, customer_type_id: e.target.value })}
          >
            <MenuItem value="">All</MenuItem>
            {guestTypes.map((type) => (
              <MenuItem key={type.id} value={type.id}>{type.name}</MenuItem>
            ))}
          </Select>
        </FormControl>
      </div>
      <TableContainer component={Paper} style={{ fontSize: '0.9rem' }}>
        <Table size="small">
          <TableHead className="bg-blue-600">
            <TableRow>
              <TableCell className="text-white font-bold">ID</TableCell>
              <TableCell className="text-white font-bold">Name</TableCell>
              <TableCell className="text-white font-bold">Email</TableCell>
              <TableCell className="text-white font-bold">Phone</TableCell>
              <TableCell className="text-white font-bold">Type</TableCell>
              <TableCell className="text-white font-bold">Last Visit</TableCell>
              <TableCell className="text-white font-bold">Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {loading ? (
              <TableRow>
                <TableCell colSpan={7} align="center">
                  <CircularProgress size={24} />
                </TableCell>
              </TableRow>
            ) : !guests || guests.length === 0 ? (
              <TableRow>
                <TableCell colSpan={7} align="center">
                  <Typography variant="body2">No guests found</Typography>
                </TableCell>
              </TableRow>
            ) : (
              guests.map((guest) => (
                <TableRow key={guest.id}>
                  <TableCell>{guest.id}</TableCell>
                  <TableCell>{`${guest.firstname} ${guest.surname}`}</TableCell>
                  <TableCell>{guest.email}</TableCell>
                  <TableCell>{guest.phone}</TableCell>
                  <TableCell>
                    {guestTypes.find(type => type.id === guest.customer_type_id)?.name || 'N/A'}
                  </TableCell>
                  <TableCell>
                    {guest.last_visit_date && guest.last_visit_date !== "0001-01-01T00:00:00Z"
                      ? new Date(guest.last_visit_date).toLocaleDateString()
                      : 'N/A'}
                  </TableCell>
                  <TableCell>
                    <Button
                      variant="outlined"
                      color="primary"
                      size="small"
                      onClick={() => handleEditGuest(guest.id)}
                    >
                      Detail
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

export default Guest;
