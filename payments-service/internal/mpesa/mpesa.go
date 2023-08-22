package mpesa

import (
	"net/http"

	"github.com/Mik3y-F/order-management-system/pkg"
	"github.com/jwambugu/mpesa-golang-sdk"
)

const (
	MPESA_CONSUMER_KEY    = "MPESA_CONSUMER_KEY"
	MPESA_CONSUMER_SECRET = "MPESA_CONSUMER_SECRET"
	ENVIRONMENT           = "ENVIRONMENT"
)

type Mpesa struct {
	app *mpesa.Mpesa
}

func NewMpesaService() *Mpesa {

	consumerKey := pkg.MustGetEnv(MPESA_CONSUMER_KEY)
	consumerSecret := pkg.MustGetEnv(MPESA_CONSUMER_SECRET)

	env := pkg.MustGetEnv(ENVIRONMENT)

	var mpesaEnv mpesa.Environment
	if env == "prod" {
		mpesaEnv = mpesa.Production
	} else {
		mpesaEnv = mpesa.Sandbox
	}

	mpesaApp := mpesa.NewApp(http.DefaultClient, consumerKey, consumerSecret, mpesaEnv)

	return &Mpesa{
		app: mpesaApp,
	}
}
