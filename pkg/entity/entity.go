package entity

//Item is the structure of a general shopping item
type Item struct {
	ID           string  `json:id`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	Rating       string  `json:"rating"`
	Summary      string  `json:"summary"`
	Manufacturer string  `json:"manufacturer"`
	Sources      Source  `json:"source"`
}

//Source is the structure of provider of an item
type Source struct {
	Name          string `json:"name"`
	Author        string `json:"author"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	URL           string `json:"url"`
	DatePublished string `json:"date"`
}

//Items is a list if items
type Items []Item

// APIError is to record errors
type APIError struct {
	ErrorCode    int
	ErrorMessage string
}

func (I Items) append(itm Item) {
	I = append(I, itm)
}

//NewItems Return list of Items
func NewItems() Items {
	source := Source{
		Name:          "N1",
		Author:        "A1",
		Title:         "T1",
		Content:       "C1",
		URL:           "https://test",
		DatePublished: "09-08-2019",
	}

	item1 := Item{
		ID:      "10001",
		Name:    "Book",
		Price:   40,
		Summary: "My Book",
		Sources: source,
	}

	item2 := Item{
		ID:      "10002",
		Name:    "Book2",
		Price:   20,
		Summary: "My Book 2",
		Sources: source,
	}

	ItemList := Items{
		item1,
		item2,
	}

	return ItemList
}
