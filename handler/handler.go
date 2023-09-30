package handler

import (
	"fmt"
	"net/http"
	"tugaspijar/repository"

	"github.com/labstack/echo/v4"
)

type simpleRequest struct {
	User string `json:"user"`
	Name string `json:"name"`
}

type simpleResponse struct {
	ID      string `json:"id"`
	User    string `json:"user"`
	Name    string `json:"name"`
	Created string `json:"create"`
	Updated string `json:"updated"`
}

type Handler struct {
	repo repository.Repository
}

func NewHandler() Handler {
	repo := repository.NewRepository()
	return Handler{repo: repo}
}

func (h Handler) Index(c echo.Context) error {
	result, err := h.repo.Simple()
	if err != nil {
		fmt.Println(err)
		return echo.ErrServiceUnavailable
	}

	res := make([]simpleResponse, 0)
	for _, v := range result {
		res = append(res, mapResponse(v))
	}

	return c.JSON(http.StatusOK, res)
}

func (h Handler) Create(c echo.Context) error {
	var req simpleRequest
	if err := c.Bind(&req); err != nil {
		fmt.Println(err)
		return echo.ErrBadRequest
	}

	result, err := h.repo.Insert(req.User, req.Name)
	if err != nil {
		fmt.Println(err)
		return echo.ErrServiceUnavailable
	}
	res := mapResponse(result)

	return c.JSON(http.StatusCreated, res)
}

func (h Handler) Update(c echo.Context) error {
	var req simpleRequest
	if err := c.Bind(&req); err != nil {
		fmt.Println(err)
		return echo.ErrBadRequest
	}

	id := c.Param("id")
	ent := repository.SimpleEntity{
		User: req.User,
		Name: req.Name,
	}
	result, err := h.repo.Update(id, ent)
	if err != nil {
		fmt.Println(err)
		return echo.ErrServiceUnavailable
	}
	res := mapResponse(result)

	return c.JSON(http.StatusCreated, res)
}

func (h Handler) SimpleOne(c echo.Context) error {
	id := c.Param("id")

	result, err := h.repo.SimpleOne(id)
	if err != nil {
		fmt.Println(err)
		return echo.ErrNotFound
	}
	res := mapResponse(result)

	return c.JSON(http.StatusCreated, res)
}

func mapResponse(ent repository.SimpleEntity) simpleResponse {
	return simpleResponse{
		ID:      ent.ID.Hex(),
		User:    ent.User,
		Name:    ent.Name,
		Created: ent.Created.String(),
		Updated: ent.Updated.String(),
	}
}