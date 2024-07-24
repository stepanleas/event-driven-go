package http

import (
	"fmt"
	"net/http"
	"tickets/db/read_model"
	"time"

	"github.com/labstack/echo/v4"
)

type OpsBookingController struct {
	opsReadModel read_model.OpsBookingReadModel
}

func NewOpsBookingController(opsReadModel read_model.OpsBookingReadModel) OpsBookingController {
	return OpsBookingController{opsReadModel: opsReadModel}
}

func (ctrl OpsBookingController) FindAll(c echo.Context) error {
	receiptIssueDate := c.QueryParam("receipt_issue_date")
	if receiptIssueDate != "" {
		_, err := time.Parse("2006-01-02", receiptIssueDate)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid receipt_issue_date format, expected RFC3339 date: ", err.Error())
		}
	}

	reservations, err := ctrl.opsReadModel.AllReservations(receiptIssueDate)
	if err != nil {
		return fmt.Errorf("failed to find reservations: %w", err)
	}

	return c.JSON(http.StatusOK, reservations)
}

func (ctrl OpsBookingController) FindByID(c echo.Context) error {
	reservation, err := ctrl.opsReadModel.BookingReadModel(c.Request().Context(), c.Param("id"))
	if err != nil {
		return fmt.Errorf("failed to find reservation: %w", err)
	}

	return c.JSON(http.StatusOK, reservation)
}
