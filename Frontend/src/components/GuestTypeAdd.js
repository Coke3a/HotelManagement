import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Box,
  Typography,
  TextField,
  Button,
  CircularProgress,
} from '@mui/material';

const GuestTypeAdd = () => {
  const navigate = useNavigate();
  const [guestType, setGuestType] = useState({
    name: '',
    description: '',
  });
  const [loading, setLoading] = useState(false);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setGuestType({ ...guestType, [name]: value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    try {
      const response = await fetch('http://localhost:8080/v1/customer-types/', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(guestType),
        credentials: 'include'
      });

      if (response.status === 409) {
        throw new Error('Guest type already exists. Please use a different name.');
      }

      if (!response.ok) {
        throw new Error('Failed to create guest type');
      }

      const data = await response.json();
      console.log('Guest type created:', data);
      navigate('/guest-type');
    } catch (error) {
      console.error('Error adding guest type:', error);
      alert(error.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Box className="form-container">
      <Typography variant="h4" gutterBottom className="form-title">
        Add New Guest Type
      </Typography>
      {loading ? (
        <Box display="flex" justifyContent="center" alignItems="center" height="300px">
          <CircularProgress />
        </Box>
      ) : (
        <form onSubmit={handleSubmit} className="form">
          <TextField
            fullWidth
            label="Name"
            name="name"
            value={guestType.name}
            onChange={handleChange}
            margin="normal"
            required
            className="form-input"
          />
          <TextField
            fullWidth
            label="Description"
            name="description"
            value={guestType.description}
            onChange={handleChange}
            margin="normal"
            multiline
            rows={3}
            className="form-input"
          />
          <Button type="submit" variant="contained" color="primary" className="form-submit" disabled={loading}>
            {loading ? 'Adding...' : 'Add Guest Type'}
          </Button>
        </form>
      )}
    </Box>
  );
};

export default GuestTypeAdd;
