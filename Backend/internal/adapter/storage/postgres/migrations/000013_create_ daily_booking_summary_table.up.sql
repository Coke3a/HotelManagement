CREATE TABLE daily_booking_summary (
    summary_date DATE PRIMARY KEY,
    created_bookings TEXT NOT NULL, -- Store booking IDs as a comma-separated string
    completed_bookings TEXT NOT NULL, -- Store booking IDs as a comma-separated string
    canceled_bookings TEXT NOT NULL, -- Store booking IDs as a comma-separated string
    total_amount DECIMAL(10, 2) NOT NULL,
    status INT DEFAULT 0, -- 0: Unchecked, 1: Checked, 2: Confirmed
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
