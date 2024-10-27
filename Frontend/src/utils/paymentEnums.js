export const PaymentStatus = {
  PENDING: 1,
  COMPLETED: 2,
  FAILED: 3,
  REFUNDED: 4
};

export const getPaymentStatusMessage = (status) => {
  switch (status) {
    case PaymentStatus.PENDING:
      return 'Pending';
    case PaymentStatus.COMPLETED:
      return 'Completed';
    case PaymentStatus.FAILED:
      return 'Failed';
    case PaymentStatus.REFUNDED:
      return 'Refunded';
    default:
      return 'Unknown';
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
