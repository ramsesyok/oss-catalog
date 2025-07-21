//go:generate go tool oapi-codegen -config cfg.yaml internal/api/openapi.yaml
package main

import (
	"flag"
	"log"
	"net"
	"runtime"

	gen "github.com/ramsesyok/oss-catalog/internal/api/gen"
	"github.com/ramsesyok/oss-catalog/internal/api/handler"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	middleware "github.com/oapi-codegen/echo-middleware"
)

const serviceName = "oss-catalog"

func runServer() error {
	// OASテンプレートの読み込み
	swagger, err := gen.GetSwagger()
	if err != nil {
		return err
	}
	// ホスト名での検証は行わない
	swagger.Servers = nil

	h := handler.Handler{}

	e := echo.New()
	e.Use(echomiddleware.Logger())
	// OASテンプレートで指定したスキーマによる検証を行う
	e.Use(middleware.OapiRequestValidator(swagger))
	gen.RegisterHandlers(e, &h)
	return e.Start(net.JoinHostPort("0.0.0.0", "8080"))
}

func main() {
	svcFlag := flag.String("service", "", "windows service control (install|uninstall)")
	flag.Parse()

	if runtime.GOOS == "windows" {
		switch *svcFlag {
		case "install":
			if err := installService(serviceName, "OSS Catalog service"); err != nil {
				log.Fatalf("install failed: %v", err)
			}
			return
		case "uninstall":
			if err := removeService(serviceName); err != nil {
				log.Fatalf("uninstall failed: %v", err)
			}
			return
		}

		isSvc, err := isWindowsService()
		if err == nil && isSvc {
			if err := runService(serviceName); err != nil {
				log.Fatalf("service run failed: %v", err)
			}
			return
		}
	}

	if err := runServer(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
