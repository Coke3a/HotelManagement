CREATE TABLE daily_booking_summary (
    summary_date DATE PRIMARY KEY,
    total_bookings INT NOT NULL,
    total_amount DECIMAL(10, 2) NOT NULL,
    pending_bookings INT NOT NULL,
    confirmed_bookings INT NOT NULL,
    checked_in_bookings INT NOT NULL,
    checked_out_bookings INT NOT NULL,
    canceled_bookings INT NOT NULL,
    completed_bookings INT NOT NULL,
    booking_ids TEXT NOT NULL, -- Store booking IDs as a comma-separated string
    status INT DEFAULT 0, -- 0: Unchecked, 1: Checked, 2: Confirmed
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
