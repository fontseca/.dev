package transfer

import (
  "time"
)

// ArticleCreation represents the data required to create a new article entry.
type ArticleCreation struct {
  Title     string `json:"title" binding:"required,max=256"`
  Slug      string
  ReadTime  int
  Content   string `json:"content"`
  DraftedAt time.Time
}

// ArticleUpdate represents the data required to update an existing article entry.
type ArticleUpdate struct {
  Title       string `json:"title"`
  Slug        string
  ReadTime    int
  Content     string `json:"content"`
  DraftedAt   time.Time
  PinnedAt    *time.Time
  ArchivedAt  *time.Time
  PublishedAt *time.Time
  ModifiedAt  *time.Time
}