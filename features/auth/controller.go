package auth

import (
	"net/http"

	"github.com/sedind/flow/app"
	"github.com/sedind/flow/app/router"
	"github.com/sedind/flow/features/auth/jwtauth"
)

// Controller -
type Controller struct {
	*app.Context
	tokenAuth *jwtauth.JWTAuth
}

// New creates controller instance
func New(ctx *app.Context) *Controller {
	secret := ctx.AppSetting("jwt_secret")
	return &Controller{
		Context:   ctx,
		tokenAuth: jwtauth.New("HS256", []byte(secret), nil),
	}
}

// TokenAuth returns tokenauth instance
func (ctrl *Controller) TokenAuth() *jwtauth.JWTAuth {
	return ctrl.tokenAuth
}

// JWTVerifierHandler verifies a JWT string from a http request
func (ctrl *Controller) JWTVerifierHandler() func(http.Handler) http.Handler {
	return jwtauth.Verifier(ctrl.tokenAuth)
}

// JWTAuthenticatorMiddleware uses a default authentication middleware to enforce access from the
// Verifier middleware request context values. The Authenticator sends a 401 Unauthorized
// response for any unverified tokens and passes the good ones through. It's just fine
// until you decide to write something similar and customize your client response.
func (ctrl *Controller) JWTAuthenticatorMiddleware(next http.Handler) http.Handler {
	return jwtauth.Authenticator(next)
}

// Routes returns list of defined feature routes
func (ctrl *Controller) Routes() *router.Mux {
	cr := router.New()
	cr.Get("/", ctrl.HomeGetAction)

	return cr
}

// HomeGetAction handles HTTP GET methos on / route
func (ctrl *Controller) HomeGetAction(w http.ResponseWriter, r *http.Request) {
	_, tokenString, _ := ctrl.tokenAuth.Encode(jwtauth.Claims{"user_id": 123})
	ctrl.JSON(w, 200, ctrl.ResponseData(tokenString))
}
