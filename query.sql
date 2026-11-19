
-- name: CreatePosting :exec
INSERT INTO grailed_postings_db (
    id, title, designer_names, department, category, sold, price, price_drop, image_url, posting_url, notified_at, flagged
) VALUES (
    ?,?,?,?,?,?,?,?,?,?,?,?
    );


-- name: CalculateMedianPriceOfBrand :one
SELECT AVG(price)
FROM (SELECT price
    FROM grailed_postings_db
    ORDER BY price
    LIMIT 2 - ( SELECT COUNT(price) FROM grailed_postings_db) % 2
    OFFSET (SELECT (COUNT(price) - 1) / 2
        FROM grailed_postings_db));


-- name: FlagDeals :exec
UPDATE grailed_postings_db
SET flagged = 'true'
WHERE price BETWEEN CAST(sqlc.arg(low) AS REAL) AND CAST(sqlc.arg(high) AS REAL);
