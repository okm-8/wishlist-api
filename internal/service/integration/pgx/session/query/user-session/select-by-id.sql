SELECT
    user_id,
    user_email,
    user_name,
    user_admin,
    user_password_hash,
    created_at,
    expire_at
FROM user_sessions_view
WHERE id = @id;
