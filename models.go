package main

type Root struct {
	Results []Results `json:"results,omitempty"`
}

type Hits struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	DesignerNames string `json:"designer_names"`
	Department    string `json:"department"`
	Catagory	  string `json:"category"`
	Sold          bool   `json:"sold"`
	Price         int    `json:"price"`
	PriceDrop     []int  `json:"price_drop"`
	PostingUrl    string
	Notified_at   string

	CoverPhoto struct {
		ImageUrl string `json:"image_url"`
	} `json:"cover_photo"`
	
}

type Results struct {
	Hits    []Hits `json:"hits,omitempty"`
	Page    int    `json:"page"`
	NbPages int    `json:"nbPages"`
	NbHits  int    `json:"nbHits"`
}
