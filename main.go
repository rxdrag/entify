package main

import (
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"rxdrag.com/entity-engine/authentication"
	"rxdrag.com/entity-engine/handler"
	"rxdrag.com/entity-engine/resolve"
	"rxdrag.com/entity-engine/schema"
)

func main() {
	h := handler.New(&handler.Config{
		SchemaResolveFn: schema.ResolveSchema,
		Pretty:          true,
	})

	http.Handle("/graphql",
		authentication.CorsMiddleware(
			authentication.AuthMiddleware(
				resolve.LoadersMiddleware(h),
			),
		),
	)
	http.HandleFunc("/subscriptions", handler.CreateSubscriptionsHandler(schema.ResolveSchema))
	fmt.Println("Running a GraphQL API server at http://localhost:8080/graphql")
	err2 := http.ListenAndServe(":8080", nil)
	if err2 != nil {
		fmt.Printf("启动失败:%s", err2)
	}

}
