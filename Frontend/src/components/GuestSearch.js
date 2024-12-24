import React, { useState } from 'react';
import { 
  TextField, 
  Button, 
  Box, 
  Typography, 
  CircularProgress,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions
} from '@mui/material';
import SearchIcon from '@mui/icons-material/Search';
import PersonAddIcon from '@mui/icons-material/PersonAdd';
import GuestAdd from './GuestAdd';

const GuestSearch = ({ onGuestSelected, showTable = true }) => {
  const [identityNumber, setIdentityNumber] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [guests, setGuests] = useState([]);
  const [openGuestModal, setOpenGuestModal] = useState(false);
  const token = localStorage.getItem('token');

  const handleSearch = async () => {
    if (!identityNumber.trim()) {
      setError('Please enter an identity number');
      return;
    }

    setLoading(true);
    setError('');

    try {
      const queryParams = new URLSearchParams({
        skip: '0',
        limit: '100',
        identity_number: identityNumber
      });

      const response = await fetch(`http://localhost:8080/v1/customers/?${queryParams.toString()}`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        throw new Error('Failed to search for guests');
      }

      const data = await response.json();
      if (data.success && data.data && data.data.customers) {
        setGuests(data.data.customers);
        if (data.data.customers.length === 0) {
          setError('No guests found with this identity number.');
        }
      } else {
        setError('No guests found with this identity number.');
        setGuests([]);
      }
    } catch (error) {
      setError(error.message);
      setGuests([]);
    } finally {
      setLoading(false);
    }
  };

  const handleSelectGuest = (guest) => {
    onGuestSelected(guest);
  };

  const handleAddNewGuest = () => {
    setOpenGuestModal(true);
  };

  const handleCloseGuestModal = () => {
    setOpenGuestModal(false);
  };

  const handleGuestAdded = (newGuest) => {
    setOpenGuestModal(false);
    onGuestSelected(newGuest);
  };

  return (
    <Box>
      <Box display="flex" gap={2} alignItems="center" mb={2}>
        <TextField
          fullWidth
          label="Guest Identity Number"
          value={identityNumber}
          onChange={(e) => setIdentityNumber(e.target.value)}
          disabled={loading}
          error={!!error}
          helperText={error}
          size="small"
        />
        <Button
          variant="contained"
          onClick={handleSearch}
          disabled={loading}
          startIcon={<SearchIcon />}
        >
          {loading ? <CircularProgress size={24} /> : 'Search'}
        </Button>
        <Button
          variant="outlined"
          onClick={handleAddNewGuest}
          disabled={loading}
          startIcon={<PersonAddIcon />}
        >
          Add New Guest
        </Button>
      </Box>

      {showTable && guests.length > 0 && (
        <TableContainer component={Paper}>
          <Table size="small">
            <TableHead>
              <TableRow>
                <TableCell>ID</TableCell>
                <TableCell>Name</TableCell>
                <TableCell>Identity Number</TableCell>
                <TableCell>Action</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {guests.map((guest) => (
                <TableRow key={guest.id}>
                  <TableCell>{guest.id}</TableCell>
                  <TableCell>{`${guest.firstname} ${guest.surname}`}</TableCell>
                  <TableCell>{guest.identity_number}</TableCell>
                  <TableCell>
                    <Button
                      variant="outlined"
                      size="small"
                      onClick={() => handleSelectGuest(guest)}
                    >
                      Select
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      )}

      <Dialog open={openGuestModal} onClose={handleCloseGuestModal} maxWidth="md" fullWidth>
        <DialogTitle>Add New Guest</DialogTitle>
        <DialogContent>
          <GuestAdd onGuestAdded={handleGuestAdded} isFromBooking={true} />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseGuestModal}>Cancel</Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default GuestSearch; 