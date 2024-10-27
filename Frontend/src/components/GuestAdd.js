import React, { useState, useEffect } from 'react';
import { TextField, Button, Box, Typography, CircularProgress, Select, MenuItem, FormControl, InputLabel, Grid } from '@mui/material';
import { LocalizationProvider, DatePicker } from '@mui/x-date-pickers';
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns';
import { useNavigate } from 'react-router-dom';

const GuestAdd = ({ onGuestAdded, isFromBooking = false }) => {
  const navigate = useNavigate();
  const [guest, setGuest] = useState({
    firstname: '',
    surname: '',
    identity_number: '',
    email: '',
    phone: '',
    date_of_birth: null,
    gender: '',
    customer_type_id: '',
    preferences: '',
  });
  const [loading, setLoading] = useState(false);
  const [guestTypes, setGuestTypes] = useState([]);
  const [guestTypesLoading, setGuestTypesLoading] = useState(true);

  useEffect(() => {
    const fetchGuestTypes = async () => {
      setGuestTypesLoading(true);
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
      } finally {
        setGuestTypesLoading(false);
      }
    };

    fetchGuestTypes();
  }, []);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setGuest({ ...guest, [name]: value });
  };

  const handleDateChange = (date) => {
    setGuest({ ...guest, date_of_birth: date ? new Date(date) : null });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    try {
      const response = await fetch('http://localhost:8080/v1/customers/', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Access-Control-Expose-Headers': 'X-My-Custom-Header, X-Another-Custom-Header'
        },
        body: JSON.stringify(guest),
        credentials: 'include'
      });

      if (response.status === 409) {
        throw new Error('A guest with this information already exists.');
      }

      if (!response.ok) {
        throw new Error('Failed to create guest');
      }

      const data = await response.json();
      console.log('Guest created:', data);

      if (isFromBooking && onGuestAdded) {
        onGuestAdded(data.data);
      } else {
        navigate('/guest');
      }
    } catch (error) {
      console.error('Error adding guest:', error);
      alert(error.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Box className="form-container">
      <Typography variant="h4" gutterBottom className="form-title">
        Add New Guest
      </Typography>
      {loading || guestTypesLoading ? (
        <Box display="flex" justifyContent="center" alignItems="center" height="300px">
          <CircularProgress />
        </Box>
      ) : !guestTypes || guestTypes.length === 0 ? (
        <Box display="flex" flexDirection="column" alignItems="center">
          <Typography variant="h6" gutterBottom>
            No guest types available. Please add a guest type first.
          </Typography>
          <Button
            variant="contained"
            color="primary"
            onClick={() => navigate('/guest-type/add')}
          >
            Add Guest Type
          </Button>
        </Box>
      ) : (
        <form onSubmit={handleSubmit} className="form">
          <Grid container spacing={2}>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="First Name"
                name="firstname"
                value={guest.firstname}
                onChange={handleChange}
                margin="normal"
                required
                className="form-input"
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Surname"
                name="surname"
                value={guest.surname}
                onChange={handleChange}
                margin="normal"
                required
                className="form-input"
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Identity Number"
                name="identity_number"
                value={guest.identity_number}
                onChange={handleChange}
                margin="normal"
                required
                className="form-input"
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <FormControl fullWidth margin="normal" required className="form-input">
                <InputLabel>Guest Type</InputLabel>
                <Select
                  name="customer_type_id"
                  value={guest.customer_type_id}
                  onChange={handleChange}
                >
                  {guestTypes.map((type) => (
                    <MenuItem key={type.id} value={type.id}>{type.name}</MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} sm={6}>
              <FormControl fullWidth margin="normal" className="form-input">
                <InputLabel>Gender</InputLabel>
                <Select
                  name="gender"
                  value={guest.gender}
                  onChange={handleChange}
                >
                  <MenuItem value="male">Male</MenuItem>
                  <MenuItem value="female">Female</MenuItem>
                  <MenuItem value="other">Other</MenuItem>
                </Select>
              </FormControl>
            </Grid>
            <Grid item xs={12} sm={6}>
              <LocalizationProvider dateAdapter={AdapterDateFns}>
                <DatePicker
                  label="Date of Birth"
                  value={guest.date_of_birth}
                  onChange={handleDateChange}
                  renderInput={(params) => <TextField {...params} fullWidth margin="normal" className="form-input" />}
                  inputFormat="yyyy-MM-dd"
                />
              </LocalizationProvider>
            </Grid>

            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Phone"
                name="phone"
                value={guest.phone}
                onChange={handleChange}
                margin="normal"
                className="form-input"
              />
            </Grid>
            <Grid item xs={12} sm={6}>
              <TextField
                fullWidth
                label="Email"
                name="email"
                value={guest.email}
                onChange={handleChange}
                margin="normal"
                className="form-input"
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Preferences"
                name="preferences"
                value={guest.preferences}
                onChange={handleChange}
                margin="normal"
                multiline
                rows={3}
                className="form-input"
              />
            </Grid>
            <Grid item xs={12}>
              <Button type="submit" variant="contained" color="primary" className="form-submit" disabled={loading}>
                {loading ? 'Adding...' : 'Add Guest'}
              </Button>
            </Grid>
          </Grid>
        </form>
      )}
    </Box>
  );
};

export default GuestAdd;
