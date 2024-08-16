SELECT
    id,
    created_at,
    expire_at
FROM signup_sessions_view
WHERE email = @email;