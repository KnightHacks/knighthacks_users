package main

import (
	"context"
	"errors"
	"github.com/99designs/gqlgen/graphql"
	"github.com/KnightHacks/knighthacks_shared/auth"
	databaseUtils "github.com/KnightHacks/knighthacks_shared/database"
	"github.com/KnightHacks/knighthacks_shared/pagination"
	"github.com/KnightHacks/knighthacks_shared/utils"
	"github.com/KnightHacks/knighthacks_users/graph/model"
	"github.com/KnightHacks/knighthacks_users/repository/database"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"log"
	"os"
	"runtime/debug"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/KnightHacks/knighthacks_users/graph"
	"github.com/KnightHacks/knighthacks_users/graph/generated"
)

const defaultPort = "8080"

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	pool, err := databaseUtils.ConnectWithRetries(utils.GetEnvOrDie("DATABASE_URI"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	newAuth, err := auth.NewAuthWithEnvironment()
	if err != nil {
		log.Fatalf("An error occured when trying to create an instance of Auth: %s\n", err)
	}
	ginRouter := gin.Default()
	ginRouter.Use(auth.AuthContextMiddleware(newAuth))
	ginRouter.Use(utils.GinContextMiddleware())

	ginRouter.POST("/query", graphqlHandler(newAuth, pool))
	ginRouter.GET("/", playgroundHandler())

	log.Fatalln(ginRouter.Run(":" + port))
}

func graphqlHandler(a *auth.Auth, pool *pgxpool.Pool) gin.HandlerFunc {
	hasRoleDirective := auth.HasRoleDirective{GetUserId: func(ctx context.Context, obj interface{}) (string, error) {
		switch t := obj.(type) {
		case *model.User:
			return t.ID, nil
		default:
			// shouldn't happen, you must implement the new object with the ID field
			return "", errors.New("this shouldn't happen")
		}
	}, Queryable: pool}

	repository, err := database.NewDatabaseRepository(context.Background(), pool)
	if err != nil {
		log.Fatalf("error occured while initializing database repository err = %v\n", err)
	}
	config := generated.Config{
		Resolvers: &graph.Resolver{
			Repository: repository,
			Auth:       a,
		},
		Directives: generated.DirectiveRoot{
			HasRole:    hasRoleDirective.Direct,
			Pagination: pagination.Pagination,
		},
	}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(config))
	srv.SetErrorPresenter(func(ctx context.Context, err error) *gqlerror.Error {
		log.Println("Error presented: ", err)
		debug.PrintStack()
		return graphql.DefaultErrorPresenter(ctx, err)
	})
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
