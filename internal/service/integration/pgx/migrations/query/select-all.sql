SELECT
    id,
    filename,
    executed_at
FROM migrations
ORDER BY executed_at DESC;
