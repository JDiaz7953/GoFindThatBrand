package main

type Root struct {
	Results []Results `json:"results,omitempty"`
}

type Hits struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	DesignerNames string `json:"designer_names"`
	Department    string `json:"department"`
	Sold          bool   `json:"sold"`
	Price         int    `json:"price"`
	PriceDrop     []int  `json:"price_drop"`
	PostingUrl    string

	CoverPhoto struct {
		ImageUrl string `json:"image_url"`
	} `json:"cover_photo"`
	// I need those photos in those lisitng cover_photo : {image_url: <url>}
}

type Results struct {
	Hits    []Hits `json:"hits,omitempty"`
	Page    int    `json:"page"`
	NbPages int    `json:"nbPages"`
	NbHits  int    `json:"nbHits"`
}

// Have sqlite calculate the median price then search the db against items that fall below that median price
