export const BookingStatus = {
  UNCHECK_IN: 1,
  CHECKED_IN: 2,
  CHECKED_OUT: 3,
  CANCELED: 4,
  COMPLETED: 5
};

export const getBookingStatusMessage = (status) => {
  switch (status) {
    case BookingStatus.UNCHECK_IN:
      return 'Uncheck-in';
    case BookingStatus.CHECKED_IN:
      return 'Check-in';
    case BookingStatus.CHECKED_OUT:
      return 'Check-out';
    case BookingStatus.CANCELED:
      return 'Cancelled';
    case BookingStatus.COMPLETED:
      return 'Completed';
    default:
      return 'Unknown';
  }
};

export const getBookingStatusColor = (status) => {
  switch (parseInt(status)) {
    case BookingStatus.UNCHECK_IN:
      return {
        chipColor: 'warning'
      };
    case BookingStatus.CHECKED_IN:
      return {
        chipColor: 'info'
      };
    case BookingStatus.CHECKED_OUT:
      return {
        chipColor: 'default'
      };
    case BookingStatus.CANCELED:
      return {
        chipColor: 'error'
      };
    case BookingStatus.COMPLETED:
      return {
        chipColor: 'success'
      };
    default:
      return {
        chipColor: 'default'
      };
  }
};

