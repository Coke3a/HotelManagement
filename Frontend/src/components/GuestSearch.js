import React, { useState, useEffect } from 'react';
import { 
  TextField, 
  Button, 
  Box,
  Grid,
  Typography,
  Card,
  CardContent,
  Divider,
  Alert
} from '@mui/material';
import SearchIcon from '@mui/icons-material/Search';
import PersonAddIcon from '@mui/icons-material/PersonAdd';
import { useNavigate } from 'react-router-dom';
import { handleTokenExpiration } from '../utils/api';

const GuestSearch = ({ onGuestSelected, currentGuestId, guests, setGuests, onAddNewGuest, initialSelectedGuest }) => {
  const [identityNumber, setIdentityNumber] = useState('');
  const [error, setError] = useState('');
  const [selectedGuest, setSelectedGuest] = useState(initialSelectedGuest || null);
  const token = localStorage.getItem('token');
  const navigate = useNavigate();

  useEffect(() => {
    if (guests.length > 0) {
      const currentGuest = guests.find(g => g.id === currentGuestId);
      if (currentGuest) {
        setSelectedGuest(currentGuest);
        setIdentityNumber(currentGuest.identity_number);
      }
    }
  }, [guests, currentGuestId]);

  const handleSearch = async () => {
    if (!identityNumber.trim()) return;

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

      if (response.status === 401) {
        handleTokenExpiration(new Error("access token has expired"), navigate);
        return;
      }

      if (!response.ok) {
        throw new Error('Failed to search for guests');
      }

      const data = await response.json();
      if (data.success && data.data && data.data.customers && data.data.customers.length > 0) {
        const foundGuest = data.data.customers[0]; // Get the first matching guest
        setSelectedGuest(foundGuest);
        onGuestSelected(foundGuest.id);
        setError('');
      } else {
        setSelectedGuest(null);
        onGuestSelected('');
        setError('No guest found with this identity number');
      }
    } catch (error) {
      console.error('Error searching guests:', error);
      setSelectedGuest(null);
      onGuestSelected('');
      setError(error.message);
    }
  };

  return (
    <Box sx={{ p: 2, border: '1px solid #e0e0e0', borderRadius: 1 }}>
      <Typography variant="subtitle1" gutterBottom sx={{ fontWeight: 'medium', color: 'text.secondary' }}>
        Guest Information
      </Typography>
      <Grid container spacing={2} alignItems="center">
        <Grid item xs={12} md={6}>
          <TextField
            fullWidth
            size="small"
            label="Guest Identity Number"
            value={identityNumber}
            onChange={(e) => setIdentityNumber(e.target.value)}
            error={!!error}
            helperText={error}
            onKeyPress={(e) => {
              if (e.key === 'Enter') {
                handleSearch();
              }
            }}
          />
        </Grid>
        <Grid item xs={12} md={3}>
          <Button
            fullWidth
            variant="contained"
            onClick={handleSearch}
            startIcon={<SearchIcon />}
            size="small"
            sx={{ height: '40px' }}
          >
            Search Guest
          </Button>
        </Grid>
        <Grid item xs={12} md={3}>
          <Button
            fullWidth
            variant="outlined"
            onClick={onAddNewGuest}
            startIcon={<PersonAddIcon />}
            size="small"
            sx={{ height: '40px' }}
          >
            Add New Guest
          </Button>
        </Grid>
        
        {selectedGuest && (
          <Grid item xs={12}>
            <Card variant="outlined" sx={{ mt: 2 }}>
              <CardContent>
                <Grid container spacing={2}>
                  <Grid item xs={12} md={6}>
                    <Typography variant="body2" color="text.secondary">Name</Typography>
                    <Typography variant="body1">{`${selectedGuest.firstname} ${selectedGuest.surname}`}</Typography>
                  </Grid>
                  <Grid item xs={12} md={6}>
                    <Typography variant="body2" color="text.secondary">Identity Number</Typography>
                    <Typography variant="body1">{selectedGuest.identity_number}</Typography>
                  </Grid>
                  {selectedGuest.phone && (
                    <Grid item xs={12} md={6}>
                      <Typography variant="body2" color="text.secondary">Phone</Typography>
                      <Typography variant="body1">{selectedGuest.phone}</Typography>
                    </Grid>
                  )}
                  {selectedGuest.email && (
                    <Grid item xs={12} md={6}>
                      <Typography variant="body2" color="text.secondary">Email</Typography>
                      <Typography variant="body1">{selectedGuest.email}</Typography>
                    </Grid>
                  )}
                </Grid>
              </CardContent>
            </Card>
          </Grid>
        )}
      </Grid>
    </Box>
  );
};

export default GuestSearch; 