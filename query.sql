
-- name: CreatePosting :exec
INSERT INTO grailed_postings (
    id, title, designer_names, department, sold, price, price_drop, image_url
) VALUES (
    ?,?,?,?,?,?,?,?
    )
