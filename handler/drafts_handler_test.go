package handler

import (
  "errors"
  "fontseca.dev/mocks"
  "fontseca.dev/model"
  "fontseca.dev/problem"
  "fontseca.dev/transfer"
  "github.com/gin-gonic/gin"
  "github.com/google/uuid"
  "github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/mock"
  "net/http"
  "net/http/httptest"
  "testing"
)

func TestDraftsHandler_Start(t *testing.T) {
  const (
    routine = "Draft"
    method  = http.MethodPost
    target  = "/archive.drafts.start"
  )

  creation := &transfer.ArticleCreation{
    Title:   "Title",
    Content: "Content",
  }

  id := uuid.New()

  request := httptest.NewRequest(method, target, nil)
  _ = request.ParseForm()

  request.PostForm.Add("title", creation.Title)
  request.PostForm.Add("content", creation.Content)

  t.Run("success", func(t *testing.T) {
    expectedStatusCode := http.StatusCreated
    expectedResponse := string(marshal(t, gin.H{"draft_uuid": id.String()}))

    s := mocks.NewDraftsService()
    s.On(routine, mock.AnythingOfType("*gin.Context"), creation).Return(id, nil)

    engine := gin.Default()
    engine.POST(target, NewDraftsHandler(s).Start)

    recorder := httptest.NewRecorder()

    engine.ServeHTTP(recorder, request)

    assert.Equal(t, expectedStatusCode, recorder.Code)
    assert.Equal(t, expectedResponse, recorder.Body.String())
    assert.Empty(t, recorder.Result().Cookies())
  })

  t.Run("expected problem detail", func(t *testing.T) {
    expectedStatusCode := http.StatusBadRequest
    expectBodyContains := "Expected problem detail."

    expected := &problem.Problem{}
    expected.Status(expectedStatusCode)
    expected.Detail(expectBodyContains)

    s := mocks.NewDraftsService()
    s.On(routine, mock.AnythingOfType("*gin.Context"), creation).Return(uuid.Nil, expected)

    engine := gin.Default()
    engine.POST(target, NewDraftsHandler(s).Start)

    recorder := httptest.NewRecorder()

    engine.ServeHTTP(recorder, request)

    assert.Equal(t, expectedStatusCode, recorder.Code)
    assert.Contains(t, recorder.Body.String(), expectBodyContains)
    assert.Empty(t, recorder.Result().Cookies())
    assert.Contains(t, recorder.Result().Header.Get("Content-Type"), "application/problem+json")
  })

  t.Run("unexpected error", func(t *testing.T) {
    unexpected := errors.New("unexpected error")
    expectedStatusCode := http.StatusInternalServerError
    expectBodyContains := "An unexpected error occurred while processing your request"

    s := mocks.NewDraftsService()
    s.On(routine, mock.AnythingOfType("*gin.Context"), creation).Return(uuid.Nil, unexpected)

    engine := gin.Default()
    engine.POST(target, NewDraftsHandler(s).Start)

    recorder := httptest.NewRecorder()

    engine.ServeHTTP(recorder, request)

    assert.Equal(t, expectedStatusCode, recorder.Code)
    assert.Contains(t, recorder.Body.String(), expectBodyContains)
    assert.Empty(t, recorder.Result().Cookies())
    assert.Contains(t, recorder.Result().Header.Get("Content-Type"), "application/problem+json")
  })
}

