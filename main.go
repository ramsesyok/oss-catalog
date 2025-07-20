package main

import (
	//...
	"log"
	"net"

	gen "github.com/ramsesyok/oss-catalog/internal/api/gen"
	"github.com/ramsesyok/oss-catalog/internal/api/handler"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	middleware "github.com/oapi-codegen/echo-middleware"
)

func main() {
	// OASテンプレートの読み込み
	swagger, err := gen.GetSwagger()
	if err != nil {
		log.Fatalf("Error loading OAS template: %s", err)
	}
	// ホスト名での検証は行わない
	swagger.Servers = nil

	handler := handler.Handler{}

	e := echo.New()
	e.Use(echomiddleware.Logger())
	// OASテンプレートで指定したスキーマによる検証を行う
	e.Use(middleware.OapiRequestValidator(swagger))
	gen.RegisterHandlers(e, &handler)
	e.Logger.Fatal(e.Start(net.JoinHostPort("0.0.0.0", "8080")))
}
