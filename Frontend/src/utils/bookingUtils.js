export const countBookings = (bookingString) => {
  if (!bookingString || bookingString === '') {
    return 0;
  }
  return bookingString.split(';').length;
};

export const getStatusText = (status) => {
  const statusMap = {
    0: 'Unchecked',
    1: 'Checked',
    2: 'Confirmed'
  };
  return statusMap[status] || 'Unknown';
};

export const getStatusColor = (status) => {
  const colorMap = {
    0: 'default',
    1: 'processing',
    2: 'success'
  };
  return colorMap[status] || 'default';
}; 