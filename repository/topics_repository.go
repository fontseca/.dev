package repository

import (
  "context"
  "database/sql"
  "fontseca.dev/model"
  "fontseca.dev/transfer"
  "github.com/google/uuid"
  "log/slog"
  "time"
)

// TopicsRepository is a low level API that provides methods for interacting
// with topics in the database.
type TopicsRepository interface {
  // Add adds a new topic.
  Add(ctx context.Context, creation *transfer.TopicCreation) (insertedUUID string, err error)

  // Get retrieves all the topics.
  Get(ctx context.Context) (topics []*model.Topic, err error)

  // Update updates an existing topic.
  Update(ctx context.Context, id string, update *transfer.TopicUpdate) error

  // Remove removes a topic and detaches it from any article that uses it.
  Remove(ctx context.Context, id string) error
}

type topicsRepository struct {
  db *sql.DB
}

func NewTopicsRepository(db *sql.DB) TopicsRepository {
  return &topicsRepository{db}
}

func (r *topicsRepository) Add(ctx context.Context, creation *transfer.TopicCreation) (insertedUUID string, err error) {
  tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
  if nil != err {
    return uuid.Nil.String(), err
  }

  defer tx.Rollback()

  addTopicQuery := `
  INSERT INTO "topic" (name)
               VALUES (@name)
    RETURNING "uuid";`

  ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
  defer cancel()

  result := tx.QueryRowContext(ctx, addTopicQuery,
    sql.Named("name", creation.Name),
  )

  err = result.Scan(&insertedUUID)
  if nil != err {
    slog.Error(err.Error())
    return uuid.Nil.String(), err
  }

  if err = tx.Commit(); nil != err {
    slog.Error(err.Error())
    return uuid.Nil.String(), err
  }

  return insertedUUID, nil
}

func (r *topicsRepository) Get(ctx context.Context) (topics []*model.Topic, err error) {
  // TODO implement me
  panic("implement me")
}

func (r *topicsRepository) Update(ctx context.Context, id string, update *transfer.TopicUpdate) error {
  // TODO implement me
  panic("implement me")
}

func (r *topicsRepository) Remove(ctx context.Context, id string) error {
  // TODO implement me
  panic("implement me")
}
