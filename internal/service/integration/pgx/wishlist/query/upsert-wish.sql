INSERT INTO wishes(id, wishlist_id, name, description, hidden, fulfilled, assignee_id, created_at)
VALUES (
    @id,
    @wishlistId,
    @name,
    @description,
    @hidden,
    @fulfilled,
    @assigneeId,
    @createdAt
)
ON CONFLICT (id) DO UPDATE SET
    wishlist_id = @wishlistId,
    name = @name,
    description = @description,
    hidden = @hidden,
    fulfilled = @fulfilled,
    assignee_id = @assigneeId
;
