// Copyright 2017 Emir Ribic. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// cryptocom - Go(lang) restful starter kit
//
// API Docs for cryptocom v1
//
//		 Terms Of Service:  N/A
//	    Schemes: http
//	    Version: 2.0.0
//	    License: MIT http://opensource.org/licenses/MIT
//	    Contact: Emir Ribic <ribice@gmail.com> https://ribice.ba
//	    Host: localhost:8080
//
//	    Consumes:
//	    - application/json
//
//	    Produces:
//	    - application/json
//
//	    Security:
//	    - bearer: []
//
//	    SecurityDefinitions:
//	    bearer:
//	         type: apiKey
//	         name: Authorization
//	         in: header
//
// swagger:meta
package api

import (
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/sappy5678/cryptocom/pkg/api/wallet"
	wl "github.com/sappy5678/cryptocom/pkg/api/wallet/logging"
	wt "github.com/sappy5678/cryptocom/pkg/api/wallet/transport"
	"github.com/sappy5678/cryptocom/pkg/utl/config"
	"github.com/sappy5678/cryptocom/pkg/utl/postgres"
	"github.com/sappy5678/cryptocom/pkg/utl/server"
	"github.com/sappy5678/cryptocom/pkg/utl/zlog"
)

// Start starts the API service
func Start(cfg *config.Configuration) error {
	_, err := postgres.New(os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}

	log := zlog.New()

	e := server.New()
	e.Static("/swaggerui", cfg.App.SwaggerUIPath)
	v1 := e.Group("/v1")
	wt.NewHTTP(wl.New(&wallet.User{}, log), v1)

	v1.GET("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
	server.Start(e, &server.Config{
		Port:                cfg.Server.Port,
		ReadTimeoutSeconds:  cfg.Server.ReadTimeout,
		WriteTimeoutSeconds: cfg.Server.WriteTimeout,
		Debug:               cfg.Server.Debug,
	})

	return nil
}
