package http

import (
	"fmt"
	"net/http"
	"tickets/entities"
	"tickets/message/contracts"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type showCreateRequest struct {
	DeadNationID    uuid.UUID `json:"dead_nation_id"`
	NumberOfTickets int       `json:"number_of_tickets"`
	StartTime       time.Time `json:"start_time"`
	Title           string    `json:"title"`
	Venue           string    `json:"venue"`
}

type ShowController struct {
	repo contracts.ShowRepository
}

func NewShowController(repo contracts.ShowRepository) ShowController {
	return ShowController{repo: repo}
}

func (ctrl ShowController) Store(c echo.Context) error {
	var request showCreateRequest
	if err := c.Bind(&request); err != nil {
		return err
	}

	show := entities.Show{
		ShowID:          uuid.New(),
		DeadNationID:    request.DeadNationID,
		NumberOfTickets: request.NumberOfTickets,
		StartTime:       request.StartTime,
		Title:           request.Title,
		Venue:           request.Venue,
	}

	if err := ctrl.repo.Add(c.Request().Context(), show); err != nil {
		return fmt.Errorf("failed to store show: %w", err)
	}

	return c.JSON(http.StatusCreated, show)
}
