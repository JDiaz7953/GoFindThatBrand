package models

type Root struct {
	Results []Results `json:"results,omitempty"`
}

type Hits struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	DesignerNames string `json:"designer_names"`
	Department    string `json:"department"`
	Catagory      string `json:"category"`
	CatagoryPath  string `json:"category_path"`
	Sold          bool   `json:"sold"`
	Price         int    `json:"price"`
	PriceDrop     []int  `json:"price_drop"`
	PostingUrl    string
	Notified_at   string
	Flagged       bool

	CoverPhoto struct {
		ImageUrl string `json:"image_url"`
	} `json:"cover_photo"`
}

type Results struct {
	Hits    []Hits `json:"hits,omitempty"`
	Page    int    `json:"page"`
	NbPages int    `json:"nbPages"`
	NbHits  int    `json:"nbHits"`

	FacetStats struct {
		Price_i struct {
			Max int `json:"max"`
		} `json:"price_i"`
	} `json:"facets_stats"`

	Facets struct {
		DesignerNames map[string]int `json:"designers.name"`
		CategoryPath  map[string]int `json:"category_path"`
	} `json:"facets"`
}

type DiscordMessage struct {
	Url		   string
	Title	   string
	PriceDesc  int
	Image	   *MessageEmbedImage 
	

}

type MessageEmbedImage struct {
	URL       string `json:"url"`
	ProxyURL  string `json:"proxy_url,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
}