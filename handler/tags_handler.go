package handler

import (
  "fontseca.dev/problem"
  "fontseca.dev/service"
  "fontseca.dev/transfer"
  "github.com/gin-gonic/gin"
  "net/http"
)

type TagsHandler struct {
  tags service.TagsService
}

func NewTagsHandler(tags service.TagsService) *TagsHandler {
  return &TagsHandler{tags: tags}
}

func (h *TagsHandler) Add(c *gin.Context) {
  var creation transfer.TagCreation

  if err := bindPostForm(c, &creation); check(err, c.Writer) {
    return
  }

  if err := validateStruct(&creation); check(err, c.Writer) {
    return
  }

  if err := h.tags.Add(c, &creation); check(err, c.Writer) {
    return
  }

  c.Status(http.StatusCreated)
}

func (h *TagsHandler) Get(c *gin.Context) {
  tags, err := h.tags.Get(c)

  if check(err, c.Writer) {
    return
  }

  c.JSON(http.StatusOK, tags)
}

func (h *TagsHandler) Update(c *gin.Context) {
  var update transfer.TagUpdate

  tag, ok := c.GetPostForm("tag_id")

  if !ok {
    problem.NewMissingParameter("tag_id").Emit(c.Writer)
    return
  }

  if err := bindPostForm(c, &update); check(err, c.Writer) {
    return
  }

  if err := validateStruct(&update); check(err, c.Writer) {
    return
  }

  if err := h.tags.Update(c, tag, &update); check(err, c.Writer) {
    return
  }

  c.Status(http.StatusNoContent)
}

func (h *TagsHandler) Remove(c *gin.Context) {
  tag, ok := c.GetPostForm("tag_id")

  if !ok {
    problem.NewMissingParameter("tag_id").Emit(c.Writer)
    return
  }

  if err := h.tags.Remove(c, tag); check(err, c.Writer) {
    return
  }

  c.Status(http.StatusNoContent)
}
