package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var db = make(map[string]string)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	engine := gin.Default()

	// Ping test
	engine.GET("/ping", func(context *gin.Context) {
		context.String(http.StatusOK, "pong")
	})

	// Get user value
	engine.GET("/user/:name", func(context *gin.Context) {
		user := context.Params.ByName("name")
		value, ok := db[user]
		if ok {
			context.JSON(http.StatusOK, gin.H{"user": user, "value": value})
		} else {
			context.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
		}
	})

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := engine.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	// TODO: try to create a router group with a relative path
	// TODO: try to craete a router group with no middleware
	// TODO: try to create a full crud mock group
	// TODO: do a little more research into gin middlewares and handlerfuncs
	authorized := engine.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	}))

	// Could be done like this too
	//accounts := gin.Accounts{
	//	"foo":  "bar",
	//	"manu": "123",
	//}
	//authMiddleware := gin.BasicAuth(accounts)
	//authorized := engine.Group("/", authMiddleware)

	/* example curl for /admin with basicauth header
	   Zm9vOmJhcg== is base64("foo:bar")

		curl -X POST \
	  	http://localhost:8080/admin \
	  	-H 'authorization: Basic Zm9vOmJhcg==' \
	  	-H 'content-type: application/json' \
	  	-d '{"value":"bar"}'
	*/
	authorized.POST("admin", func(context *gin.Context) {
		// https://go.dev/ref/spec#Type_assertions
		user := context.MustGet(gin.AuthUserKey).(string) // .(string) é "type assertion" ele garante que o resultado do MustGet não é nill e é uma 'string'

		// Parse JSON
		// TODO: try to add more values to this json object
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if context.Bind(&json) == nil {
			db[user] = json.Value
			context.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})

	return engine
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	err := r.Run(":8080")

	if err != nil {
		return
	}
}
