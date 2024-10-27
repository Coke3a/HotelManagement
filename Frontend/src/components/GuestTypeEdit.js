import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import {
  Box,
  Typography,
  TextField,
  Button,
  CircularProgress,
} from '@mui/material';

const GuestTypeEdit = () => {
  const navigate = useNavigate();
  const { id } = useParams();
  const [guestType, setGuestType] = useState({
    id: '',
    name: '',
    description: '',
  });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchGuestType = async () => {
      setLoading(true);
      try {
        const response = await fetch(`http://localhost:8080/v1/customer-types/${id}`, {
          headers: {
            'Content-Type': 'application/json',
          },
          credentials: 'include'
        });

        if (!response.ok) {
          throw new Error('Failed to fetch guest type');
        }

        const data = await response.json();
        setGuestType(data.data);
      } catch (error) {
        console.error('Error fetching guest type:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchGuestType();
  }, [id]);

  const handleChange = (e) => {
    const { name, value } = e.target;
    setGuestType({ ...guestType, [name]: value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    try {
      const response = await fetch(`http://localhost:8080/v1/customer-types/`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(guestType),
        credentials: 'include'
      });

      if (!response.ok) {
        throw new Error('Failed to update guest type');
      }

      const data = await response.json();
      console.log('Guest type updated:', data);
      navigate('/guest-type');
    } catch (error) {
      console.error('Error updating guest type:', error);
      alert(error.message);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async () => {
    if (window.confirm('Are you sure you want to delete this guest type?')) {
      setLoading(true);
      try {
        const response = await fetch(`http://localhost:8080/v1/customer-types/${id}`, {
          method: 'DELETE',
          credentials: 'include'
        });

        if (!response.ok) {
          throw new Error('Failed to delete guest type');
        }

        console.log('Guest type deleted');
        navigate('/guest-type');
      } catch (error) {
        console.error('Error deleting guest type:', error);
        alert(error.message);
      } finally {
        setLoading(false);
      }
    }
  };

  return (
    <Box className="form-container">
      <Typography variant="h4" gutterBottom className="form-title">
        Edit Guest Type
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
          <Box display="flex" justifyContent="space-between" mt={2}>
            <Button
              variant="contained"
              color="error"
              onClick={handleDelete}
              disabled={loading}
            >
              Delete Guest Type
            </Button>
            <Button type="submit" variant="contained" color="primary" disabled={loading}>
              {loading ? 'Updating...' : 'Update Guest Type'}
            </Button>
          </Box>
        </form>
      )}
    </Box>
  );
};

export default GuestTypeEdit;
