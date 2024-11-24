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
  Pagination,
  Stack,
} from '@mui/material';
import { getUserRole } from '../utils/auth';
import { UserRoleEnum, UserRoleName } from '../utils/userRoleEnum';
import { handleTokenExpiration } from '../utils/api';

const User = () => {
  const navigate = useNavigate();
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [userRole, setUserRole] = useState('');
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [totalCount, setTotalCount] = useState(0);
  const [error, setError] = useState(null);

  useEffect(() => {
    setUserRole(getUserRole());
  }, []);

  const fetchUsers = async (page, rowsPerPage) => {
    setLoading(true);
    try {
      const skip = page * rowsPerPage;
      const queryParams = new URLSearchParams({
        skip: skip.toString(),
        limit: rowsPerPage.toString(),
      });
      const token = localStorage.getItem('token');
      const response = await fetch(`http://localhost:8080/v1/users/?${queryParams.toString()}`, {
        method: 'GET',
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
        throw new Error('Failed to fetch users');
      }

      const data = await response.json();
      if (data && data.data && data.data.users) {
        setUsers(data.data.users);
      } else {
        setUsers([]);
      }

      if (data && data.data && data.data.meta && typeof data.data.meta.total === 'number') {
        setTotalCount(data.data.meta.total);
      } else {
        setTotalCount(data.data?.users?.length || 0);
      }
    } catch (error) {
      console.error('Error fetching users:', error);
      setError(error.message);
      setUsers([]);
      setTotalCount(0);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchUsers(page, rowsPerPage);
  }, [page, rowsPerPage, navigate]);

  const handlePageChange = (event, newPage) => {
    setPage(newPage - 1);
  };

  const handleAddUser = () => {
    navigate('/user/add');
  };

  const handleEditUser = (id) => {
    navigate(`/user/edit/${id}`);
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
        {/* Left side - Add User button */}
        {userRole === UserRoleEnum.ADMIN && (
          <Button
            variant="contained"
            className="bg-blue-600 text-white hover:bg-blue-700"
            sx={{ mb: 1 }}
            onClick={handleAddUser}
          >
            Add User
          </Button>
        )}
        
        {/* Right side - Pagination controls */}
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
            page={page + 1} // Convert 0-based to 1-based
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
              <TableCell>Username</TableCell>
              <TableCell>Role</TableCell>
              <TableCell>Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {loading ? (
              <TableRow>
                <TableCell colSpan={4} align="center">
                  <CircularProgress size={24} />
                </TableCell>
              </TableRow>
            ) : users.length === 0 ? (
              <TableRow>
                <TableCell colSpan={4} align="center">
                  <Typography variant="body2">No users found</Typography>
                </TableCell>
              </TableRow>
            ) : (
              users.map((user) => (
                <TableRow key={user.id}>
                  <TableCell>{user.id}</TableCell>
                  <TableCell>{user.username}</TableCell>
                  <TableCell>{UserRoleName[user.role]}</TableCell>
                  <TableCell>
                    <Button
                      variant="outlined"
                      color="primary"
                      size="small"
                      onClick={() => handleEditUser(user.id)}
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

      {error && (
        <Typography color="error" align="center" style={{ padding: '1rem' }}>
          {error}
        </Typography>
      )}
    </div>
  );
};

export default User;
