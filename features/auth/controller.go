package auth

import (
	"net/http"

	"github.com/sedind/flow/features/auth/googleauth"

	"github.com/sedind/flow/app"
	"github.com/sedind/flow/app/router"
	"github.com/sedind/flow/features/auth/jwtauth"
)

// Controller -
type Controller struct {
	*app.Context
	tokenAuth *jwtauth.JWTAuth
	gAuth     *googleauth.GoogleAuth
}

// New creates controller instance
func New(ctx *app.Context) *Controller {
	secret := ctx.AppSetting("jwt_secret")
	if secret == "" {
		ctx.Logger.Warn("jwt_secret key not provided in app_settings")
	}

	gAuthID := ctx.AppSetting("google_auth_id")
	gAuthSecret := ctx.AppSetting("google_auth_secret")
	gAuthRedirectURL := ctx.AppSetting("google_auth_redirect_url")
	if gAuthID == "" {
		ctx.Logger.Warn("google_auth_id key not provided in app_settings")
	}

	if gAuthSecret == "" {
		ctx.Logger.Warn("google_auth_secret key not provided in app_settings")
	}
	if gAuthRedirectURL == "" {
		ctx.Logger.Warn("google_auth_redirect_url key not provided in app_settings")
	}

	return &Controller{
		Context:   ctx,
		tokenAuth: jwtauth.New("HS256", []byte(secret), nil),
		gAuth:     googleauth.New(gAuthID, gAuthSecret, gAuthRedirectURL),
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
	code := r.URL.Query().Get("code")
	if len(code) > 0 {
		user, err := ctrl.gAuth.Authenticate(code)
		if err != nil {
			ctrl.JSON(w, 400, err.Error())
		}
		ctrl.JSON(w, 200, ctrl.ResponseData(user))
		return
	}
	//_, tokenString, _ := ctrl.tokenAuth.Encode(jwtauth.Claims{"user_id": 123})

	//db := ctrl.DefaultConnection().DB

	//var currentDatabase string
	//var count int

	//db.QueryRow("SELECT DATABASE()").Scan(&currentDatabase)

	//db.QueryRow("SELECT count(*) FROM INFORMATION_SCHEMA.TABLES WHERE table_schema = ? AND table_name = ?", currentDatabase, "test").Scan(&count)

	ctrl.JSON(w, 200, ctrl.ResponseData(ctrl.gAuth.LoginURL("")))
}
