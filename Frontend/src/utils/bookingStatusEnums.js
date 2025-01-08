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
  switch (status) {
    case BookingStatus.UNCHECK_IN:
      return {
        backgroundColor: '#FFF4E5',
        textColor: '#663C00'
      };
    case BookingStatus.CHECKED_IN:
      return {
        backgroundColor: '#E8F4FD',
        textColor: '#0A4FA6'
      };
    case BookingStatus.CHECKED_OUT:
      return {
        backgroundColor: '#EDF7ED',
        textColor: '#1E4620'
      };
    case BookingStatus.CANCELED:
      return {
        backgroundColor: '#FEEBEE',
        textColor: '#7F1D1D'
      };
    case BookingStatus.COMPLETED:
      return {
        backgroundColor: '#E5D7FD',
        textColor: '#4A1D96'
      };
    default:
      return {
        backgroundColor: '#F5F5F5',
        textColor: '#666666'
      };
  }
};