func TestDraftsHandler_Publish(t *testing.T) {
  const (
    routine = "Publish"
    method  = http.MethodPost
    target  = "/archive.drafts.publish"
  )

  id := uuid.NewString()

  request := httptest.NewRequest(method, target, nil)
  _ = request.ParseForm()

  request.PostForm.Add("draft_uuid", id)

  t.Run("success", func(t *testing.T) {
    expectedStatusCode := http.StatusNoContent

    s := mocks.NewDraftsService()
    s.On(routine, mock.AnythingOfType("*gin.Context"), id).Return(nil)

    engine := gin.Default()
    engine.POST(target, NewDraftsHandler(s).Publish)

    recorder := httptest.NewRecorder()

    engine.ServeHTTP(recorder, request)

    assert.Equal(t, expectedStatusCode, recorder.Code)
    assert.Empty(t, recorder.Body.String())
    assert.Empty(t, recorder.Result().Cookies())
  })

  t.Run("expected problem detail", func(t *testing.T) {
    expectedStatusCode := http.StatusBadRequest
    expectBodyContains := "Expected problem detail."

    expected := &problem.Problem{}
    expected.Status(expectedStatusCode)
    expected.Detail(expectBodyContains)

    s := mocks.NewDraftsService()
    s.On(routine, mock.AnythingOfType("*gin.Context"), id).Return(expected)

    engine := gin.Default()
    engine.POST(target, NewDraftsHandler(s).Publish)

    recorder := httptest.NewRecorder()

    engine.ServeHTTP(recorder, request)

    assert.Equal(t, expectedStatusCode, recorder.Code)
    assert.Contains(t, recorder.Body.String(), expectBodyContains)
    assert.Empty(t, recorder.Result().Cookies())
    assert.Contains(t, recorder.Result().Header.Get("Content-Type"), "application/problem+json")
  })

  t.Run("unexpected error", func(t *testing.T) {
    unexpected := errors.New("unexpected error")
    expectedStatusCode := http.StatusInternalServerError
    expectBodyContains := "An unexpected error occurred while processing your request"

    s := mocks.NewDraftsService()
    s.On(routine, mock.AnythingOfType("*gin.Context"), id).Return(unexpected)

    engine := gin.Default()
    engine.POST(target, NewDraftsHandler(s).Publish)

    recorder := httptest.NewRecorder()

    engine.ServeHTTP(recorder, request)

    assert.Equal(t, expectedStatusCode, recorder.Code)
    assert.Contains(t, recorder.Body.String(), expectBodyContains)
    assert.Empty(t, recorder.Result().Cookies())
    assert.Contains(t, recorder.Result().Header.Get("Content-Type"), "application/problem+json")
  })
}

