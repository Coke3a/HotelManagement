import React from 'react';
import { Modal, Paper, Typography, Button } from '@mui/material';
import { format } from 'date-fns';

const ReceiptModal = ({ open, onClose, booking }) => {
  const receiptContent = `
  ,,,,,,,,,
,,,,,,,เลขที่ (Receipt NO.),,${booking.booking_id}
,,,,,,,,วันที่ (Date),${format(new Date(), 'yyyy-MM-dd')}
,,,,,,,,,
"ใบเสร็จรับเงิน
Receipt",,,,,,,,,
โรงแรม มาย เมมโมรี่ เมืองเลย,,,,,,,,,
เลขที่  454 หมู่ 2 ตำบลนาอาน อำเภอเมืองเลย จังหวัดเลย 42000,,,,,,,,,
"โทรศัพท์. 065-6249290, 042-030290, 081-9190968, 
 E-mail: Mymemory.loei@gmail.com
 Line@924hxdqn
ทะเบียนเลขที่ : 0425563000601 ",,,,,,,,,
ชื่อ( Customer Name ),${booking.customer_firstname} ${booking.customer_surname},,,,,,,,
ที่อยู่ ( Address ),,,,,,,,,
"เลขประจำตัวผู้เสียภาษี
 ( Tax ID )",,,,,,,,,
,,,,,,,,,
ลำดับที่,รายการ ,,,,,,,,จ��นวนเงิน
1,ผู้พัก ${booking.customer_firstname} ${booking.customer_surname},,,,,,,, ${booking.booking_price}
,ค่าห้องพัก,,,,,,,,
,[ จำนวน ${calculateStayDuration(booking.check_in_date, booking.check_out_date)} คืน คืนละ ${calculatePricePerNight(booking.booking_price, booking.check_in_date, booking.check_out_date)} บาท ],,,,,,,,
,"Check In:
",,วันที่ ${format(new Date(booking.check_in_date), 'yyyy-MM-dd')},,,,,,
,Check Out: ,,วันที่ ${format(new Date(booking.check_out_date), 'yyyy-MM-dd')},,,,,,
,,,,,,,,,
,,,,,,,,,
,,,,,,,,"จำนวนเงินรวม
Total",, ${booking.booking_price}
,,,,,,,,,
,,,,,,ลงชื่อ,______________________________,,
,,,,,,วันที่,${format(new Date(), 'yyyy-MM-dd')},,
  `;

  const handlePrint = () => {
    const printWindow = window.open('', '_blank');
    printWindow.document.write('<pre>' + receiptContent + '</pre>');
    printWindow.document.close();
    printWindow.print();
  };

  return (
    <Modal open={open} onClose={onClose}>
      <Paper style={{ padding: '20px', maxWidth: '800px', margin: '20px auto' }}>
        <Typography variant="h5" gutterBottom>Receipt</Typography>
        <pre style={{ whiteSpace: 'pre-wrap', fontFamily: 'monospace', fontSize: '12px' }}>
          {receiptContent}
        </pre>
        <Button variant="contained" color="primary" onClick={handlePrint}>
          Print Receipt
        </Button>
      </Paper>
    </Modal>
  );
};

export default ReceiptModal;

function calculateStayDuration(checkInDate, checkOutDate) {
  const start = new Date(checkInDate);
  const end = new Date(checkOutDate);
  const diffTime = Math.abs(end - start);
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
  return diffDays;
}

function calculatePricePerNight(totalPrice, checkInDate, checkOutDate) {
  const nights = calculateStayDuration(checkInDate, checkOutDate);
  return (totalPrice / nights).toFixed(2);
}
