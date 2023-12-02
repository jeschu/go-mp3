package engine

import (
	"context"
	"errors"
	health "github.com/AppsFlyer/go-sundheit"
	healthhttp "github.com/AppsFlyer/go-sundheit/http"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go-mp3/engine/controller"
	hlth "go-mp3/engine/health"
	"go-mp3/engine/middlewares"
	"log"
	"net/http"
	"time"
)

//goland:noinspection GoExportedFuncWithUnexportedType
func NewServer() *serverBuilder {
	return &serverBuilder{
		addr:            ":8080",
		managementAddr:  ":8081",
		shutdownTimeout: 5 * time.Second,
		withRequestId:   true,
	}
}

type Server struct {
	server           *http.Server
	engine           *gin.Engine
	managementServer *http.Server
	managementEngine *gin.Engine
	shutdownTimeout  time.Duration
	controllers      []Controller
}

func newServer(builder serverBuilder) (*Server, error) {
	engine := buildEngine(builder.withRequestId)
	managementEngine, err := buildManagementEngine()
	if err != nil {
		return nil, err
	}
	server := &Server{
		server:           &http.Server{Addr: builder.addr, Handler: engine},
		engine:           engine,
		managementServer: &http.Server{Addr: builder.managementAddr, Handler: managementEngine},
		managementEngine: managementEngine,
		shutdownTimeout:  builder.shutdownTimeout,
		controllers:      make([]Controller, 0),
	}
	server.engine.NoRoute(func(c *gin.Context) { controller.AbortWithCode(c, 404) })
	return server, nil

}

func buildEngine(withRequestId bool) *gin.Engine {
	var ngine = gin.New()
	ngine.Use(
		gin.Logger(),
		gin.Recovery(),
		middlewares.DefaultHeaders("", withRequestId),
		middlewares.Metrics(),
	)
	return ngine
}

func buildManagementEngine() (engine *gin.Engine, err error) {
	var h health.Health

	engine = gin.New()
	engine.Use(
		gin.Logger(),
		gin.Recovery(),
	)

	if h, err = hlth.RegisterHealthChecks(); err != nil {
		return
	}

	engine.GET("/health", gin.WrapF(healthhttp.HandleHealthJSON(h)))
	engine.GET("/metrics", gin.WrapH(promhttp.Handler()))
	return
}

func (server *Server) RegisterControllers(controllers ...Controller) {
	for _, c := range controllers {
		server.RegisterController(c)
	}
}
func (server *Server) RegisterController(controller Controller) {
	controller.RegisterRoutes(server.engine)
	server.controllers = append(server.controllers, controller)
}

func (server *Server) Start() {
	go func() { listenAndServe("web", server.server) }()
	go func() { listenAndServe("management", server.managementServer) }()
}

func (server *Server) Shutdown() {
	errs := make(chan error, 2)
	go shutdown("web", errs, server.shutdownTimeout, server.server)
	go shutdown("managment", errs, server.shutdownTimeout, server.managementServer)
	for e := 1; e <= 2; e++ {
		if err := <-errs; err != nil {
			log.Print(err)
		}
	}
	close(errs)
}

func listenAndServe(name string, server *http.Server) {
	log.Printf("%s server listening as %s", name, server.Addr)
	if err := server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
		log.Printf("unable to listen: %v", err)
	}
}

func shutdown(name string, errors chan error, timeout time.Duration, server *http.Server) {
	log.Printf("shutdown %s server at %v", name, server.Addr)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	errors <- server.Shutdown(ctx)
}
