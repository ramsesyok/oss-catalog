//go:generate go tool oapi-codegen -config cfg.yaml internal/api/openapi.yaml
package main

import (
	"flag"
	"log"
	"net"
	"runtime"

	gen "github.com/ramsesyok/oss-catalog/internal/api/gen"
	"github.com/ramsesyok/oss-catalog/internal/api/handler"
	"github.com/ramsesyok/oss-catalog/internal/config"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	middleware "github.com/oapi-codegen/echo-middleware"
)

const serviceName = "oss-catalog"

func runServer(host, port string) error {
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
	return e.Start(net.JoinHostPort(host, port))
}

func main() {
	cfgPath := flag.String("config", "", "config file path")
	svcFlag := flag.String("service", "", "windows service control (install|uninstall)")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

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
			if err := runService(serviceName, cfg.Server.Host, cfg.Server.Port); err != nil {
				log.Fatalf("service run failed: %v", err)
			}
			return
		}
	}

	if err := runServer(cfg.Server.Host, cfg.Server.Port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
