SELECT
    id,
    name,
    email,
    admin,
    password_hash
FROM users
LIMIT @limit
OFFSET @offset;
