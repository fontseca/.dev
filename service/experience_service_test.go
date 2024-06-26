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
  "testing"
  "time"
)

func TestExperienceService_Get(t *testing.T) {
  const routine = "Get"

  t.Run("success", func(t *testing.T) {
    var ctx = context.Background()
    var exp = make([]*model.Experience, 0)

    var r = mocks.NewExperienceRepository()
    r.On(routine, ctx, true).Return(exp, nil)
    res, err := NewExperienceService(r).Get(ctx, true)
    assert.NotNil(t, res)
    assert.NoError(t, err)

    r = mocks.NewExperienceRepository()
    r.On(routine, ctx, false).Return(exp, nil)
    res, err = NewExperienceService(r).Get(ctx, false)
    assert.NotNil(t, res)
    assert.NoError(t, err)

    r = mocks.NewExperienceRepository()
    r.On(routine, ctx, false).Return(exp, nil)
    res, err = NewExperienceService(r).Get(ctx)
    assert.NotNil(t, res)
    assert.NoError(t, err)
  })

  t.Run("error", func(t *testing.T) {
    var unexpected = errors.New("unexpected error")
    var ctx = context.Background()

    var r = mocks.NewExperienceRepository()
    r.On(routine, ctx, false).Return(nil, unexpected)
    res, err := NewExperienceService(r).Get(ctx)
    assert.Nil(t, res)
    assert.ErrorIs(t, err, unexpected)
  })
}

func TestExperienceService_GetByID(t *testing.T) {
  const routine = "GetByID"
  var id = uuid.New().String()

  t.Run("success", func(t *testing.T) {
    var r = mocks.NewExperienceRepository()
    var ctx = context.Background()
    r.On(routine, ctx, id).Return(new(model.Experience), nil)
    res, err := NewExperienceService(r).GetByID(ctx, id)
    assert.NotNil(t, res)
    assert.NoError(t, err)
  })

  t.Run("error", func(t *testing.T) {
    var unexpected = errors.New("unexpected error")
    var r = mocks.NewExperienceRepository()
    var ctx = context.Background()
    r.On(routine, ctx, id).Return(nil, unexpected)
    res, err := NewExperienceService(r).GetByID(ctx, id)
    assert.Nil(t, res)
    assert.ErrorIs(t, err, unexpected)
  })
}

