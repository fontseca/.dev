package service

import (
  "context"
  "errors"
  "fontseca.dev/mocks"
  "fontseca.dev/model"
  "fontseca.dev/transfer"
  "github.com/google/uuid"
  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/mock"
  "strings"
  "testing"
)

func TestPatchesService_Get(t *testing.T) {
  const routine = "GetPatches"

  ctx := context.TODO()

  t.Run("success", func(t *testing.T) {
    expectedPatches := make([]*model.ArticlePatch, 3)

    r := mocks.NewArchiveRepository()
    r.On(routine, ctx).Return(expectedPatches, nil)

    articles, err := NewPatchesService(r).Get(ctx)

    assert.Equal(t, expectedPatches, articles)
    assert.NoError(t, err)
  })

  t.Run("gets a repository failure", func(t *testing.T) {
    unexpected := errors.New("unexpected error")

    r := mocks.NewArchiveRepository()
    r.On(routine, ctx).Return(nil, unexpected)

    articles, err := NewPatchesService(r).Get(ctx)

    assert.Nil(t, articles)
    assert.ErrorIs(t, err, unexpected)
  })
}

func TestPatchesService_Revise(t *testing.T) {
  const routine = "Revise"
  ctx := context.TODO()
  id := uuid.NewString()

  t.Run("success", func(t *testing.T) {
    revision := &transfer.ArticleRevision{
      Title:    "Consectetur! Adipiscing... Quis nostrud: ELIT?",
      Slug:     "consectetur-adipiscing-quis-nostrud-elit",
      ReadTime: 11,
      Content:  strings.Repeat("word ", 1999) + "word",
    }

    dirty := &transfer.ArticleRevision{
      Title:   " \t\n " + revision.Title + " \t\n ",
      Content: " \t\n " + revision.Content + " \t\n ",
    }

    r := mocks.NewArchiveRepository()
    r.On(routine, ctx, id, revision).Return(nil)

    assert.NoError(t, NewPatchesService(r).Revise(ctx, id, dirty))
  })

  t.Run("success: changing title", func(t *testing.T) {
    revision := &transfer.ArticleRevision{
      Title:    "Consectetur-Adipiscing!!... Quis nostrud: ELIT??? +-'\"",
      Slug:     "consectetur-adipiscing-quis-nostrud-elit",
      ReadTime: 1,
    }

    dirty := &transfer.ArticleRevision{
      Title: " \t\n " + revision.Title + " \t\n ",
    }

    r := mocks.NewArchiveRepository()
    r.On(routine, mock.Anything, mock.Anything, revision).Return(nil)

    assert.NoError(t, NewPatchesService(r).Revise(ctx, id, dirty))
  })

  t.Run("success: changing content", func(t *testing.T) {
    revision := &transfer.ArticleRevision{
      Content:  strings.Repeat("word ", 299) + "word",
      ReadTime: 2,
    }

    dirty := &transfer.ArticleRevision{
      Content: " \t\n " + revision.Content + " \t\n ",
    }

    r := mocks.NewArchiveRepository()
    r.On(routine, mock.Anything, mock.Anything, revision).Return(nil)

    assert.NoError(t, NewPatchesService(r).Revise(ctx, id, dirty))
  })

  t.Run("nil parameter: revision", func(t *testing.T) {
    r := mocks.NewArchiveRepository()
    r.AssertNotCalled(t, routine)
    assert.ErrorContains(t, NewPatchesService(r).Revise(ctx, id, nil), "nil value")
  })

  t.Run("wrong uuid: id", func(t *testing.T) {
    r := mocks.NewArchiveRepository()
    r.AssertNotCalled(t, routine)
    assert.Error(t, NewPatchesService(r).Revise(ctx, "x", &transfer.ArticleRevision{}))
  })

  t.Run("gets a repository failure", func(t *testing.T) {
    unexpected := errors.New("unexpected error")

    r := mocks.NewArchiveRepository()
    r.On(routine, mock.Anything, mock.Anything, mock.Anything).Return(unexpected)

    assert.ErrorIs(t, NewPatchesService(r).Revise(ctx, id, &transfer.ArticleRevision{}), unexpected)
  })
}

func TestPatchesService_Share(t *testing.T) {
  const routine = "Share"

  ctx := context.TODO()
  id := uuid.NewString()

  t.Run("success", func(t *testing.T) {
    expectedLink := "link-to-resource"

    r := mocks.NewArchiveRepository()
    r.On(routine, ctx, id).Return(expectedLink, nil)

    link, err := NewPatchesService(r).Share(ctx, id)

    assert.Equal(t, expectedLink, link)
    assert.NoError(t, err)
  })

  t.Run("wrong patch uuid", func(t *testing.T) {
    id = "e4d06ba7-f086-47dc-9f5e"

    r := mocks.NewArchiveRepository()
    r.AssertNotCalled(t, routine)

    link, err := NewPatchesService(r).Share(ctx, id)

    assert.Error(t, err)
    assert.Equal(t, "about:blank", link)
  })

  t.Run("gets a repository failure", func(t *testing.T) {
    unexpected := errors.New("unexpected error")

    r := mocks.NewArchiveRepository()
    r.On(routine, mock.Anything, mock.Anything, mock.Anything).Return("", unexpected)

    link, err := NewPatchesService(r).Share(ctx, uuid.NewString())

    assert.Equal(t, "about:blank", link)
    assert.ErrorIs(t, err, unexpected)
  })
}

func TestPatchesService_Discard(t *testing.T) {
  const routine = "Discard"

  ctx := context.TODO()
  id := uuid.NewString()

  t.Run("success", func(t *testing.T) {
    r := mocks.NewArchiveRepository()
    r.On(routine, ctx, id).Return(nil)

    assert.NoError(t, NewPatchesService(r).Discard(ctx, id))
  })

  t.Run("gets a repository failure", func(t *testing.T) {
    unexpected := errors.New("unexpected error")

    r := mocks.NewArchiveRepository()
    r.On(routine, mock.Anything, mock.Anything).Return(unexpected)

    assert.ErrorIs(t, NewPatchesService(r).Discard(ctx, id), unexpected)
  })

  t.Run("wrong uuid", func(t *testing.T) {
    id = "e4d06ba7-f086-47dc-9f5e"

    r := mocks.NewArchiveRepository()
    r.AssertNotCalled(t, routine)

    assert.Error(t, NewPatchesService(r).Discard(ctx, id))
  })
}

func TestPatchesService_Release(t *testing.T) {
  const routine = "Release"

  ctx := context.TODO()
  id := uuid.NewString()

  t.Run("success", func(t *testing.T) {
    r := mocks.NewArchiveRepository()
    r.On(routine, ctx, id).Return(nil)

    assert.NoError(t, NewPatchesService(r).Release(ctx, id))
  })

  t.Run("gets a repository failure", func(t *testing.T) {
    unexpected := errors.New("unexpected error")

    r := mocks.NewArchiveRepository()
    r.On(routine, mock.Anything, mock.Anything).Return(unexpected)

    assert.ErrorIs(t, NewPatchesService(r).Release(ctx, id), unexpected)
  })

  t.Run("wrong uuid", func(t *testing.T) {
    id = "e4d06ba7-f086-47dc-9f5e"

    r := mocks.NewArchiveRepository()
    r.AssertNotCalled(t, routine)

    assert.Error(t, NewPatchesService(r).Release(ctx, id))
  })
}
