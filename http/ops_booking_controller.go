package http

import (
	"fmt"
	"net/http"
	"tickets/db/read_model"

	"github.com/labstack/echo/v4"
)

type OpsBookingController struct {
	opsReadModel read_model.OpsBookingReadModel
}

func NewOpsBookingController(opsReadModel read_model.OpsBookingReadModel) OpsBookingController {
	return OpsBookingController{opsReadModel: opsReadModel}
}

func (ctrl OpsBookingController) FindAll(c echo.Context) error {
	reservations, err := ctrl.opsReadModel.AllReservations()
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
