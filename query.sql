
-- name: CreatePosting :exec
INSERT INTO grailed_postings_db (
    id, title, designer_names, department, category, sold, price, price_drop, image_url, posting_url, notified_at, flagged
) VALUES (
    ?,?,?,?,?,?,?,?,?,?,?,?
    )
ON CONFLICT(id) DO UPDATE SET
    title      = excluded.title,
    price      = excluded.price,
    price_drop = excluded.price_drop,
    sold       = excluded.sold,
    image_url  = excluded.image_url;

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

-- Name: InsertFilters :exec
INSERT INTO scrape_queue(
    id, designer_names, category
    )
    VALUES (?,?,?);


-- name: ScrapingQueue :one
SELECT id, designer_names, category FROM scrape_queue
ORDER BY last_scraped, id
LIMIT 1;

-- name: MarkJobCompleted :exec
UPDATE scrape_queue
SET last_scraped = datetime('now')
WHERE id = ?;
