CREATE VIEW booking_room_rate_price AS
SELECT
    b.id AS booking_id,
    b.status AS booking_status,
    b.check_in_date,
    b.check_out_date,
    r.id AS room_id,
    r.room_number,
    r.status AS room_status,
    rp.id AS rate_price_id,
    rp.price_per_night,
    rp.name AS rate_price_name,
    rt.id AS room_type_id,
    rt.name AS room_type_name
FROM
    bookings b
    JOIN rate_prices rp ON b.rate_prices_id = rp.id
    JOIN room_types rt ON rp.room_type_id = rt.id
    JOIN rooms r ON r.type_id = rt.id;