CREATE VIEW booking_customer_payment AS
SELECT
    b.id AS booking_id,
    b.customer_id,
    b.total_amount AS booking_price,
    b.status AS booking_status,
    b.check_in_date,
    b.check_out_date,
    b.created_at AS booking_created_at,
    b.updated_at AS booking_updated_at,
    b.room_id,
    r.room_number,
    r.type_id AS room_type_id,
    rt.name AS room_type_name,
    r.floor,
    b.rate_prices_id,
    c.firstname AS customer_firstname,
    c.surname AS customer_surname,
    c.identity_number AS customer_identity_number,
    c.address AS customer_address,
    p.id AS payment_id,
    p.status AS payment_status,
    p.updated_at AS payment_update_date
FROM
    bookings b
    JOIN customers c ON b.customer_id = c.id
    LEFT JOIN payments p ON b.id = p.booking_id
    JOIN rooms r ON b.room_id = r.id
    JOIN room_types rt ON r.type_id = rt.id;