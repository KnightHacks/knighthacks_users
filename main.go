package main

import (
	"context"
	"github.com/KnightHacks/knighthacks_shared/auth"
	"github.com/KnightHacks/knighthacks_shared/models"
	"github.com/KnightHacks/knighthacks_shared/utils"
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

	pool, err := pgxpool.Connect(context.Background(), getEnvOrDie("DATABASE_URI"))
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
			ClientID:     getEnvOrDie("OAUTH_GITHUB_CLIENT_ID"),
			ClientSecret: getEnvOrDie("OAUTH_GITHUB_CLIENT_SECRET"),
			RedirectURL:  getEnvOrDie("OAUTH_GITHUB_REDIRECT_URL"),
			Endpoint:     github.Endpoint,
			Scopes: []string{
				"read:user",
				"user:email",
			},
		},
	}

	newAuth, err := auth.NewAuth(getEnvOrDie("JWT_SIGNING_KEY"), getEnvOrDie("AES_CIPHER"), oauthConfigMap)
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
	config := generated.Config{
		Resolvers: &graph.Resolver{
			Repository: repository.NewDatabaseRepository(pool),
			Auth:       *newAuth,
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

func getEnvOrDie(key string) string {
	env, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("You must provide the %s environmental variable\n", key)
	}
	return env
}
