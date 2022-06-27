package main

import (
	"context"
	"errors"
	"github.com/KnightHacks/knighthacks_shared/auth"
	"github.com/KnightHacks/knighthacks_shared/models"
	"github.com/KnightHacks/knighthacks_shared/pagination"
	"github.com/KnightHacks/knighthacks_shared/utils"
	"github.com/KnightHacks/knighthacks_users/graph/model"
	"github.com/KnightHacks/knighthacks_users/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"log"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/KnightHacks/knighthacks_users/graph"
	"github.com/KnightHacks/knighthacks_users/graph/generated"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	pool, err := pgxpool.Connect(context.Background(), utils.GetEnvOrDie("DATABASE_URI"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	oauthConfigMap := map[models.Provider]oauth2.Config{
		//TODO: implement gmail auth, github is priority
		//auth.GmailAuthProvider: {
		//	ClientID:     "",
		//	ClientSecret: "",
		//	Endpoint:     oauth2.Endpoint{},
		//	RedirectURL:  "",
		//	Scopes:       nil,
		//},
		models.ProviderGithub: {
			ClientID:     utils.GetEnvOrDie("OAUTH_GITHUB_CLIENT_ID"),
			ClientSecret: utils.GetEnvOrDie("OAUTH_GITHUB_CLIENT_SECRET"),
			RedirectURL:  utils.GetEnvOrDie("OAUTH_GITHUB_REDIRECT_URL"),
			Endpoint:     github.Endpoint,
			Scopes: []string{
				"read:user",
				"user:email",
			},
		},
	}

	newAuth, err := auth.NewAuth(utils.GetEnvOrDie("JWT_SIGNING_KEY"), utils.GetEnvOrDie("AES_CIPHER"), oauthConfigMap)
	if err != nil {
		log.Fatalf("An error occured when trying to create an instance of Auth: %s\n", err)
	}
	ginRouter := gin.Default()
	ginRouter.Use(auth.AuthContextMiddleware(newAuth))
	ginRouter.Use(utils.GinContextMiddleware())

	ginRouter.POST("/query", graphqlHandler(newAuth, pool))
	ginRouter.GET("/", playgroundHandler())

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatalln(ginRouter.Run())
}

func graphqlHandler(newAuth *auth.Auth, pool *pgxpool.Pool) gin.HandlerFunc {
	hasRoleDirective := auth.HasRoleDirective{GetUserId: func(ctx context.Context, obj interface{}) (string, error) {
		switch t := obj.(type) {
		case *model.User:
			return t.ID, nil
		default:
			// shouldn't happen, you must implement the new object with the ID field
			return "", errors.New("this shouldn't happen")
		}
	}}

	config := generated.Config{
		Resolvers: &graph.Resolver{
			Repository: repository.NewDatabaseRepository(pool),
			Auth:       *newAuth,
		},
		Directives: generated.DirectiveRoot{
			HasRole:    hasRoleDirective.Direct,
			Pagination: pagination.Pagination,
		},
	}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(config))

	return func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	}
}

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
