export const PaymentStatus = {
  UNPAID: 1,
  PAID: 2,
  FAILED: 3,
  REFUNDED: 4
};

export const getPaymentStatusMessage = (status) => {
  switch (parseInt(status)) {
    case PaymentStatus.UNPAID:
      return 'Unpaid';
    case PaymentStatus.PAID:
      return 'Paid';
    case PaymentStatus.FAILED:
      return 'Failed';
    case PaymentStatus.REFUNDED:
      return 'Refunded';
    default:
      return 'Unknown';
  }
};

export const getPaymentStatusColor = (status) => {
  switch (parseInt(status)) {
    case PaymentStatus.UNPAID:
      return {
        backgroundColor: '#FFF4E5',
        textColor: '#663C00',
        chipColor: 'warning'
      };
    case PaymentStatus.PAID:
      return {
        backgroundColor: '#EDF7ED',
        textColor: '#1E4620',
        chipColor: 'success'
      };
    case PaymentStatus.FAILED:
      return {
        backgroundColor: '#FEEBEE',
        textColor: '#7F1D1D',
        chipColor: 'error'
      };
    case PaymentStatus.REFUNDED:
      return {
        backgroundColor: '#E5D7FD',
        textColor: '#4A1D96',
        chipColor: 'info'
      };
    default:
      return {
        backgroundColor: '#F5F5F5',
        textColor: '#666666',
        chipColor: 'default'
      };
  }
};

export const PaymentMethod = {
  NOT_SPECIFIED: 0,
  CREDIT_CARD: 1,
  DEBIT_CARD: 2,
  CASH: 3,
  BANK_TRANSFER: 4
};

export const getPaymentMethodMessage = (method) => {
  switch (parseInt(method)) {
    case PaymentMethod.NOT_SPECIFIED:
      return 'Not specified';
    case PaymentMethod.CREDIT_CARD:
      return 'Credit Card';
    case PaymentMethod.DEBIT_CARD:
      return 'Debit Card';
    case PaymentMethod.CASH:
      return 'Cash';
    case PaymentMethod.BANK_TRANSFER:
      return 'Bank Transfer';
    default:
      return `Unknown (${method})`;
  }
};