func TestExperienceService_Save(t *testing.T) {
  const routine = "Save"

  t.Run("success", func(t *testing.T) {
    var expected = transfer.ExperienceCreation{
      Starts:   2020,
      Ends:     2023,
      JobTitle: "JobTitle",
      Company:  "Company",
      Country:  "Country",
      Summary:  "Summary",
    }
    var dirty = transfer.ExperienceCreation{
      Starts:   expected.Starts,
      Ends:     expected.Ends,
      JobTitle: " \n\t " + expected.JobTitle + " \n\t ",
      Company:  " \n\t " + expected.Company + " \n\t ",
      Country:  " \n\t " + expected.Country + " \n\t ",
      Summary:  " \n\t " + expected.Summary + " \n\t ",
    }
    var ctx = context.Background()
    var r = mocks.NewExperienceRepository()
    r.On(routine, ctx, &expected).Return(true, nil)
    res, err := NewExperienceService(r).Save(ctx, &dirty)
    assert.NoError(t, err)
    assert.True(t, res)
  })

  t.Run("error on nil creation", func(t *testing.T) {
    var r = mocks.NewExperienceRepository()
    var ctx = context.Background()
    r.AssertNotCalled(t, routine)
    res, err := NewExperienceService(r).Save(ctx, nil)
    assert.ErrorContains(t, err, "nil value for parameter: creation")
    assert.False(t, res)
  })

  t.Run("creation.Starts validations", func(t *testing.T) {
    var creation = transfer.ExperienceCreation{Starts: 2020, Ends: 2023}

    t.Run("2017<creation.Starts<=current_year", func(t *testing.T) {
      t.Run("fails:creation.Starts=2016", func(t *testing.T) {
        creation.Starts = 2016
        var ctx = context.Background()
        var r = mocks.NewExperienceRepository()
        r.On(routine, ctx, mock.AnythingOfType("*transfer.ExperienceCreation")).Return(true, nil)
        res, err := NewExperienceService(r).Save(ctx, &creation)
        assert.ErrorContains(t, err, "The provided data does not meet the required validation criteria")
        assert.False(t, res)
      })

      t.Run("meets:creation.Starts=2020", func(t *testing.T) {
        creation.Starts = 2020
        var ctx = context.Background()
        var r = mocks.NewExperienceRepository()
        r.On(routine, ctx, mock.AnythingOfType("*transfer.ExperienceCreation")).Return(true, nil)
        res, err := NewExperienceService(r).Save(ctx, &creation)
        assert.NoError(t, err)
        assert.True(t, res)
      })

      t.Run("meets:creation.Starts=current_year", func(t *testing.T) {
        creation.Starts = time.Now().Year()
        creation.Ends = creation.Starts
        var ctx = context.Background()
        var r = mocks.NewExperienceRepository()
        r.On(routine, ctx, mock.AnythingOfType("*transfer.ExperienceCreation")).Return(true, nil)
        res, err := NewExperienceService(r).Save(ctx, &creation)
        assert.NoError(t, err)
        assert.True(t, res)
      })

      t.Run("fails:creation.Starts=1+current_year", func(t *testing.T) {
        creation.Starts = 1 + time.Now().Year()
        var ctx = context.Background()
        var r = mocks.NewExperienceRepository()
        r.On(routine, ctx, mock.AnythingOfType("*transfer.ExperienceCreation")).Return(true, nil)
        res, err := NewExperienceService(r).Save(ctx, &creation)
        assert.ErrorContains(t, err, "The provided data does not meet the required validation criteria")
        assert.False(t, res)
      })
    })
  })

  t.Run("creation.Ends validations", func(t *testing.T) {
    var creation = transfer.ExperienceCreation{Starts: 2020, Ends: 2023}

    t.Run("creation.Starts<=creation.Ends<=current_year", func(t *testing.T) {
      t.Run("fails:creation.Ends=2016", func(t *testing.T) {
        creation.Ends = 2016
        var ctx = context.Background()
        var r = mocks.NewExperienceRepository()
        r.On(routine, ctx, mock.AnythingOfType("*transfer.ExperienceCreation")).Return(true, nil)
        res, err := NewExperienceService(r).Save(ctx, &creation)
        assert.ErrorContains(t, err, "The provided data does not meet the required validation criteria")
        assert.False(t, res)
      })

      t.Run("meets:creation.Ends=creation.Starts", func(t *testing.T) {
        creation.Ends = creation.Starts
        var ctx = context.Background()
        var r = mocks.NewExperienceRepository()
        r.On(routine, ctx, mock.AnythingOfType("*transfer.ExperienceCreation")).Return(true, nil)
        res, err := NewExperienceService(r).Save(ctx, &creation)
        assert.NoError(t, err)
        assert.True(t, res)
      })

      t.Run("meets:creation.Ends=1+creation.Starts", func(t *testing.T) {
        creation.Ends = 1 + creation.Starts
        var ctx = context.Background()
        var r = mocks.NewExperienceRepository()
        r.On(routine, ctx, mock.AnythingOfType("*transfer.ExperienceCreation")).Return(true, nil)
        res, err := NewExperienceService(r).Save(ctx, &creation)
        assert.NoError(t, err)
        assert.True(t, res)
      })

      t.Run("meets:creation.Ends=current_year", func(t *testing.T) {
        creation.Ends = time.Now().Year()
        var ctx = context.Background()
        var r = mocks.NewExperienceRepository()
        r.On(routine, ctx, mock.AnythingOfType("*transfer.ExperienceCreation")).Return(true, nil)
        res, err := NewExperienceService(r).Save(ctx, &creation)
        assert.NoError(t, err)
        assert.True(t, res)
      })

      t.Run("fails:creation.Ends=1+current_year", func(t *testing.T) {
        creation.Ends = 1 + time.Now().Year()
        var ctx = context.Background()
        var r = mocks.NewExperienceRepository()
        r.On(routine, ctx, mock.AnythingOfType("*transfer.ExperienceCreation")).Return(true, nil)
        res, err := NewExperienceService(r).Save(ctx, &creation)
        assert.ErrorContains(t, err, "The provided data does not meet the required validation criteria")
        assert.False(t, res)
      })
    })
  })

  t.Run("error", func(t *testing.T) {
    var unexpected = errors.New("unexpected error")
    var r = mocks.NewExperienceRepository()
    var ctx = context.Background()
    r.On(routine, ctx, mock.AnythingOfType("*transfer.ExperienceCreation")).Return(false, unexpected)
    res, err := NewExperienceService(r).Save(ctx, new(transfer.ExperienceCreation))
    assert.False(t, res)
    assert.ErrorIs(t, err, unexpected)
  })
}

