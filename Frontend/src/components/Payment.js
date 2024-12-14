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
  Chip,
  IconButton,
} from '@mui/material';
import { getPaymentStatusMessage, getPaymentMethodMessage, PaymentStatus, PaymentMethod } from '../utils/paymentEnums';
import { getUserRole } from '../utils/auth';
import EditIcon from '@mui/icons-material/Edit';
import { handleTokenExpiration } from '../utils/api';

const Payment = () => {
  const token = localStorage.getItem('token');
  const navigate = useNavigate();
  const [payments, setPayments] = useState([]);
  const [loading, setLoading] = useState(true);
  const [userRole, setUserRole] = useState('');
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [totalCount, setTotalCount] = useState(0);
  const [error, setError] = useState(null);

  const fetchPayments = async (page, rowsPerPage) => {
    setLoading(true);
    try {
      const skip = page * rowsPerPage;
      const queryParams = new URLSearchParams({
        skip: skip.toString(),
        limit: rowsPerPage.toString(),
      });

      const response = await fetch(`http://localhost:8080/v1/payments/?${queryParams.toString()}`, {
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
        throw new Error('Failed to fetch payments');
      }

      const data = await response.json();
      if (data && data.data && data.data.payments) {
        setPayments(data.data.payments);
        if (data.data.meta && typeof data.data.meta.total === 'number') {
          setTotalCount(data.data.meta.total);
        } else {
          setTotalCount(data.data.payments.length);
        }
      } else {
        setPayments([]);
        setTotalCount(0);
      }
    } catch (error) {
      console.error('Error fetching payments:', error);
      setError(error.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPayments(page, rowsPerPage);
  }, [page, rowsPerPage, navigate]);

  useEffect(() => {
    setUserRole(getUserRole());
  }, []);

  const handlePageChange = (event, newPage) => {
    setPage(newPage - 1);
  };

  const renderPaymentStatus = (status) => {
    const statusMessage = getPaymentStatusMessage(status);
    let statusColor;

    switch (status) {
      case PaymentStatus.UNPAID:
        statusColor = 'warning';
        break;
      case PaymentStatus.PAID:
        statusColor = 'success';
        break;
      default:
        statusColor = 'default';
    }

    return (
      <Chip 
        label={statusMessage} 
        color={statusColor} 
        size="small"
      />
    );
  };

  const renderPaymentMethod = (method) => {
    const methodMessage = getPaymentMethodMessage(method);
    let methodColor;

    switch (method) {
      case PaymentMethod.CREDIT_CARD:
        methodColor = 'primary';
        break;
      case PaymentMethod.DEBIT_CARD:
        methodColor = 'info';
        break;
      case PaymentMethod.CASH:
        methodColor = 'success';
        break;
      case PaymentMethod.BANK_TRANSFER:
        methodColor = 'secondary';
        break;
      case PaymentMethod.NOT_SPECIFIED:
      default:
        methodColor = 'default';
    }

    return (
      <Chip 
        label={methodMessage} 
        color={methodColor} 
        size="small"
      />
    );
  };

  const handleEditPayment = (id) => {
    navigate(`/payment/edit/${id}`);
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
              <TableCell>Amount</TableCell>
              <TableCell>Payment Method</TableCell>
              <TableCell>Status</TableCell>
              <TableCell>Payment Date</TableCell>
              <TableCell>Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {loading ? (
              <TableRow>
                <TableCell colSpan={6} align="center">
                  <Box display="flex" justifyContent="center" alignItems="center" height="50px">
                    <CircularProgress />
                  </Box>
                </TableCell>
              </TableRow>
            ) : payments.length === 0 ? (
              <TableRow>
                <TableCell colSpan={6} align="center">
                  <Typography variant="body1">No payments found</Typography>
                </TableCell>
              </TableRow>
            ) : (
              payments.map((payment, index) => (
                <TableRow key={payment.id} className={index % 2 === 0 ? "bg-gray-50" : "bg-white"}>
                  <TableCell>{payment.id}</TableCell>
                  <TableCell>{payment.amount}</TableCell>
                  <TableCell>{renderPaymentMethod(payment.payment_method)}</TableCell>
                  <TableCell>{renderPaymentStatus(payment.status)}</TableCell>
                  <TableCell>{new Date(payment.payment_date).toLocaleString()}</TableCell>
                  <TableCell>
                    <Button
                      variant="outlined"
                      color="primary"
                      size="small"
                      onClick={() => handleEditPayment(payment.id)}
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

export default Payment;
