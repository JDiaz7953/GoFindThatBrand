CREATE TABLE "grailed_postings_db" (
	"id"	INTEGER NOT NULL,
	"title"	TEXT NOT NULL,
	"designer_names"	TEXT NOT NULL,
	"department"	TEXT NOT NULL,
	"category"	TEXT NOT NULL,
	"sold"	TEXT NOT NULL,
	"price"	INTEGER NOT NULL,
	"price_drop"	TEXT,
	"image_url"	TEXT NOT NULL,
	"posting_url"	TEXT NOT NULL,
	"notified_at"	TEXT,
    "flagged"	TEXT NOT NULL,
	PRIMARY KEY("id")
);