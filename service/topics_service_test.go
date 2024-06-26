package service

import (
  "context"
  "errors"
  "fontseca.dev/mocks"
  "fontseca.dev/model"
  "fontseca.dev/transfer"
  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/mock"
  "testing"
)

func TestTopicsService_Add(t *testing.T) {
  const routine = "Add"

  ctx := context.TODO()
  creation := &transfer.TopicCreation{
    Name: "Consectetur! Adipiscing... Quis nostrud: ELIT?",
    ID:   "consectetur-adipiscing-quis-nostrud-elit",
  }

  t.Run("success", func(t *testing.T) {
    dirty := &transfer.TopicCreation{
      Name: " \n\t " + creation.Name + " \n\t ",
    }

    r := mocks.NewTopicsRepository()

    r.On(routine, ctx, creation).Return(nil)
    r.On("Get", ctx).Return([]*model.Topic{{}, {}}, nil)

    err := NewTopicsService(r).Add(ctx, dirty)

    assert.NoError(t, err)
  })

  t.Run("gets a repository failure", func(t *testing.T) {
    unexpected := errors.New("unexpected error")
    r := mocks.NewTopicsRepository()
    r.On(routine, ctx, mock.Anything).Return(unexpected)

    err := NewTopicsService(r).Add(ctx, creation)
    assert.ErrorIs(t, err, unexpected)
  })
}

func TestTopicsService_Get(t *testing.T) {
  const routine = "Get"

  ctx := context.TODO()

  t.Run("success without cache", func(t *testing.T) {
    expectedTopics := []*model.Topic{{}, {}, {}}

    r := mocks.NewTopicsRepository()
    r.On(routine, ctx).Return(expectedTopics, nil)

    s := NewTopicsService(r).(*topicsService)

    s.cache = nil

    topics, err := s.Get(ctx)

    assert.Equal(t, expectedTopics, topics)
    assert.NoError(t, err)
  })

  t.Run("success with cache", func(t *testing.T) {
    expectedTopics := []*model.Topic{{}, {}, {}}

    r := mocks.NewTopicsRepository()
    r.AssertNotCalled(t, routine)

    s := NewTopicsService(r).(*topicsService)

    s.cache = expectedTopics

    topics, err := s.Get(ctx)

    assert.Equal(t, expectedTopics, topics)
    assert.NoError(t, err)
  })

  t.Run("gets a repository failure", func(t *testing.T) {
    unexpected := errors.New("unexpected error")

    r := mocks.NewTopicsRepository()
    r.On(routine, ctx).Return(nil, unexpected)

    topics, err := NewTopicsService(r).Get(ctx)

    assert.Nil(t, topics)
    assert.ErrorIs(t, err, unexpected)
  })
}

func TestTopicsService_Update(t *testing.T) {
  const routine = "Update"

  ctx := context.TODO()
  id := "consectetur-adipiscing-quis-nostrud-elit"

  update := &transfer.TopicUpdate{
    ID:   "consectetur-adipiscing-quis-nostrud-elit",
    Name: "Consectetur! Adipiscing... Quis nostrud: ELIT?",
  }

  t.Run("success", func(t *testing.T) {
    dirty := &transfer.TopicUpdate{
      Name: " \n\t " + update.Name + " \n\t ",
    }

    r := mocks.NewTopicsRepository()

    r.On(routine, ctx, id, update).Return(nil)
    r.On("Get", ctx).Return([]*model.Topic{{}, {}}, nil)

    err := NewTopicsService(r).Update(ctx, id, dirty)

    assert.NoError(t, err)
  })

  t.Run("gets a repository failure", func(t *testing.T) {
    unexpected := errors.New("unexpected error")
    r := mocks.NewTopicsRepository()
    r.On(routine, ctx, mock.Anything, mock.Anything).Return(unexpected)

    err := NewTopicsService(r).Update(ctx, id, update)
    assert.ErrorIs(t, err, unexpected)
  })
}

func TestTopicsService_Remove(t *testing.T) {
  const routine = "Remove"

  ctx := context.TODO()
  id := "id"

  t.Run("success", func(t *testing.T) {
    r := mocks.NewTopicsRepository()
    r.On(routine, ctx, id).Return(nil)
    r.On("Get", ctx).Return([]*model.Topic{{}, {}}, nil)

    assert.NoError(t, NewTopicsService(r).Remove(ctx, id))
  })

  t.Run("gets a repository failure", func(t *testing.T) {
    unexpected := errors.New("unexpected error")

    r := mocks.NewTopicsRepository()
    r.On(routine, ctx, id).Return(unexpected)

    assert.ErrorIs(t, NewTopicsService(r).Remove(ctx, id), unexpected)
  })
}
