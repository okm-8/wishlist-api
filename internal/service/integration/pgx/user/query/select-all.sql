SELECT
    id,
    email,
    name,
    admin,
    password_hash
FROM users
LIMIT @limit
OFFSET @offset;
