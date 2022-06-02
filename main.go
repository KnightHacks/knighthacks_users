package main

import (
	"context"
	"github.com/KnightHacks/knighthacks_shared/auth"
	"github.com/KnightHacks/knighthacks_users/repository"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"log"
	"net/http"
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

	oauthConfigMap := map[auth.Provider]oauth2.Config{
		//TODO: implement gmail auth, github is priority
		//auth.GmailAuthProvider: {
		//	ClientID:     "",
		//	ClientSecret: "",
		//	Endpoint:     oauth2.Endpoint{},
		//	RedirectURL:  "",
		//	Scopes:       nil,
		//},
		auth.GitHubAuthProvider: {
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

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: &graph.Resolver{
			Repository: repository.NewDatabaseRepository(pool),
			Auth:       auth.Auth{ConfigMap: oauthConfigMap},
		},
	}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getEnvOrDie(key string) string {
	env, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("You must provide the %s environmental variable\n", env)
	}
	return env
}
