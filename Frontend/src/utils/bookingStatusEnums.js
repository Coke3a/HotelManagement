export const BookingStatus = {
  PENDING: 1,
  CONFIRMED: 2,
  CHECKED_IN: 3,
  CHECKED_OUT: 4,
  CANCELED: 5,
  COMPLETED: 6,
};

export const getBookingStatusMessage = (status) => {
  switch (status) {
    case BookingStatus.PENDING:
      return 'Pending';
    case BookingStatus.CONFIRMED:
      return 'Confirmed';
    case BookingStatus.CHECKED_IN:
      return 'Checked In';
    case BookingStatus.CHECKED_OUT:
      return 'Checked Out';
    case BookingStatus.CANCELED:
      return 'Canceled';
    case BookingStatus.COMPLETED:
      return 'Completed';
    default:
      return 'Unknown';
  }
};

