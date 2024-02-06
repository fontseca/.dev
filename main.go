package main

import (
  "context"
  "database/sql"
  "fmt"
  "github.com/google/uuid"
  "github.com/mattn/go-sqlite3"
  "io"
  "log"
  "log/slog"
  "os"
  "path/filepath"
  "time"
)

// table contains information about a relation in the database.
type table struct {
  name  string
  query string
}

// exists checks if the table t.name is already created in the transaction tx.
func (t *table) exists(ctx context.Context, tx *sql.Tx) bool {
  if nil == tx {
    return false
  }
  var query = `
  SELECT count (1)
    FROM "sqlite_master"
   WHERE "type" = 'table'
     AND "name" = $1;`
  ctx, cancel := context.WithTimeout(ctx, time.Second)
  defer cancel()
  var result = tx.QueryRowContext(ctx, query, t.name)
  var err = result.Err()
  if nil != err {
    err = fmt.Errorf("checking existence of table %q: %v", t.name, err)
    if rollbackErr := tx.Rollback(); nil != rollbackErr {
      log.Fatalf("unable to rollback: %v: %v", err, rollbackErr)
    }
    log.Fatal(err)
  }
  var n int
  err = result.Scan(&n)
  if nil != err {
    log.Fatal(err)
  }
  return n >= 1
}

// create attempts to create the table t in the transaction tx.
func (t *table) create(ctx context.Context, tx *sql.Tx) {
  if nil == tx {
    return
  }
  ctx, cancel := context.WithTimeout(ctx, time.Second)
  defer cancel()
  if _, err := tx.ExecContext(ctx, t.query); nil != err {
    err = fmt.Errorf("creating table %q: %v", t.name, err)
    if rollbackErr := tx.Rollback(); nil != rollbackErr {
      log.Fatalf("unable to rollback: %v: %v", err, rollbackErr)
    }
    log.Fatal(err)
  }
}

