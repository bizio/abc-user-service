// @title ABC User Service API
// @version 1.0
// @description This is a sample user service for ABC.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /v1
package main

import (
	"fmt"
	"os"

	_ "github.com/bizio/abc-user-service/docs"
	cmd "github.com/bizio/abc-user-service/pkg/cmd"
)

func main() {
	if err := cmd.RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
