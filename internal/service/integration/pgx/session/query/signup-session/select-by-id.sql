SELECT
    email,
    created_at,
    expire_at
FROM signup_sessions_view
WHERE id = @id;