INSERT INTO migrations(filename)
VALUES (@filename)
ON CONFLICT (filename) DO UPDATE
SET executed_at = NOW();