func main() {
  sql.Register("sqlite3_custom", &sqlite3.SQLiteDriver{
    ConnectHook: func(conn *sqlite3.SQLiteConn) error {
      if err := conn.RegisterFunc(
        "uuid_generate_v4",
        func() string { return uuid.New().String() },
        true,
      ); nil != err {
        return err
      }

      if err := conn.RegisterFunc(
        "uuid_nil",
        func() string { return uuid.Nil.String() },
        true,
      ); nil != err {
        return err
      }

      return nil
    },
  })

  var db, err = sql.Open("sqlite3_custom", "./db.sqlite")
  if nil != err {
    log.Fatal(err)
  }

  defer func(db *sql.DB) {
    err := db.Close()
    if err != nil {
      log.Fatal(err)
    }
  }(db)

  if err = db.Ping(); nil != err {
    log.Fatal(err)
  }

  var tables = []table{
    {
      name: "me",
      query: `
      CREATE TABLE "me"
      (
        "username"      VARCHAR(64) NOT NULL DEFAULT 'fontseca.dev',
        "name"          VARCHAR(31) NOT NULL DEFAULT 'Jeremy Fonseca',
        "summary"       VARCHAR(1024) NOT NULL,
        "job_title"     VARCHAR(64) NOT NULL DEFAULT 'Back-End Software Developer',
        "email"         VARCHAR(254) NOT NULL,
        "photo_url"     VARCHAR(2048) NOT NULL DEFAULT 'about:blank',
        "resume_url"    VARCHAR(2048) NOT NULL DEFAULT 'about:blank',
        "coding_since"  INT NOT NULL DEFAULT 2017,
        "company"       VARCHAR(64),
        "location"      VARCHAR(64),
        "hireable"      BOOLEAN NOT NULL DEFAULT TRUE,
        "github_url"    VARCHAR(2048) NOT NULL DEFAULT 'https://github.com/fontseca',
        "linkedin_url"  VARCHAR(2048) NOT NULL DEFAULT 'about:blank',
        "youtube_url"   VARCHAR(2048) NOT NULL DEFAULT 'about:blank',
        "twitter_url"   VARCHAR(2048) NOT NULL DEFAULT 'about:blank',
        "instagram_url" VARCHAR(2048) NOT NULL DEFAULT 'about:blank',
        "created_at"    TIMESTAMP NOT NULL DEFAULT current_timestamp,
        "updated_at"    TIMESTAMP NOT NULL DEFAULT current_timestamp,
        CHECK ("coding_since" = 2017)
      );`,
    },
    {
      name: "experience",
      query: `
      CREATE TABLE "experience"
      (
        "id"         VARCHAR(36) NOT NULL PRIMARY KEY DEFAULT (uuid_generate_v4 ()),
        "starts"     INT NOT NULL,
        "ends"       INT NULL,
        "job_title"  VARCHAR(64) NOT NULL DEFAULT 'Back-End Software Developer',
        "company"    VARCHAR(64) NOT NULL,
        "country"    VARCHAR(64),
        "summary"    TEXT NOT NULL,
        "active"     BOOLEAN DEFAULT FALSE,
        "hidden"     BOOLEAN DEFAULT FALSE,
        "created_at" TIMESTAMP NOT NULL DEFAULT current_timestamp,
        "updated_at" TIMESTAMP NOT NULL DEFAULT current_timestamp,
        CHECK ("starts" > 2003 AND "ends" > 2003)
      );`,
    },
    {
      name: "project",
      query: `
      CREATE TABLE "project"
      (
        "id"               VARCHAR(36) NOT NULL PRIMARY KEY DEFAULT (uuid_generate_v4 ()),
        "name"             VARCHAR(64) NOT NULL,
        "homepage"         VARCHAR(2048) NOT NULL DEFAULT 'about:blank',
        "language"         VARCHAR(64),
        "summary"          VARCHAR(1024) NOT NULL,
        "content"          TEXT NOT NULL,
        "estimated_time"   INT,
        "first_image_url"  VARCHAR(2048) NOT NULL DEFAULT 'about:blank',
        "second_image_url" VARCHAR(2048) NOT NULL DEFAULT 'about:blank',
        "github_url"       VARCHAR(2048) NOT NULL DEFAULT 'about:blank',
        "collection_url"   VARCHAR(2048) NOT NULL DEFAULT 'about:blank',
        "playground_url"   VARCHAR(2048) NOT NULL DEFAULT 'about:blank',
        "playable"         BOOLEAN NOT NULL DEFAULT FALSE,
        "archived"         BOOLEAN NOT NULL DEFAULT FALSE,
        "finished"         BOOLEAN DEFAULT FALSE,
        "created_at"       TIMESTAMP NOT NULL DEFAULT current_timestamp,
        "updated_at"       TIMESTAMP NOT NULL DEFAULT current_timestamp,
        CHECK ("estimated_time" > 0)
      );`,
    },
    {
      name: "technology_tag",
      query: `
      CREATE TABLE "technology_tag"
      (
        "id"   VARCHAR(36) NOT NULL PRIMARY KEY DEFAULT (uuid_generate_v4 ()),
        "name" VARCHAR(64) NOT NULL
      );`,
    },
    {
      name: "project_technology_tag",
      query: `
      CREATE TABLE "project_technology_tag"
      (
        "project_id"        VARCHAR(36) NOT NULL REFERENCES "project" ("id"),
        "technology_tag_id" VARCHAR(36) NOT NULL REFERENCES "technology_tag" ("id")
      );`,
    },
    {
      name: "article",
      query: `
      CREATE TABLE "article"
      (
        "id"           VARCHAR(36) NOT NULL PRIMARY KEY DEFAULT (uuid_generate_v4 ()),
        "title"        VARCHAR(64) NOT NULL,
        "author"       VARCHAR(512) NOT NULL REFERENCES "me" ("username"),
        "slug"         VARCHAR(2024) NOT NULL,
        "cover_url"    VARCHAR(2048) NOT NULL DEFAULT 'about:blank',
        "reading_time" INT NOT NULL,
        "description"  VARCHAR(1024) NOT NULL,
        "content"      TEXT NOT NULL,
        "pinned"       BOOLEAN NOT NULL DEFAULT FALSE,
        "draft"        BOOLEAN NOT NULL DEFAULT TRUE,
        "archived"     BOOLEAN NOT NULL DEFAULT FALSE,
        "published_at" TIMESTAMP NOT NULL,
        "modified_at"  TIMESTAMP NOT NULL,
        "created_at"   TIMESTAMP NOT NULL DEFAULT current_timestamp,
        "updated_at"   TIMESTAMP NOT NULL DEFAULT current_timestamp
        CHECK ("reading_time" > 0)
      );`,
    },
    {
      name: "tag",
      query: `
      CREATE TABLE "tag"
      (
        "id"   VARCHAR(36) NOT NULL PRIMARY KEY DEFAULT (uuid_generate_v4 ()),
        "name" VARCHAR(64) NOT NULL
      );`,
    },
    {
      name: "article_tag",
      query: `
      CREATE TABLE "article_tag"
      (
        "article_id" VARCHAR(36) NOT NULL REFERENCES "article" ("id"),
        "tag_id"     VARCHAR(36) NOT NULL REFERENCES "tag" ("id")
      );`,
    },
  }

  var ctx = context.Background()
  tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
  if nil != err {
    log.Fatal(err)
  }

  for _, t := range tables {
    if !t.exists(ctx, tx) {
      t.create(ctx, tx)
    }
  }

  if err = tx.Commit(); nil != err {
    log.Fatal(err)
  }

  logfile, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
  if nil != err {
    log.Fatal(err)
  }
  defer logfile.Close()

  var multiWriter = io.MultiWriter(os.Stderr, logfile)
  var logger = slog.New(slog.NewJSONHandler(multiWriter,
    &slog.HandlerOptions{
      AddSource: true,
      ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
        if slog.SourceKey == a.Key {
          var source, _ = a.Value.Any().(*slog.Source)
          if nil != source {
            source.File = filepath.Base(source.File)
          }
        }
        return a
      },
    }))

  slog.SetDefault(logger)
}