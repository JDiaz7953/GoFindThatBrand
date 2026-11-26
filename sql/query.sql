
-- name: CreatePosting :exec
INSERT INTO grailed_postings_db (
    id, title, designer_names, department, category, category_path, sold, price, price_drop, image_url, posting_url, notified_at, flagged
) VALUES (
    ?,?,?,?,?,?,?,?,?,?,?,?,?
    )
ON CONFLICT(id) DO UPDATE SET
    title      = excluded.title,
    price      = excluded.price,
    price_drop = excluded.price_drop,
    sold       = excluded.sold,
    image_url  = excluded.image_url;

-- name: CalculateMedianPriceOfBrand :one
SELECT AVG(price)
FROM (SELECT price, designer_names, category_path
    FROM grailed_postings_db as p
    WHERE p.designer_names = sqlc.arg(brand) AND p.category_path = sqlc.arg(category)
    ORDER BY p.price
    LIMIT 2 - ( SELECT COUNT(c1.price) FROM grailed_postings_db as c1
    WHERE c1.designer_names = sqlc.arg(brand) AND c1.category_path = sqlc.arg(category)) % 2
    OFFSET (SELECT (COUNT(c2.price) - 1) / 2
        FROM grailed_postings_db as c2  WHERE c2.designer_names = sqlc.arg(brand) AND c2.category_path = sqlc.arg(category))
        );

-- name: FlagDeals :exec
UPDATE grailed_postings_db
SET flagged = CASE
    WHEN price BETWEEN CAST(sqlc.arg(low) AS REAL) AND CAST(sqlc.arg(high) AS REAL)
    THEN 'true' ELSE 'false' END
WHERE designer_names = sqlc.arg(brand) AND category_path = sqlc.arg(category);

-- name: InsertFilters :exec
INSERT INTO scrape_queue(
    designer_names, category
    )
    VALUES (?,?);

-- name: GetFlaggedItem :many
SELECT title, posting_url, price, image_url
FROM grailed_postings_db
WHERE category_path = sqlc.arg(category) 
AND designer_names = sqlc.arg(brand)
AND flagged = true
ORDER BY price;

-- name: ScrapingQueue :one
SELECT id, designer_names, category FROM scrape_queue
ORDER BY last_scraped, id
LIMIT 1;

-- name: MarkJobCompleted :exec
UPDATE scrape_queue
SET last_scraped = datetime('now')
WHERE id = ?;
