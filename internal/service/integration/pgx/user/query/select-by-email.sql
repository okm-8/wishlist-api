SELECT
    id,
    name,
    admin,
    password_hash
FROM users
WHERE email = @email;
