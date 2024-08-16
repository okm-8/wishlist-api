SELECT
    email,
    name,
    admin,
    password_hash
FROM users
WHERE id = @id;
