package firebase

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	pkg "github.com/Mik3y-F/order-management-system/pkg"
)

const (
	//  These are the environment variable names that need to be set in order for the
	//  Firebase SDK to work.
	GOOGLE_APPLICATION_CREDENTIALS = "GOOGLE_APPLICATION_CREDENTIALS" // #nosec G101 - This is an env variable name
	GOOGLE_CLOUD_PROJECT           = "GOOGLE_CLOUD_PROJECT"           // #nosec G101 - This is an env variable name
	FIRESTORE_EMULATOR_HOST        = "FIRESTORE_EMULATOR_HOST"
	ENVIRONMENT                    = "ENVIRONMENT"
)

type FirebaseService struct {
	app *firebase.App
}

func NewFirebaseService() *FirebaseService {
	// Check preconditions
	// https://firebase.google.com/docs/admin/setup/#initialize_the_sdk_in_non-google_environments
	c := pkg.MustGetEnv(GOOGLE_APPLICATION_CREDENTIALS)
	log.Println(c)

	env := pkg.MustGetEnv(ENVIRONMENT)
	if env == "dev" || env == "test" {
		pkg.MustGetEnv(FIRESTORE_EMULATOR_HOST)
	}

	projectID := pkg.MustGetEnv(GOOGLE_CLOUD_PROJECT)

	conf := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(context.Background(), conf)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	return &FirebaseService{
		app: app,
	}
}

func (s *FirebaseService) GetApp() *firebase.App {
	return s.app
}
