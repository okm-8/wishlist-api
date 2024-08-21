SELECT
    id,
    name,
    description,
    hidden,
    fulfilled,
    assignee_id,
    assignee_email,
    assignee_name,
    wishlist_id,
    wishlist_name,
    wishlist_description,
    wishlist_hidden,
    wisher_id,
    wisher_email,
    wisher_name
FROM wishes_view
WHERE id = @id
