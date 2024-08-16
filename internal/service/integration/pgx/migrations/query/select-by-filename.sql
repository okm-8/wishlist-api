SELECT
    id,
    executed_at
FROM migrations WHERE filename=@filename;
