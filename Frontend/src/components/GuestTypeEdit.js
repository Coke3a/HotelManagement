import React, { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import {
  Box,
  Typography,
  TextField,
  Button,
  CircularProgress,
  Paper,
  Grid,
} from '@mui/material';

const GuestTypeEdit = () => {
  const token = localStorage.getItem('token');
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
            'Authorization': `Bearer ${token}`,
          },
        });

        if (response.status === 401) {
          handleTokenExpiration(new Error("access token has expired"), navigate);
          return;
        }

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
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify(guestType),
      });

      if (response.status === 401) {
        handleTokenExpiration(new Error("access token has expired"), navigate);
        return;
      }

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
          headers: {
            'Authorization': `Bearer ${token}`,
          },
        });

        if (response.status === 401) {
          handleTokenExpiration(new Error("access token has expired"), navigate);
          return;
        }

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
    <Box sx={{ maxWidth: '800px', margin: '0 auto', padding: '20px' }}>
      <Paper elevation={3} sx={{ padding: '24px', backgroundColor: '#f8f9fa' }}>
        <Typography variant="h5" gutterBottom className="form-title">
          Edit Guest Type
        </Typography>
        {loading ? (
          <Box display="flex" justifyContent="center" alignItems="center" height="300px">
            <CircularProgress />
          </Box>
        ) : (
          <form onSubmit={handleSubmit} className="form">
            <Grid container spacing={3}>
              <Grid item xs={12} md={6}>
                <TextField
                  fullWidth
                  label="Name"
                  name="name"
                  value={guestType.name}
                  onChange={handleChange}
                  required
                  className="form-input"
                  size="small"
                />
              </Grid>
              <Grid item xs={12}>
                <TextField
                  fullWidth
                  label="Description"
                  name="description"
                  value={guestType.description}
                  onChange={handleChange}
                  multiline
                  rows={3}
                  className="form-input"
                  size="small"
                />
              </Grid>
            </Grid>
            <Box display="flex" justifyContent="flex-end" mt={3} gap={2}>
              <Button
                variant="outlined"
                color="error"
                onClick={handleDelete}
                disabled={loading}
                size="medium"
              >
                Delete Guest Type
              </Button>
              <Button
                type="submit"
                variant="contained"
                color="primary"
                disabled={loading}
                size="medium"
              >
                {loading ? 'Updating...' : 'Update Guest Type'}
              </Button>
            </Box>
          </form>
        )}
      </Paper>
    </Box>
  );
};

export default GuestTypeEdit;
