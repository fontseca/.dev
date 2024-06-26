package repository

import (
  "context"
  "database/sql"
  "fontseca.dev/model"
  "fontseca.dev/problem"
  "fontseca.dev/transfer"
  "log/slog"
  "net/http"
  "time"
)

// TagsRepository is a low level API that provides methods for interacting
// with tags in the database.
type TagsRepository interface {
  // Add adds a new tag.
  Add(ctx context.Context, creation *transfer.TagCreation) error

  // Get retrieves all the tags.
  Get(ctx context.Context) (tags []*model.Tag, err error)

  // Update updates an existing tag.
  Update(ctx context.Context, id string, update *transfer.TagUpdate) error

  // Remove removes a tag and detaches it from any article that uses it.
  Remove(ctx context.Context, id string) error
}

type tagsRepository struct {
  db *sql.DB
}

func NewTagsRepository(db *sql.DB) TagsRepository {
  return &tagsRepository{db}
}

func (r *tagsRepository) Add(ctx context.Context, creation *transfer.TagCreation) error {
  tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
  if nil != err {
    return err
  }

  defer tx.Rollback()

  addTagQuery := `
  INSERT INTO "tag" ("id", "name")
               VALUES (@id, @name);`

  ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
  defer cancel()

  result, err := tx.ExecContext(ctx, addTagQuery,
    sql.Named("id", creation.ID),
    sql.Named("name", creation.Name),
  )

  if nil != err {
    slog.Error(err.Error())
    return err
  }

  if affected, _ := result.RowsAffected(); 1 != affected {
    p := problem.Problem{}
    p.Title("Tag not created.")
    p.Detail("Could not create tag for an unknown reason.")
    p.Status(http.StatusInternalServerError)

    return &p
  }

  if err = tx.Commit(); nil != err {
    slog.Error(err.Error())
    return err
  }

  return nil
}

func (r *tagsRepository) Get(ctx context.Context) (tags []*model.Tag, err error) {
  getTagsQuery := `
  SELECT "id",
         "name",
         "created_at",
         "updated_at"
    FROM "tag"
ORDER BY lower("name");`

  ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
  defer cancel()

  result, err := r.db.QueryContext(ctx, getTagsQuery)
  if nil != err {
    slog.Error(err.Error())
    return nil, err
  }

  defer result.Close()

  tags = make([]*model.Tag, 0)

  for result.Next() {
    var tag model.Tag

    err = result.Scan(
      &tag.ID,
      &tag.Name,
      &tag.CreatedAt,
      &tag.UpdatedAt,
    )

    if nil != err {
      slog.Error(err.Error())
      return nil, err
    }

    tags = append(tags, &tag)
  }

  return tags, nil
}

func (r *tagsRepository) Update(ctx context.Context, id string, update *transfer.TagUpdate) error {
  tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
  if nil != err {
    slog.Error(err.Error())
    return err
  }

  defer tx.Rollback()

  updateArticleTagQuery := `
  UPDATE "article_tag"
     SET "tag_id" = @new_tag_id
   WHERE "tag_id" = @tag_id;`

  ctx1, cancel := context.WithTimeout(ctx, 5*time.Second)
  defer cancel()

  result, err := tx.ExecContext(ctx1, updateArticleTagQuery,
    sql.Named("tag_id", id),
    sql.Named("new_tag_id", update.ID),
  )

  if nil != err {
    slog.Error(err.Error())
    return err
  }

  updateTagQuery := `
  UPDATE "tag"
     SET "id" = @new_tag_id,
         "name" = @name,
         "updated_at" = current_timestamp
   WHERE "id" = @id;`

  ctx, cancel = context.WithTimeout(ctx, 3*time.Second)
  defer cancel()

  result, err = tx.ExecContext(ctx, updateTagQuery,
    sql.Named("id", id),
    sql.Named("new_tag_id", update.ID),
    sql.Named("name", update.Name),
  )

  if nil != err {
    slog.Error(err.Error())
    return err
  }

  if affected, _ := result.RowsAffected(); 1 != affected {
    return problem.NewNotFound(id, "tag")
  }

  if err = tx.Commit(); nil != err {
    slog.Error(err.Error())
    return err
  }

  return nil
}

func (r *tagsRepository) Remove(ctx context.Context, id string) error {
  tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
  if nil != err {
    slog.Error(err.Error())
    return err
  }

  defer tx.Rollback()

  removeTagQuery := `
  DELETE FROM "tag"
        WHERE "id" = $1;`

  ctx1, cancel := context.WithTimeout(ctx, 3*time.Second)
  defer cancel()

  result, err := tx.ExecContext(ctx1, removeTagQuery, id)

  if nil != err {
    slog.Error(err.Error())
    return err
  }

  if affected, _ := result.RowsAffected(); 1 != affected {
    return problem.NewNotFound(id, "tag")
  }

  removeFromAttachedArticlesQuery := `
  DELETE FROM "article_tag"
        WHERE "tag_id" = $1;`

  ctx1, cancel = context.WithTimeout(ctx, 5*time.Second)
  defer cancel()

  result, err = tx.ExecContext(ctx1, removeFromAttachedArticlesQuery, id)

  if nil != err {
    slog.Error(err.Error())
    return err
  }

  if _, err = result.RowsAffected(); nil != err {
    slog.Error(err.Error())
  }

  if err = tx.Commit(); nil != err {
    slog.Error(err.Error())
    return err
  }

  return nil
}
