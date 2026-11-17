
-- name: CreatePosting :exec
INSERT INTO grailed_postings_db (
    id, title, designer_names, department, sold, price, price_drop, image_url, posting_url
) VALUES (
    ?,?,?,?,?,?,?,?,?
    )
