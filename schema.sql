CREATE TABLE grailed_postings (
	"id"	INTEGER NOT NULL,
	"title"	TEXT NOT NULL,
	"designer_names"	TEXT NOT NULL,
	"department"	TEXT NOT NULL,
	"sold"	TEXT NOT NULL,
	"price"	INTEGER NOT NULL,
	"price_drop"	TEXT,
	"image_url"	TEXT NOT NULL,
	PRIMARY KEY("id")
);