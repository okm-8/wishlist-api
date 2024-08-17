SELECT
    id,
    name,
    description,
    hidden,
    fulfilled,
    created_at,
    assignee_email,
    assignee_name,
    wishlist_id,
    wishlist_name,
    wishlist_description,
    wishlist_hidden,
    wishlist_created_at,
    wisher_id,
    wisher_email,
    wisher_name
FROM wishes_view
WHERE assignee_id = @assigneeId