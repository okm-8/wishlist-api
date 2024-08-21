UPDATE users SET email = @email, name = @name, admin = @admin, password_hash = @passwordHash
WHERE id = @id;
