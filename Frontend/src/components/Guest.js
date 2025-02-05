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
  Grid,
  TextField,
  FormControl,
  Select,
  MenuItem,
} from '@mui/material';
import { handleTokenExpiration } from '../utils/api';
import SearchIcon from '@mui/icons-material/Search';

const Guest = () => {
  const token = localStorage.getItem('token');
  const navigate = useNavigate();
  const [guests, setGuests] = useState([]);
  const [guestTypes, setGuestTypes] = useState([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(25);
  const [totalCount, setTotalCount] = useState(0);
  const [error, setError] = useState(null);
  const [filters, setFilters] = useState({
    id: '',
    firstname: '',
    surname: '',
    email: '',
    phone: '',
    customer_type_id: '',
    identity_number: '',
  });

  const fetchGuests = async (page, rowsPerPage) => {
    setLoading(true);
    try {
      const skip = page * rowsPerPage;
      const queryParams = new URLSearchParams({
        skip: skip.toString(),
        limit: rowsPerPage.toString(),
        ...(filters.id && { id: filters.id }),
        ...(filters.firstname && { firstname: filters.firstname }),
        ...(filters.surname && { surname: filters.surname }),
        ...(filters.email && { email: filters.email }),
        ...(filters.phone && { phone: filters.phone }),
        ...(filters.customer_type_id && { customer_type_id: filters.customer_type_id }),
        ...(filters.identity_number && { identity_number: filters.identity_number }),
      });

      const response = await fetch(`http://localhost:8080/v1/customers/?${queryParams.toString()}`, {
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
        throw new Error('Failed to fetch guests');
      }

      const data = await response.json();
      if (data && data.data && data.data.customers) {
        setGuests(data.data.customers);
        if (data.data.meta && typeof data.data.meta.total === 'number') {
          setTotalCount(data.data.meta.total);
        } else {
          setTotalCount(data.data.customers.length);
        }
      } else {
        setGuests([]);
        setTotalCount(0);
      }
    } catch (error) {
      console.error('Error fetching guests:', error);
      setError(error.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchGuests(page, rowsPerPage);
  }, [page, rowsPerPage, navigate]);

  const handlePageChange = (event, newPage) => {
    setPage(newPage - 1);
  };

  useEffect(() => {
    const fetchGuestTypes = async () => {
      try {
        const response = await fetch('http://localhost:8080/v1/customer-types/', {
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
        if (data && data.data && data.data.customerTypes) {
          setGuestTypes(data.data.customerTypes);
        } else {
          setGuestTypes([]);
        }
      } catch (error) {
        console.error('Error fetching guest types:', error);
        setGuestTypes([]);
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
          onClick={handleAddGuest}
        >
          Add Guest
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
              <option value={25}>25 per page</option>
              <option value={50}>50 per page</option>
              <option value={75}>75 per page</option>
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

      <Box sx={{ mb: 3 }}>
        <Paper sx={{ p: 2 }}>
          <Grid container spacing={2}>
            {/* First Row */}
            <Grid item xs={12}>
              <Stack direction={{ xs: 'column', md: 'row' }} spacing={2}>
                <TextField
                  size="small"
                  placeholder="Guest ID"
                  value={filters.id}
                  onChange={(e) => setFilters({ ...filters, id: e.target.value })}
                  sx={{ width: { xs: '100%', md: '20%' } }}
                />
                <TextField
                  size="small"
                  placeholder="First Name"
                  value={filters.firstname}
                  onChange={(e) => setFilters({ ...filters, firstname: e.target.value })}
                  sx={{ width: { xs: '100%', md: '20%' } }}
                />
                <TextField
                  size="small"
                  placeholder="Last Name"
                  value={filters.surname}
                  onChange={(e) => setFilters({ ...filters, surname: e.target.value })}
                  sx={{ width: { xs: '100%', md: '20%' } }}
                />
                <TextField
                  size="small"
                  placeholder="Identity Number"
                  value={filters.identity_number}
                  onChange={(e) => setFilters({ ...filters, identity_number: e.target.value })}
                  sx={{ width: { xs: '100%', md: '20%' } }}
                />
                <FormControl size="small" sx={{ width: { xs: '100%', md: '20%' } }}>
                  <Select
                    value={filters.customer_type_id}
                    onChange={(e) => setFilters({ ...filters, customer_type_id: e.target.value })}
                    displayEmpty
                  >
                    <MenuItem value="">All Types</MenuItem>
                    {guestTypes.map((type) => (
                      <MenuItem key={type.id} value={type.id}>{type.name}</MenuItem>
                    ))}
                  </Select>
                </FormControl>
              </Stack>
            </Grid>

            {/* Second Row */}
            <Grid item xs={12}>
              <Stack direction={{ xs: 'column', md: 'row' }} spacing={2}>
                <TextField
                  size="small"
                  placeholder="Email"
                  value={filters.email}
                  onChange={(e) => setFilters({ ...filters, email: e.target.value })}
                  sx={{ width: { xs: '100%', md: '30%' } }}
                />
                <TextField
                  size="small"
                  placeholder="Phone"
                  value={filters.phone}
                  onChange={(e) => setFilters({ ...filters, phone: e.target.value })}
                  sx={{ width: { xs: '100%', md: '30%' } }}
                />
                <Stack direction="row" spacing={1} sx={{ width: { xs: '100%', md: '40%' } }}>
                  <Button
                    variant="contained"
                    onClick={() => {
                      setPage(0);
                      fetchGuests(0, rowsPerPage);
                    }}
                    startIcon={<SearchIcon />}
                    size="small"
                    fullWidth
                  >
                    Search
                  </Button>
                  <Button
                    variant="outlined"
                    onClick={() => {
                      setFilters({
                        id: '',
                        firstname: '',
                        surname: '',
                        email: '',
                        phone: '',
                        customer_type_id: '',
                      });
                      setPage(0);
                      fetchGuests(0, rowsPerPage);
                    }}
                    size="small"
                    sx={{ minWidth: 'auto', width: '40px', p: 0 }}
                  >
                    ×
                  </Button>
                </Stack>
              </Stack>
            </Grid>
          </Grid>
        </Paper>
      </Box>

      <TableContainer component={Paper}>
        <Table size="small" className="compact-table">
          <TableHead>
            <TableRow>
              <TableCell>ID</TableCell>
              <TableCell>Name</TableCell>
              <TableCell>Email</TableCell>
              <TableCell>Phone</TableCell>
              <TableCell>Type</TableCell>
              <TableCell>Actions</TableCell>
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
