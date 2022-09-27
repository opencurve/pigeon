/*
 *  Copyright (c) 2022 NetEase Inc.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

/*
 * Project: Pigeon
 * Created Date: 2022-09-20
 * Author: Jingli Chen (Wine93)
 */

package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opencurve/pigeon/internal/configure"
	"github.com/opencurve/pigeon/pkg/log"
	"go.uber.org/zap"
)

type HTTPServer struct {
	name         []string
	configure    *configure.Configure
	serverCfg    *configure.ServerConfigure
	accessLogger *zap.Logger
	errorLogger  *zap.Logger

	engine *gin.Engine
}

func NewHTTPServer(name ...string) *HTTPServer {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	return &HTTPServer{
		name:   name,
		engine: engine,
	}
}

func (s *HTTPServer) initLogger() error {
	cfg := s.serverCfg

	// error logger
	logger, err := log.New(cfg.GetLogLevel(), cfg.GetErrorLogPath())
	if err != nil {
		return err
	}
	s.errorLogger = logger

	// access logger
	logger, err = log.New("info", cfg.GetAccessLogPath())
	if err != nil {
		return err
	}
	logger = logger.WithOptions(zap.WithCaller(false))
	s.accessLogger = logger

	return nil
}

func (s *HTTPServer) Init(config *configure.Configure) error {
	s.configure = config
	s.serverCfg = configure.DefaultServer()
	return s.initLogger()
}

func (s *HTTPServer) Start() error {
	server := &http.Server{
		Addr:    s.serverCfg.GetListenAddress(),
		Handler: s.engine,
	}
	return server.ListenAndServe()
}

func (s *HTTPServer) Route(relativePath string, handlers ...HandlerFunc) {
	router := &router{server: s, group: &s.engine.RouterGroup}
	router.Route(relativePath, handlers...)
}

func (s *HTTPServer) RouterGroup(relativePath string, handlers ...HandlerFunc) RouterGroup {
	router := &router{server: s, group: &s.engine.RouterGroup}
	return router.Group(relativePath, handlers...)
}

func (s *HTTPServer) DefaultRoute(handlers ...HandlerFunc) {
	router := &router{server: s}
	s.engine.NoRoute(router.warpHandlers(handlers))
}
