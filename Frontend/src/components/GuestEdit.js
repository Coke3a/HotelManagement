import React, { useState, useEffect } from 'react';
import { TextField, Button, Box, Typography, CircularProgress, Select, MenuItem, FormControl, InputLabel, Grid } from '@mui/material';
import { LocalizationProvider, DatePicker } from '@mui/x-date-pickers';
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns';
import { useNavigate, useParams } from 'react-router-dom';

const GuestEdit = () => {
  const token = localStorage.getItem('token');
  const navigate = useNavigate();
  const { id } = useParams();
  const [guest, setGuest] = useState({
    id: '',
    firstname: '',
    surname: '',
    identity_number: '',
    email: '',
    phone: '',
    address: '',
    date_of_birth: null,
    gender: '',
    customer_type_id: '',
    preferences: '',
    last_visit_date: null,
  });
  const [loading, setLoading] = useState(true);
  const [guestTypes, setGuestTypes] = useState([]);
  const [guestTypesLoading, setGuestTypesLoading] = useState(true);

  useEffect(() => {
    const fetchGuest = async () => {
      setLoading(true);
      try {
        const response = await fetch(`http://localhost:8080/v1/customers/${id}`, {
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
          throw new Error('Failed to fetch guest');
        }

        const data = await response.json();
        setGuest({
          ...data.data,
          date_of_birth: data.data.date_of_birth && data.data.date_of_birth !== "0001-01-01T00:00:00Z" 
            ? new Date(data.data.date_of_birth) 
            : null,
          last_visit_date: data.data.last_visit_date && data.data.last_visit_date !== "0001-01-01T00:00:00Z"
            ? new Date(data.data.last_visit_date) 
            : null,
        });
        setLoading(false);
      } catch (error) {
        console.error('Error fetching guest:', error);
        alert(error.message);
        setLoading(false);
      }
    };

    fetchGuest();
  }, [id]);

  useEffect(() => {
    const fetchGuestTypes = async () => {
      setGuestTypesLoading(true);
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
        setGuestTypes([]); // Ensure it's always an array on error
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

  const handleDateChange = (name) => (date) => {
    setGuest({ ...guest, [name]: date });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    try {
      const response = await fetch(`http://localhost:8080/v1/customers/`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify({
          ...guest,
          date_of_birth: guest.date_of_birth ? guest.date_of_birth.toISOString().split('T')[0] : null,
          last_visit_date: guest.last_visit_date ? guest.last_visit_date.toISOString().split('T')[0] : null,
        }),
      });

      if (response.status === 401) {
        handleTokenExpiration(new Error("access token has expired"), navigate);
        return;
      }

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.messages[0] || 'Failed to update guest');
      }

      const data = await response.json();
      console.log('Guest updated:', data);
      navigate('/guest');
    } catch (error) {
      console.error('Error updating guest:', error);
      alert(error.message);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async () => {
    if (window.confirm('Are you sure you want to delete this guest?')) {
      setLoading(true);
      try {
        const response = await fetch(`http://localhost:8080/v1/customers/${id}`, {
          method: 'DELETE',
          headers: {
            'Authorization': `Bearer ${token}`,
          },
        });

        if (response.status === 401) {
          handleTokenExpiration(new Error("access token has expired"), navigate);
          return;
        }

        if (!response.ok) {
          throw new Error('Failed to delete guest');
        }

        console.log('Guest deleted');
        navigate('/guest');
      } catch (error) {
        console.error('Error deleting guest:', error);
        alert(error.message);
      } finally {
        setLoading(false);
      }
    }
  };

  return (
    <Box className="form-container">
      <Typography variant="h4" gutterBottom className="form-title">
        Edit Guest
      </Typography>
      {loading || guestTypesLoading ? (
        <Box display="flex" justifyContent="center" alignItems="center" height="300px">
          <CircularProgress />
        </Box>
      ) : guestTypes.length === 0 ? (
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
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Address"
                name="address"
                value={guest.address}
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
              <Box display="flex" justifyContent="space-between" mt={2}>
                <Button
                  variant="contained"
                  color="error"
                  onClick={handleDelete}
                  disabled={loading}
                >
                  Delete Guest
                </Button>
                <Button type="submit" variant="contained" color="primary" disabled={loading}>
                  {loading ? 'Updating...' : 'Update Guest'}
                </Button>
              </Box>
            </Grid>
          </Grid>
        </form>
      )}
    </Box>
  );
};

export default GuestEdit;
