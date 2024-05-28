package transfer

// ArticleCreation represents the data required to create a new article entry.
type ArticleCreation struct {
  Title    string `json:"title" binding:"required,max=256"`
  Slug     string
  ReadTime int
  Content  string `json:"content"`
}

// ArticleUpdate represents the data required to update an existing article entry.
type ArticleUpdate struct {
  Title    string `json:"title"`
  Slug     string
  ReadTime int
  Content  string `json:"content"`
}
