INSERT INTO wishlists(id, wisher_id, name, description, hidden)
VALUES (@id, @wisherId, @name, @description, @hidden)
ON CONFLICT (id) DO UPDATE
SET name = @name,
    description = @description,
    hidden = @hidden;