package mocks

import (
  "context"
  "fontseca.dev/model"
  "fontseca.dev/transfer"
  "github.com/stretchr/testify/mock"
)

type TechnologyTagRepository struct {
  mock.Mock
}

func NewTechnologyTagRepository() *TechnologyTagRepository {
  return new(TechnologyTagRepository)
}

func (o *TechnologyTagRepository) Get(ctx context.Context) (technologies []*model.TechnologyTag, err error) {
  var args = o.Called(ctx)
  var arg0 = args.Get(0)
  if nil != arg0 {
    technologies = arg0.([]*model.TechnologyTag)
  }
  return technologies, args.Error(1)
}

func (o *TechnologyTagRepository) Add(ctx context.Context, creation *transfer.TechnologyTagCreation) (id string, err error) {
  var args = o.Called(ctx, creation)
  return args.String(0), args.Error(1)
}

func (o *TechnologyTagRepository) Exists(ctx context.Context, id string) (err error) {
  var args = o.Called(ctx, id)
  return args.Error(0)
}

func (o *TechnologyTagRepository) Update(ctx context.Context, id string, update *transfer.TechnologyTagUpdate) (updated bool, err error) {
  var args = o.Called(ctx, id, update)
  return args.Bool(0), args.Error(1)
}

func (o *TechnologyTagRepository) Remove(ctx context.Context, id string) (err error) {
  var args = o.Called(ctx, id)
  return args.Error(0)
}

type TechnologyTagService struct {
  mock.Mock
}

func NewTechnologyTagService() *TechnologyTagService {
  return new(TechnologyTagService)
}

func (o *TechnologyTagService) Get(ctx context.Context) (technologies []*model.TechnologyTag, err error) {
  var args = o.Called(ctx)
  var arg0 = args.Get(0)
  if nil != arg0 {
    technologies = arg0.([]*model.TechnologyTag)
  }
  return technologies, args.Error(1)
}

func (o *TechnologyTagService) Add(ctx context.Context, creation *transfer.TechnologyTagCreation) (id string, err error) {
  var args = o.Called(ctx, creation)
  return args.String(0), args.Error(1)
}

func (o *TechnologyTagService) Exists(ctx context.Context, id string) (err error) {
  var args = o.Called(ctx, id)
  return args.Error(0)
}

func (o *TechnologyTagService) Update(ctx context.Context, id string, update *transfer.TechnologyTagUpdate) (updated bool, err error) {
  var args = o.Called(ctx, id, update)
  return args.Bool(0), args.Error(1)
}

func (o *TechnologyTagService) Remove(ctx context.Context, id string) (err error) {
  var args = o.Called(ctx, id)
  return args.Error(0)
}
