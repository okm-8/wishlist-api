SELECT DISTINCT
    wisher_id,
    wisher_email,
    wisher_name
FROM wishes_view
WHERE wishlist_hidden = FALSE