import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { UserRoleEnum } from '../utils/userRoleEnum';
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
} from '@mui/material';
import { getUserRole } from '../utils/auth';

const GuestType = () => {
  const token = localStorage.getItem('token');
  const navigate = useNavigate();
  const [guestTypes, setGuestTypes] = useState([]);
  const [loading, setLoading] = useState(true);
  const [userRole, setUserRole] = useState('');
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [totalCount, setTotalCount] = useState(0);
  const [error, setError] = useState(null);

  useEffect(() => {
    setUserRole(getUserRole());
  }, []);

  const fetchGuestTypes = async (page, rowsPerPage) => {
    setLoading(true);
    try {
      const skip = page * rowsPerPage;
      const queryParams = new URLSearchParams({
        skip: skip.toString(),
        limit: rowsPerPage.toString(),
      });

      const response = await fetch(`http://localhost:8080/v1/customer-types/?${queryParams.toString()}`, {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
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
        if (data.data.meta && typeof data.data.meta.total === 'number') {
          setTotalCount(data.data.meta.total);
        } else {
          setTotalCount(data.data.customerTypes.length);
        }
      } else {
        setGuestTypes([]);
        setTotalCount(0);
      }
    } catch (error) {
      console.error('Error fetching guest types:', error);
      setError(error.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchGuestTypes(page, rowsPerPage);
  }, [page, rowsPerPage, navigate]);

  const handlePageChange = (event, newPage) => {
    setPage(newPage - 1);
  };

  const handleEditGuestType = (id) => {
    navigate(`/guest-type/edit/${id}`);
  };

  const handleAddGuestType = () => {
    navigate('/guest-type/add');
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
        {userRole === UserRoleEnum.ADMIN && (
          <Button
            variant="contained"
            className="bg-blue-600 text-white hover:bg-blue-700"
            onClick={handleAddGuestType}
          >
            Add Guest Type
          </Button>
        )}
        
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
              <option value={5}>5 per page</option>
              <option value={10}>10 per page</option>
              <option value={25}>25 per page</option>
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
      <TableContainer component={Paper}>
        <Table size="small" className="compact-table">
          <TableHead>
            <TableRow>
              <TableCell>ID</TableCell>
              <TableCell>Name</TableCell>
              <TableCell>Description</TableCell>
              <TableCell>Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {loading ? (
              <TableRow>
                <TableCell colSpan={4} align="center">
                  <Box display="flex" justifyContent="center" alignItems="center" height="50px">
                    <CircularProgress />
                  </Box>
                </TableCell>
              </TableRow>
            ) : !guestTypes || guestTypes.length === 0 ? (
              <TableRow>
                <TableCell colSpan={4} align="center">
                  <Typography variant="body1">No guest types found</Typography>
                </TableCell>
              </TableRow>
            ) : (
              guestTypes.map((guestType, index) => (
                <TableRow key={guestType.id} className={index % 2 === 0 ? "bg-gray-50" : "bg-white"}>
                  <TableCell>{guestType.id}</TableCell>
                  <TableCell>{guestType.name}</TableCell>
                  <TableCell>{guestType.description}</TableCell>
                  <TableCell>
                    <Button
                      variant="outlined"
                      color="primary"
                      size="small"
                      onClick={() => handleEditGuestType(guestType.id)}
                    >
                      Edit
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

export default GuestType;
