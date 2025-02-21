-- This is the SQL script that will be used to initialize the database schema.
-- We will evaluate you based on how well you design your database.
-- 1. How you design the tables.
-- 2. How you choose the data types and keys.
-- 3. How you name the fields.
-- In this assignment we will use PostgreSQL as the database.

CREATE TABLE estates (
    id UUID PRIMARY KEY,
    length INTEGER NOT NULL,
    width INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE trees (
    id UUID PRIMARY KEY,
    estate_id UUID NOT NULL,
    row INTEGER NOT NULL,
    col INTEGER NOT NULL,
    height INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (estate_id, row, col)
);

CREATE TABLE drone_routes (
    route INTEGER NOT NULL,
    estate_id UUID NOT NULL,
    row INTEGER NOT NULL,
    col INTEGER NOT NULL,
    altitude INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (estate_id, row, col)
);


-- Leveraging a materialized view to precompute aggregation values for improved read performance,  
-- assuming read operations are more frequent than writes.  
-- The materialized view is refreshed via a trigger on the trees table.  
-- However, if eventual consistency is acceptable, a periodic refresh is preferred,  
-- as triggers can negatively impact write performance.
CREATE MATERIALIZED VIEW estate_stats_mv AS
SELECT 
    estate_id,
    COUNT(height) AS tree_count,
    MAX(height) AS max_height,
    MIN(height) AS min_height,
    PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY height) AS median_height
FROM trees
GROUP BY estate_id;

CREATE UNIQUE INDEX estate_stats_mv_idx ON estate_stats_mv (estate_id);

CREATE OR REPLACE FUNCTION refresh_estate_stats_mv()
RETURNS TRIGGER AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY estate_stats_mv;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_refresh_mv
AFTER INSERT OR UPDATE OR DELETE ON trees
FOR EACH STATEMENT
EXECUTE FUNCTION refresh_estate_stats_mv();

