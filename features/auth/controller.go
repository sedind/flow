package auth

import (
	"net/http"

	"github.com/sedind/flow/app"
	"github.com/sedind/flow/app/router"
)

// Controller -
type Controller struct {
	*app.Context
}

// New creates controller instance
func New(ctx *app.Context) *Controller {
	return &Controller{ctx}
}

// Routes returns list of defined feature routes
func (ctrl *Controller) Routes() *router.Mux {
	cr := router.New()
	cr.Get("/", ctrl.HomeGetAction)

	return cr
}

// HomeGetAction handles HTTP GET methos on / route
func (ctrl *Controller) HomeGetAction(w http.ResponseWriter, r *http.Request) {
	ctrl.JSON(w, 200, ctrl.ResponseData("Hello !!!"))
}