func TestExperienceService_Update(t *testing.T) {
  const routine = "Update"
  var id = "   \t\n{7d7d4da0-093a-443b-b041-2da650381220}\n\t   "
  var expectedID = "7d7d4da0-093a-443b-b041-2da650381220"
  var ctx = context.Background()

  t.Run("success", func(t *testing.T) {
    var expected = transfer.ExperienceUpdate{
      Starts:   2020,
      Ends:     2023,
      JobTitle: "JobTitle",
      Company:  "Company",
      Country:  "Country",
      Summary:  "Summary",
    }
    var dirty = transfer.ExperienceUpdate{
      Starts:   expected.Starts,
      Ends:     expected.Ends,
      JobTitle: " \n\t " + expected.JobTitle + " \n\t ",
      Company:  " \n\t " + expected.Company + " \n\t ",
      Country:  " \n\t " + expected.Country + " \n\t ",
      Summary:  " \n\t " + expected.Summary + " \n\t ",
    }
    var r = mocks.NewExperienceRepository()
    r.On(routine, ctx, expectedID, &expected).Return(true, nil)
    res, err := NewExperienceService(r).Update(ctx, id, &dirty)
    assert.True(t, res)
    assert.NoError(t, err)
  })

  t.Run("update.Starts validations", func(t *testing.T) {
    var update = transfer.ExperienceUpdate{Starts: 2020}

    t.Run("2017<update.Starts<=current_year", func(t *testing.T) {
      t.Run("fails:update.Starts=2016", func(t *testing.T) {
        update.Starts = 2016
        var r = mocks.NewExperienceRepository()
        r.On(routine, ctx, mock.Anything, mock.AnythingOfType("*transfer.ExperienceUpdate")).Return(true, nil)
        res, err := NewExperienceService(r).Update(ctx, id, &update)
        assert.ErrorContains(t, err, "The provided data does not meet the required validation criteria")
        assert.False(t, res)
      })

      t.Run("meets:update.Starts=2020", func(t *testing.T) {
        update.Starts = 2020
        var r = mocks.NewExperienceRepository()
        r.On(routine, ctx, mock.Anything, mock.AnythingOfType("*transfer.ExperienceUpdate")).Return(true, nil)
        res, err := NewExperienceService(r).Update(ctx, id, &update)
        assert.NoError(t, err)
        assert.True(t, res)
      })

      t.Run("meets:update.Starts=current_year", func(t *testing.T) {
        update.Starts = time.Now().Year()
        update.Ends = update.Starts
        var r = mocks.NewExperienceRepository()
        r.On(routine, ctx, mock.Anything, mock.AnythingOfType("*transfer.ExperienceUpdate")).Return(true, nil)
        res, err := NewExperienceService(r).Update(ctx, id, &update)
        assert.NoError(t, err)
        assert.True(t, res)
      })

      t.Run("fails:update.Starts=1+current_year", func(t *testing.T) {
        update.Starts = 1 + time.Now().Year()
        var r = mocks.NewExperienceRepository()
        r.On(routine, ctx, mock.Anything, mock.AnythingOfType("*transfer.ExperienceUpdate")).Return(true, nil)
        res, err := NewExperienceService(r).Update(ctx, id, &update)
        assert.ErrorContains(t, err, "The provided data does not meet the required validation criteria")
        assert.False(t, res)
      })
    })
  })

  t.Run("update.Ends validations", func(t *testing.T) {
    var update = transfer.ExperienceUpdate{Ends: 2023}

    t.Run("update.Starts<=update.Ends<=current_year", func(t *testing.T) {
      t.Run("fails:update.Ends=2016", func(t *testing.T) {
        update.Ends = 2016
        var r = mocks.NewExperienceRepository()
        r.On(routine, ctx, mock.Anything, mock.AnythingOfType("*transfer.ExperienceUpdate")).Return(true, nil)
        res, err := NewExperienceService(r).Update(ctx, id, &update)
        assert.ErrorContains(t, err, "The provided data does not meet the required validation criteria")
        assert.False(t, res)
      })

      t.Run("meets:update.Ends=update.Starts", func(t *testing.T) {
        update.Starts = 2020
        update.Ends = update.Starts
        var r = mocks.NewExperienceRepository()
        r.On(routine, ctx, mock.Anything, mock.AnythingOfType("*transfer.ExperienceUpdate")).Return(true, nil)
        res, err := NewExperienceService(r).Update(ctx, id, &update)
        assert.NoError(t, err)
        assert.True(t, res)
      })

      t.Run("meets:update.Ends=1+update.Starts", func(t *testing.T) {
        update.Ends = 1 + update.Starts
        var r = mocks.NewExperienceRepository()
        r.On(routine, ctx, mock.Anything, mock.AnythingOfType("*transfer.ExperienceUpdate")).Return(true, nil)
        res, err := NewExperienceService(r).Update(ctx, id, &update)
        assert.NoError(t, err)
        assert.True(t, res)
      })

      t.Run("meets:update.Ends=current_year", func(t *testing.T) {
        update.Ends = time.Now().Year()
        var r = mocks.NewExperienceRepository()
        r.On(routine, ctx, mock.Anything, mock.AnythingOfType("*transfer.ExperienceUpdate")).Return(true, nil)
        res, err := NewExperienceService(r).Update(ctx, id, &update)
        assert.NoError(t, err)
        assert.True(t, res)
      })

      t.Run("fails:update.Starts>update.Ends", func(t *testing.T) {
        update.Starts = 2020
        update.Ends = 2019
        var r = mocks.NewExperienceRepository()
        r.On(routine, ctx, mock.Anything, mock.AnythingOfType("*transfer.ExperienceUpdate")).Return(true, nil)
        res, err := NewExperienceService(r).Update(ctx, id, &update)
        assert.ErrorContains(t, err, "The provided data does not meet the required validation criteria")
        assert.False(t, res)
      })

      t.Run("fails:update.Ends=1+current_year", func(t *testing.T) {
        update.Ends = 1 + time.Now().Year()
        var r = mocks.NewExperienceRepository()
        r.On(routine, ctx, mock.Anything, mock.AnythingOfType("*transfer.ExperienceUpdate")).Return(true, nil)
        res, err := NewExperienceService(r).Update(ctx, id, &update)
        assert.ErrorContains(t, err, "The provided data does not meet the required validation criteria")
        assert.False(t, res)
      })
    })
  })

  t.Run("error", func(t *testing.T) {
    var unexpected = errors.New("unexpected error")
    var r = mocks.NewExperienceRepository()
    r.On(routine, ctx, expectedID, mock.AnythingOfType("*transfer.ExperienceUpdate")).Return(false, unexpected)
    res, err := NewExperienceService(r).Update(ctx, id, new(transfer.ExperienceUpdate))
    assert.False(t, res)
    assert.ErrorIs(t, err, unexpected)
  })
}

func TestExperienceService_Remove(t *testing.T) {
  const routine = "Remove"
  var id = "   \t\n{7d7d4da0-093a-443b-b041-2da650381220}\n\t   "
  var expectedID = "7d7d4da0-093a-443b-b041-2da650381220"
  var ctx = context.Background()

  t.Run("success", func(t *testing.T) {
    var r = mocks.NewExperienceRepository()
    r.On(routine, ctx, expectedID).Return(nil)
    err := NewExperienceService(r).Remove(ctx, id)
    assert.NoError(t, err)
  })

  t.Run("error", func(t *testing.T) {
    var unexpected = errors.New("unexpected error")
    var r = mocks.NewExperienceRepository()
    r.On(routine, ctx, expectedID).Return(unexpected)
    err := NewExperienceService(r).Remove(ctx, id)
    assert.ErrorIs(t, err, unexpected)
  })
}
