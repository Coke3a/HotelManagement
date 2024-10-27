import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { Box, Typography, Button, CircularProgress } from '@mui/material';
import { format, differenceInDays } from 'date-fns';

const ReceiptPrint = () => {
    const { id } = useParams();
    const [booking, setBooking] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchBookingData = async () => {
            try {
                const response = await fetch(`http://localhost:8080/v1/booking/${id}/details`, {
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    credentials: 'include'
                });

                if (!response.ok) {
                    throw new Error('Failed to fetch booking data');
                }

                const data = await response.json();
                setBooking(data.data);
            } catch (error) {
                console.error('Error fetching booking data:', error);
                setError(error.message);
            } finally {
                setLoading(false);
            }
        };

        fetchBookingData();
    }, [id]);

    const handlePrint = () => {
        window.print();
    };

    const calculateStayDuration = (checkInDate, checkOutDate) => {
        const start = new Date(checkInDate);
        const end = new Date(checkOutDate);
        return differenceInDays(end, start);
    };

    const calculatePricePerNight = (totalPrice, checkInDate, checkOutDate) => {
        const nights = calculateStayDuration(checkInDate, checkOutDate);
        return (totalPrice / nights).toFixed(2);
    };

    const formatReceiptContent = () => {
        if (!booking) return null;

        return (
            <div className="receipt-page">
                <div className="receipt-header">
                    <div className="receipt-number">
                        <p>เลขที่ (Receipt NO.) {booking.booking_id}</p>
                        <p>วันที่ (Date) {format(new Date(), 'dd/MM/yyyy')}</p>
                    </div>
                    <div className="receipt-logo-container">
                        <img src="/images/hotel-logo.png" alt="Hotel Logo" className="hotel-logo" />
                    </div>
                </div>
                <h1 className="receipt-title">ใบเสร็จรับเงิน<br />Receipt</h1>
                <div className="hotel-title">
                    <h2>โรงแรม มาย เมมโมรี่ เมืองเลย</h2>
                </div>
                <div className="hotel-info">
                    <p>เลขที่ 454 หมู่ 2 ตำบลนาอาน อำเภอเมืองเลย จังหวัดเลย 42000</p>
                    <p>โทรศัพท์. 065-6249290, 042-030290, 081-9190968,</p>
                    <p>E-mail: Mymemory.loei@gmail.com</p>
                    <p>Line@924hxdqn</p>
                    <p>ทะเบียนเลขที่ : 0425563000601</p>
                </div>
                <div className="customer-info">
                    <p><span>ชื่อ (Customer Name):</span> {booking.customer_firstname} {booking.customer_surname}</p>
                    <p><span>ที่อยู่ (Address):</span> {booking.customer_address || '-'}</p>
                    <p><span>เลขประจำตัวผู้เสียภาษี (Tax ID):</span> {booking.customer_tax_id || '-'}</p>
                </div>
                <table className="receipt-details">
                    <thead>
                        <tr>
                            <th>ลำดับที่</th>
                            <th>รายการ</th>
                            <th>จำนวนเงิน</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr>
                            <td>1</td>
                            <td>
                                ผู้พัก {booking.customer_firstname} {booking.customer_surname}<br />
                                ค่าห้องพัก ({booking.room_number}) {booking.room_type_name}<br />
                                [ จำนวน {calculateStayDuration(booking.check_in_date, booking.check_out_date)} คืน คืนละ {calculatePricePerNight(booking.booking_price, booking.check_in_date, booking.check_out_date)} บาท ]<br />
                                Check In: วันที่ {format(new Date(booking.check_in_date), 'dd/MM/yyyy')}<br />
                                Check Out: วันที่ {format(new Date(booking.check_out_date), 'dd/MM/yyyy')}
                            </td>
                            <td>{booking.booking_price}</td>
                        </tr>
                    </tbody>
                    <tfoot>
                        <tr>
                            <td colSpan="2">
                                จำนวนเงินรวม / Total 
                                <span className="total-amount">({numberToThaiWords(booking.booking_price)})</span>
                            </td> 
                            <td>{booking.booking_price}</td>
                        </tr>
                    </tfoot>
                </table>
                <div className="receipt-footer">
                    <p>ลงชื่อ ______________________________</p>
                    <p>วันที่ {format(new Date(), 'dd/MM/yyyy')}</p>
                </div>
            </div>
        );
    };

    const numberToThaiWords = (num) => {
        const units = ['', 'หนึ่ง', 'สอง', 'สาม', 'สี่', 'ห้า', 'หก', 'เจ็ด', 'แปด', 'เก้า'];
        const tens = ['', 'สิบ', 'ยี่สิบ', 'สามสิบ', 'สี่สิบ', 'ห้าสิบ', 'หกสิบ', 'เจ็ดสิบ', 'แปดสิบ', 'เก้าสิบ'];
        const scales = ['', 'สิบ', 'ร้อย', 'พัน', 'หมื่น', 'แสน', 'ล้าน'];

        if (num === 0) return 'ศูนย์';

        let words = '';
        let scale = 0;

        while (num > 0) {
            const n = num % 10;
            if (n !== 0) {
                words = (scale === 1 && n === 1 ? '' : units[n]) + scales[scale] + words;
            }
            num = Math.floor(num / 10);
            scale++;
        }

        return words;
    };

    return (
        <Box>
            <Box className="no-print" style={{ marginBottom: '20px' }}>
                <Button variant="contained" color="primary" onClick={handlePrint}>
                    Print Receipt
                </Button>
            </Box>
            <Box className="receipt-content">
                {loading ? (
                    <CircularProgress />
                ) : error ? (
                    <Typography color="error">{error}</Typography>
                ) : !booking ? (
                    <Typography>No booking data found</Typography>
                ) : (
                    formatReceiptContent()
                )}
            </Box>
        </Box>
    );
};

export default ReceiptPrint;
