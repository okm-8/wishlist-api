SELECT
    id,
    email,
    name,
    admin,
    password_hash
FROM users
WHERE id = ANY(@ids);