func TestDraftsHandler_Get(t *testing.T) {
  const (
    routine = "Get"
    method  = http.MethodGet
    target  = "/archive.drafts.list"
  )

  request := httptest.NewRequest(method, target, nil)
  drafts := []*transfer.Article{{}, {}, {}}

  t.Run("success without search", func(t *testing.T) {
    expectedStatusCode := http.StatusOK
    expectedBody := string(marshal(t, drafts))

    s := mocks.NewDraftsService()
    s.On(routine, mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*transfer.ArticleFilter")).Return(drafts, nil)

    engine := gin.Default()
    engine.GET(target, NewDraftsHandler(s).Get)

    recorder := httptest.NewRecorder()

    engine.ServeHTTP(recorder, request)

    assert.Equal(t, expectedStatusCode, recorder.Code)
    assert.Equal(t, expectedBody, recorder.Body.String())
    assert.Empty(t, recorder.Result().Cookies())
  })

  t.Run("expected problem detail", func(t *testing.T) {
    expectedStatusCode := http.StatusBadRequest
    expectBodyContains := "Expected problem detail."

    expected := &problem.Problem{}
    expected.Status(expectedStatusCode)
    expected.Detail(expectBodyContains)

    s := mocks.NewDraftsService()
    s.On(routine, mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*transfer.ArticleFilter")).Return(nil, expected)

    engine := gin.Default()
    engine.GET(target, NewDraftsHandler(s).Get)

    recorder := httptest.NewRecorder()

    engine.ServeHTTP(recorder, request)

    assert.Equal(t, expectedStatusCode, recorder.Code)
    assert.Contains(t, recorder.Body.String(), expectBodyContains)
    assert.Empty(t, recorder.Result().Cookies())
    assert.Contains(t, recorder.Result().Header.Get("Content-Type"), "application/problem+json")
  })

  t.Run("unexpected error", func(t *testing.T) {
    unexpected := errors.New("unexpected error")
    expectedStatusCode := http.StatusInternalServerError
    expectBodyContains := "An unexpected error occurred while processing your request"

    s := mocks.NewDraftsService()
    s.On(routine, mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*transfer.ArticleFilter")).Return(nil, unexpected)

    engine := gin.Default()
    engine.GET(target, NewDraftsHandler(s).Get)

    recorder := httptest.NewRecorder()

    engine.ServeHTTP(recorder, request)

    assert.Equal(t, expectedStatusCode, recorder.Code)
    assert.Contains(t, recorder.Body.String(), expectBodyContains)
    assert.Empty(t, recorder.Result().Cookies())
    assert.Contains(t, recorder.Result().Header.Get("Content-Type"), "application/problem+json")
  })
}

func TestDraftsHandler_GetByID(t *testing.T) {
  const (
    routine = "GetByID"
    method  = http.MethodGet
    target  = "/archive.drafts.info"
  )

  request := httptest.NewRequest(method, target, nil)
  id := uuid.NewString()

  draft := &model.Article{
    UUID: uuid.MustParse(id),
  }

  request.URL.RawQuery = request.URL.RawQuery + "&draft_uuid=" + id

  t.Run("success", func(t *testing.T) {
    expectedStatusCode := http.StatusOK
    expectedBody := string(marshal(t, draft))

    s := mocks.NewDraftsService()
    s.On(routine, mock.AnythingOfType("*gin.Context"), id).Return(draft, nil)

    engine := gin.Default()
    engine.GET(target, NewDraftsHandler(s).GetByID)

    recorder := httptest.NewRecorder()

    engine.ServeHTTP(recorder, request)

    assert.Equal(t, expectedStatusCode, recorder.Code)
    assert.Equal(t, expectedBody, recorder.Body.String())
    assert.Empty(t, recorder.Result().Cookies())
  })

  t.Run("expected problem detail", func(t *testing.T) {
    expectedStatusCode := http.StatusBadRequest
    expectBodyContains := "Expected problem detail."

    expected := &problem.Problem{}
    expected.Status(expectedStatusCode)
    expected.Detail(expectBodyContains)

    s := mocks.NewDraftsService()
    s.On(routine, mock.AnythingOfType("*gin.Context"), id).Return(nil, expected)

    engine := gin.Default()
    engine.GET(target, NewDraftsHandler(s).GetByID)

    recorder := httptest.NewRecorder()

    engine.ServeHTTP(recorder, request)

    assert.Equal(t, expectedStatusCode, recorder.Code)
    assert.Contains(t, recorder.Body.String(), expectBodyContains)
    assert.Empty(t, recorder.Result().Cookies())
    assert.Contains(t, recorder.Result().Header.Get("Content-Type"), "application/problem+json")
  })

  t.Run("unexpected error", func(t *testing.T) {
    unexpected := errors.New("unexpected error")
    expectedStatusCode := http.StatusInternalServerError
    expectBodyContains := "An unexpected error occurred while processing your request"

    s := mocks.NewDraftsService()
    s.On(routine, mock.AnythingOfType("*gin.Context"), id).Return(nil, unexpected)

    engine := gin.Default()
    engine.GET(target, NewDraftsHandler(s).GetByID)

    recorder := httptest.NewRecorder()

    engine.ServeHTTP(recorder, request)

    assert.Equal(t, expectedStatusCode, recorder.Code)
    assert.Contains(t, recorder.Body.String(), expectBodyContains)
    assert.Empty(t, recorder.Result().Cookies())
    assert.Contains(t, recorder.Result().Header.Get("Content-Type"), "application/problem+json")
  })
}

func TestDraftsHandler_AddTag(t *testing.T) {
  const (
    routine = "AddTag"
    method  = http.MethodPost
    target  = "/archive.drafts.tags.add"
  )

  request := httptest.NewRequest(method, target, nil)
  draftUUID := uuid.NewString()
  tagID := uuid.NewString()

  _ = request.ParseForm()

  request.PostForm.Add("draft_uuid", draftUUID)
  request.PostForm.Add("tag_id", tagID)

  t.Run("success", func(t *testing.T) {
    expectedStatusCode := http.StatusNoContent

    s := mocks.NewDraftsService()
    s.On(routine, mock.AnythingOfType("*gin.Context"), draftUUID, tagID).Return(nil)

    engine := gin.Default()
    engine.POST(target, NewDraftsHandler(s).AddTag)

    recorder := httptest.NewRecorder()

    engine.ServeHTTP(recorder, request)

    assert.Equal(t, expectedStatusCode, recorder.Code)
    assert.Empty(t, recorder.Body)
    assert.Empty(t, recorder.Result().Cookies())
  })

  t.Run("expected problem detail", func(t *testing.T) {
    expectedStatusCode := http.StatusBadRequest
    expectBodyContains := "Expected problem detail."

    expected := &problem.Problem{}
    expected.Status(expectedStatusCode)
    expected.Detail(expectBodyContains)

    s := mocks.NewDraftsService()
    s.On(routine, mock.AnythingOfType("*gin.Context"), draftUUID, tagID).Return(expected)

    engine := gin.Default()
    engine.POST(target, NewDraftsHandler(s).AddTag)

    recorder := httptest.NewRecorder()

    engine.ServeHTTP(recorder, request)

    assert.Equal(t, expectedStatusCode, recorder.Code)
    assert.Contains(t, recorder.Body.String(), expectBodyContains)
    assert.Empty(t, recorder.Result().Cookies())
    assert.Contains(t, recorder.Result().Header.Get("Content-Type"), "application/problem+json")
  })

  t.Run("unexpected error", func(t *testing.T) {
    unexpected := errors.New("unexpected error")
    expectedStatusCode := http.StatusInternalServerError
    expectBodyContains := "An unexpected error occurred while processing your request"

    s := mocks.NewDraftsService()
    s.On(routine, mock.AnythingOfType("*gin.Context"), draftUUID, tagID).Return(unexpected)

    engine := gin.Default()
    engine.POST(target, NewDraftsHandler(s).AddTag)

    recorder := httptest.NewRecorder()

    engine.ServeHTTP(recorder, request)

    assert.Equal(t, expectedStatusCode, recorder.Code)
    assert.Contains(t, recorder.Body.String(), expectBodyContains)
    assert.Empty(t, recorder.Result().Cookies())
    assert.Contains(t, recorder.Result().Header.Get("Content-Type"), "application/problem+json")
  })
}

func TestDraftsHandler_RemoveTag(t *testing.T) {
  const (
    routine = "RemoveTag"
    method  = http.MethodPost
    target  = "/archive.drafts.tags.remove"
  )

  request := httptest.NewRequest(method, target, nil)
  draftUUID := uuid.NewString()
  tagID := uuid.NewString()

  _ = request.ParseForm()

  request.PostForm.Add("draft_uuid", draftUUID)
  request.PostForm.Add("tag_id", tagID)

  t.Run("success", func(t *testing.T) {
    expectedStatusCode := http.StatusNoContent

    s := mocks.NewDraftsService()
    s.On(routine, mock.AnythingOfType("*gin.Context"), draftUUID, tagID).Return(nil)

    engine := gin.Default()
    engine.POST(target, NewDraftsHandler(s).RemoveTag)

    recorder := httptest.NewRecorder()

    engine.ServeHTTP(recorder, request)

    assert.Equal(t, expectedStatusCode, recorder.Code)
    assert.Empty(t, recorder.Body)
    assert.Empty(t, recorder.Result().Cookies())
  })

  t.Run("expected problem detail", func(t *testing.T) {
    expectedStatusCode := http.StatusBadRequest
    expectBodyContains := "Expected problem detail."

    expected := &problem.Problem{}
    expected.Status(expectedStatusCode)
    expected.Detail(expectBodyContains)

    s := mocks.NewDraftsService()
    s.On(routine, mock.AnythingOfType("*gin.Context"), draftUUID, tagID).Return(expected)

    engine := gin.Default()
    engine.POST(target, NewDraftsHandler(s).RemoveTag)

    recorder := httptest.NewRecorder()

    engine.ServeHTTP(recorder, request)

    assert.Equal(t, expectedStatusCode, recorder.Code)
    assert.Contains(t, recorder.Body.String(), expectBodyContains)
    assert.Empty(t, recorder.Result().Cookies())
    assert.Contains(t, recorder.Result().Header.Get("Content-Type"), "application/problem+json")
  })

  t.Run("unexpected error", func(t *testing.T) {
    unexpected := errors.New("unexpected error")
    expectedStatusCode := http.StatusInternalServerError
    expectBodyContains := "An unexpected error occurred while processing your request"

    s := mocks.NewDraftsService()
    s.On(routine, mock.AnythingOfType("*gin.Context"), draftUUID, tagID).Return(unexpected)

    engine := gin.Default()
    engine.POST(target, NewDraftsHandler(s).RemoveTag)

    recorder := httptest.NewRecorder()

    engine.ServeHTTP(recorder, request)

    assert.Equal(t, expectedStatusCode, recorder.Code)
    assert.Contains(t, recorder.Body.String(), expectBodyContains)
    assert.Empty(t, recorder.Result().Cookies())
    assert.Contains(t, recorder.Result().Header.Get("Content-Type"), "application/problem+json")
  })
}

func TestDraftsHandler_Share(t *testing.T) {
  const (
    routine = "Share"
    method  = http.MethodPost
    target  = "/archive.drafts.share"
    link    = "/link/to/draft"
  )

  request := httptest.NewRequest(method, target, nil)
  id := uuid.NewString()

  _ = request.ParseForm()

  request.PostForm.Add("draft_uuid", id)

  t.Run("success", func(t *testing.T) {
    expectedStatusCode := http.StatusOK
    expectedBody := string(marshal(t, gin.H{"shareable_link": link}))

    s := mocks.NewDraftsService()
    s.On(routine, mock.AnythingOfType("*gin.Context"), id).Return(link, nil)

    engine := gin.Default()
    engine.POST(target, NewDraftsHandler(s).Share)

    recorder := httptest.NewRecorder()

    engine.ServeHTTP(recorder, request)

    assert.Equal(t, expectedStatusCode, recorder.Code)
    assert.Equal(t, expectedBody, recorder.Body.String())
    assert.Empty(t, recorder.Result().Cookies())
  })

  t.Run("expected problem detail", func(t *testing.T) {
    expectedStatusCode := http.StatusBadRequest
    expectBodyContains := "Expected problem detail."

    expected := &problem.Problem{}
    expected.Status(expectedStatusCode)
    expected.Detail(expectBodyContains)

    s := mocks.NewDraftsService()
    s.On(routine, mock.AnythingOfType("*gin.Context"), id).Return("about:blank", expected)

    engine := gin.Default()
    engine.POST(target, NewDraftsHandler(s).Share)

    recorder := httptest.NewRecorder()

    engine.ServeHTTP(recorder, request)

    assert.Equal(t, expectedStatusCode, recorder.Code)
    assert.Contains(t, recorder.Body.String(), expectBodyContains)
    assert.Empty(t, recorder.Result().Cookies())
    assert.Contains(t, recorder.Result().Header.Get("Content-Type"), "application/problem+json")
  })

  t.Run("unexpected error", func(t *testing.T) {
    unexpected := errors.New("unexpected error")
    expectedStatusCode := http.StatusInternalServerError
    expectBodyContains := "An unexpected error occurred while processing your request"

    s := mocks.NewDraftsService()
    s.On(routine, mock.AnythingOfType("*gin.Context"), id).Return("about:blank", unexpected)

    engine := gin.Default()
    engine.POST(target, NewDraftsHandler(s).Share)

    recorder := httptest.NewRecorder()

    engine.ServeHTTP(recorder, request)

    assert.Equal(t, expectedStatusCode, recorder.Code)
    assert.Contains(t, recorder.Body.String(), expectBodyContains)
    assert.Empty(t, recorder.Result().Cookies())
    assert.Contains(t, recorder.Result().Header.Get("Content-Type"), "application/problem+json")
  })
}

func TestDraftsHandler_Discard(t *testing.T) {
  const (
    routine = "Discard"
    method  = http.MethodPost
    target  = "/archive.drafts.discard"
  )

  request := httptest.NewRequest(method, target, nil)
  id := uuid.NewString()

  _ = request.ParseForm()

  request.PostForm.Add("draft_uuid", id)

  t.Run("success", func(t *testing.T) {
    expectedStatusCode := http.StatusNoContent

    s := mocks.NewDraftsService()
    s.On(routine, mock.AnythingOfType("*gin.Context"), id).Return(nil)

    engine := gin.Default()
    engine.POST(target, NewDraftsHandler(s).Discard)

    recorder := httptest.NewRecorder()

    engine.ServeHTTP(recorder, request)

    assert.Equal(t, expectedStatusCode, recorder.Code)
    assert.Empty(t, recorder.Body)
    assert.Empty(t, recorder.Result().Cookies())
  })

  t.Run("expected problem detail", func(t *testing.T) {
    expectedStatusCode := http.StatusBadRequest
    expectBodyContains := "Expected problem detail."

    expected := &problem.Problem{}
    expected.Status(expectedStatusCode)
    expected.Detail(expectBodyContains)

    s := mocks.NewDraftsService()
    s.On(routine, mock.AnythingOfType("*gin.Context"), id).Return(expected)

    engine := gin.Default()
    engine.POST(target, NewDraftsHandler(s).Discard)

    recorder := httptest.NewRecorder()

    engine.ServeHTTP(recorder, request)

    assert.Equal(t, expectedStatusCode, recorder.Code)
    assert.Contains(t, recorder.Body.String(), expectBodyContains)
    assert.Empty(t, recorder.Result().Cookies())
    assert.Contains(t, recorder.Result().Header.Get("Content-Type"), "application/problem+json")
  })

  t.Run("unexpected error", func(t *testing.T) {
    unexpected := errors.New("unexpected error")
    expectedStatusCode := http.StatusInternalServerError
    expectBodyContains := "An unexpected error occurred while processing your request"

    s := mocks.NewDraftsService()
    s.On(routine, mock.AnythingOfType("*gin.Context"), id).Return(unexpected)

    engine := gin.Default()
    engine.POST(target, NewDraftsHandler(s).Discard)

    recorder := httptest.NewRecorder()

    engine.ServeHTTP(recorder, request)

    assert.Equal(t, expectedStatusCode, recorder.Code)
    assert.Contains(t, recorder.Body.String(), expectBodyContains)
    assert.Empty(t, recorder.Result().Cookies())
    assert.Contains(t, recorder.Result().Header.Get("Content-Type"), "application/problem+json")
  })
}

func TestDraftsHandler_Revise(t *testing.T) {
  const (
    routine = "Revise"
    method  = http.MethodPost
    target  = "/archive.drafts.revise"
  )

  revision := &transfer.ArticleRevision{
    Title:   "Title",
    Content: "Content",
  }

  id := uuid.NewString()

  request := httptest.NewRequest(method, target, nil)
  _ = request.ParseForm()

  request.PostForm.Add("draft_uuid", id)
  request.PostForm.Add("title", revision.Title)
  request.PostForm.Add("content", revision.Content)

  t.Run("success", func(t *testing.T) {
    expectedStatusCode := http.StatusNoContent

    s := mocks.NewDraftsService()
    s.On(routine, mock.AnythingOfType("*gin.Context"), id, revision).Return(nil)

    engine := gin.Default()
    engine.POST(target, NewDraftsHandler(s).Revise)

    recorder := httptest.NewRecorder()

    engine.ServeHTTP(recorder, request)

    assert.Equal(t, expectedStatusCode, recorder.Code)
    assert.Empty(t, recorder.Body)
    assert.Empty(t, recorder.Result().Cookies())
  })

  t.Run("expected problem detail", func(t *testing.T) {
    expectedStatusCode := http.StatusBadRequest
    expectBodyContains := "Expected problem detail."

    expected := &problem.Problem{}
    expected.Status(expectedStatusCode)
    expected.Detail(expectBodyContains)

    s := mocks.NewDraftsService()
    s.On(routine, mock.AnythingOfType("*gin.Context"), id, revision).Return(expected)

    engine := gin.Default()
    engine.POST(target, NewDraftsHandler(s).Revise)

    recorder := httptest.NewRecorder()

    engine.ServeHTTP(recorder, request)

    assert.Equal(t, expectedStatusCode, recorder.Code)
    assert.Contains(t, recorder.Body.String(), expectBodyContains)
    assert.Empty(t, recorder.Result().Cookies())
    assert.Contains(t, recorder.Result().Header.Get("Content-Type"), "application/problem+json")
  })

  t.Run("unexpected error", func(t *testing.T) {
    unexpected := errors.New("unexpected error")
    expectedStatusCode := http.StatusInternalServerError
    expectBodyContains := "An unexpected error occurred while processing your request"

    s := mocks.NewDraftsService()
    s.On(routine, mock.AnythingOfType("*gin.Context"), id, revision).Return(unexpected)

    engine := gin.Default()
    engine.POST(target, NewDraftsHandler(s).Revise)

    recorder := httptest.NewRecorder()

    engine.ServeHTTP(recorder, request)

    assert.Equal(t, expectedStatusCode, recorder.Code)
    assert.Contains(t, recorder.Body.String(), expectBodyContains)
    assert.Empty(t, recorder.Result().Cookies())
    assert.Contains(t, recorder.Result().Header.Get("Content-Type"), "application/problem+json")
  })
}
