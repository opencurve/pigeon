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
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/opencurve/pigeon/internal/configure"
	"github.com/opencurve/pigeon/internal/utils"
	"github.com/opencurve/pigeon/pkg/log"
	"go.uber.org/zap"
)

type InitFunc func(*configure.ServerConfigure) error

type HTTPServer struct {
	name   string
	cfg    *configure.ServerConfigure
	initer []InitFunc
	enable bool

	errorLogger  *zap.Logger
	accessLogger *zap.Logger
	engine       *gin.Engine
}

func NewHTTPServer(name ...string) *HTTPServer {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	return &HTTPServer{
		name:   utils.FirstOne(name...),
		engine: engine,
		initer: []InitFunc{},
		enable: true,
	}
}

func (s *HTTPServer) Name() string {
	return s.name
}

func (s *HTTPServer) Initer(i InitFunc) {
	s.initer = append(s.initer, i)
}

func (s *HTTPServer) Enable() bool {
	return s.enable
}

func (s *HTTPServer) Logger() *zap.Logger {
	return s.errorLogger
}

func (s *HTTPServer) initLogger() error {
	cfg := s.cfg

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

func (s *HTTPServer) createDir(dirs []string) error {
	for _, dir := range dirs {
		err := os.MkdirAll(dir, 0644)
		if err != nil {
			s.Logger().Error("create index failed", log.Field("error", err))
			return err
		}
	}
	return nil
}

func (s *HTTPServer) routePProf() {
	if !s.cfg.GetPProfEnable() {
		return
	}

	prefix := s.cfg.GetPProfPrefix()
	s.errorLogger.Info("enable pprof, prefix=" + prefix)
	group := (&s.engine.RouterGroup).Group(prefix)
	{
		group.GET("/", gin.WrapF(pprof.Index))
		group.GET("/cmdline", gin.WrapF(pprof.Cmdline))
		group.GET("/profile", gin.WrapF(pprof.Profile))
		group.POST("/symbol", gin.WrapF(pprof.Symbol))
		group.GET("/symbol", gin.WrapF(pprof.Symbol))
		group.GET("/trace", gin.WrapF(pprof.Trace))
		group.GET("/allocs", gin.WrapH(pprof.Handler("allocs")))
		group.GET("/block", gin.WrapH(pprof.Handler("block")))
		group.GET("/goroutine", gin.WrapH(pprof.Handler("goroutine")))
		group.GET("/heap", gin.WrapH(pprof.Handler("heap")))
		group.GET("/mutex", gin.WrapH(pprof.Handler("mutex")))
		group.GET("/threadcreate", gin.WrapH(pprof.Handler("threadcreate")))
	}
}

func (s *HTTPServer) Init(cfg *configure.Configure) error {
	s.cfg = cfg.GetDefaultServer()
	if len(s.name) > 0 {
		s.cfg = cfg.GetServer(s.name)
	}

	if s.cfg == nil {
		s.enable = false
		return nil
	}

	// init logger
	err := s.initLogger()
	if err != nil {
		return err
	}

	// create directory
	dirs := []string{
		s.cfg.GetIndex(),
		path.Dir(cfg.GetPidFile()),
		path.Dir(s.cfg.GetErrorLogPath()),
		path.Dir(s.cfg.GetAccessLogPath()),
	}
	err = s.createDir(dirs)
	if err != nil {
		return err
	}

	// configure server
	s.engine.MaxMultipartMemory = s.cfg.GetMultipartMaxMemory()
	os.Setenv("TMPDIR", s.cfg.GetMultipartTempPath())

	// add router to pprof
	s.routePProf()

	// invoke initer
	for _, fn := range s.initer {
		err = fn(s.cfg)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *HTTPServer) Server() *http.Server {
	s.Logger().Info(fmt.Sprintf("ready to start server %s: %s",
		s.Name(), s.cfg.GetListenAddress()))
	return &http.Server{
		Addr:    s.cfg.GetListenAddress(),
		Handler: s.engine,
	}
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
	s.engine.NoRoute(router.wrapHandlers(handlers))
}
