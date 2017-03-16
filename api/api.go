package api

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/justinas/alice"
	"github.com/rs/rest-layer/resource"
	"github.com/rs/rest-layer/rest"
	mgo "gopkg.in/mgo.v2"

	"github.com/ultimateboy/discapi/config"
	"github.com/ultimateboy/discapi/entities"
	"github.com/ultimateboy/discapi/hooks"
	"github.com/ultimateboy/discapi/middleware"
	"github.com/ultimateboy/discapi/storage"
)

// API encapulates the requirements for serving the API
type API struct {
	// Either a pointer to the MongoDB session, or nil if in-memory
	StorageSession *mgo.Session

	// All the configurations
	Config *config.Config

	ctx context.Context
}

// NewAPI creates a new API with the provided configurations
func NewAPI(cfg *config.Config) (*API, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	session, err := storage.NewSession(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage session: %v", err)
	}

	return &API{
		StorageSession: session,
		Config:         cfg,
		ctx:            ctx,
	}, nil
}

// Start builds the resource index and starts serving the api
func (a API) Start() error {
	index := resource.NewIndex()

	// User Resource
	userHandler := storage.NewHandler(a.ctx, a.StorageSession, a.Config, "users")
	userResource := index.Bind("users", entities.Users, userHandler, resource.Conf{
		AllowedModes: resource.ReadWrite,
	})
	err := userResource.Use(hooks.UserHook{
		UserResource: userResource,
	})
	if err != nil {
		return fmt.Errorf("failed to attach user hook to user resource: %v", err)
	}

	// Auth Resource
	authHandler := storage.NewHandler(a.ctx, a.StorageSession, a.Config, "auth")
	authResource := index.Bind("auth", entities.Auth, authHandler, resource.Conf{
		AllowedModes: resource.WriteOnly,
	})
	err = authResource.Use(hooks.AuthLogin{
		SigningKey:   []byte(a.Config.JWTSigningKey),
		UserResource: userResource,
	})
	if err != nil {
		return fmt.Errorf("failed to attach hook to auth resource: %v", err)
	}

	// Company resource
	companyHandler := storage.NewHandler(a.ctx, a.StorageSession, a.Config, "companies")
	companyResource := index.Bind("companies", entities.Companies, companyHandler, resource.Conf{
		AllowedModes: resource.ReadWrite,
	})

	// Disc resource
	discHandler := storage.NewHandler(a.ctx, a.StorageSession, a.Config, "discs")
	companyResource.Bind("discs", "company", entities.Discs, discHandler, resource.Conf{
		AllowedModes: resource.ReadWrite,
	})

	jwtHandler := middleware.NewJWTHandler(userResource, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, jwt.ErrInvalidKey
		}
		return []byte(a.Config.JWTSigningKey), nil
	})

	chain := alice.New()
	chain = chain.Append(jwtHandler)

	// Create and register the http handler for the API
	api, err := rest.NewHandler(index)
	if err != nil {
		return fmt.Errorf("invalid API configuration: %v", err)
	}
	http.Handle("/", chain.Then(api))

	// Create a simple healthz endpoint
	// @todo extend this to actually verify health of the api
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err = io.WriteString(w, `{"status": "ok"}`)
		if err != nil {
			log.Fatalf("Failed to write status to healthz endpoint's io writer: %s\n", err)
		}
	})

	// Serve the API on the configured port
	log.Printf("Serving API on http://localhost:%d\n", a.Config.APIPort)
	if err = http.ListenAndServe(fmt.Sprintf(":%d", a.Config.APIPort), nil); err != nil {
		return fmt.Errorf("failed to serve the API: %v", err)
	}

	return nil
}
