package handler

import (
  "fontseca/service"
  "fontseca/transfer"
  "github.com/gin-gonic/gin"
  "net/http"
)

type ProjectsHandler struct {
  s service.ProjectsService
}

func NewProjectsHandler(service service.ProjectsService) *ProjectsHandler {
  return &ProjectsHandler{
    s: service,
  }
}

func (h *ProjectsHandler) Get(c *gin.Context) {
  var projects, err = h.s.Get(c)
  if check(err, c.Writer) {
    return
  }
  c.JSON(http.StatusOK, projects)
}

func (h *ProjectsHandler) GetArchived(c *gin.Context) {
  var projects, err = h.s.Get(c, true)
  if check(err, c.Writer) {
    return
  }
  c.JSON(http.StatusOK, projects)
}

func (h *ProjectsHandler) GetByID(c *gin.Context) {
  var id = c.Query("id")
  var project, err = h.s.GetByID(c, id)
  if check(err, c.Writer) {
    return
  }
  c.JSON(http.StatusOK, project)
}

func (h *ProjectsHandler) Add(c *gin.Context) {
  var creation = transfer.ProjectCreation{}
  if err := bindPostForm(c, &creation); check(err, c.Writer) {
    return
  }
  if err := validateStruct(&creation); check(err, c.Writer) {
    return
  }
  var insertedID, err = h.s.Add(c, &creation)
  if check(err, c.Writer) {
    return
  }
  c.JSON(http.StatusOK, gin.H{"inserted_id": insertedID})
}
