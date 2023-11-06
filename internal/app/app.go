package app

import (
	"context"
	"log/slog"
	"net/http"

	linkRouter "RapidURL/internal/api/http/link"
	metricsMiddleware "RapidURL/internal/api/http/middleware/metrics"
	userRouter "RapidURL/internal/api/http/user"
	"RapidURL/internal/config"
	metricsStruct "RapidURL/internal/metrics"
	linkRepository "RapidURL/internal/repository/link"
	memcachedLinkRepository "RapidURL/internal/repository/link/memcached"
	postgresLinkRepository "RapidURL/internal/repository/link/postgres"
	userRepository "RapidURL/internal/repository/user/postgres"
	"RapidURL/internal/usecase/link"
	"RapidURL/internal/usecase/user"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	prometheus_middleware "github.com/slok/go-http-metrics/middleware"
	"github.com/slok/go-http-metrics/middleware/std"
)

type App struct {
	pool  *pgxpool.Pool
	cache *memcache.Client
	log   *slog.Logger
	cfg   *config.Config
	srv   *http.Server
}

func New(pool *pgxpool.Pool, cache *memcache.Client, log *slog.Logger, cfg *config.Config) *App {
	return &App{
		pool:  pool,
		cache: cache,
		log:   log,
		cfg:   cfg,
	}
}

// Start init app and start http server in goroutine
func (a *App) Start() {
	reg := prometheus.NewRegistry()
	metrics := metricsStruct.New(reg)

	r := a.initRouter(metrics)
	a.initLink(r, metrics)
	a.initUser(r)
	a.log.Info("Starting server...")
	a.startServer(r, reg)
}

func (a *App) initLink(r *chi.Mux, metricsStruct *metricsStruct.Metrics) {
	linkStorage := postgresLinkRepository.New(a.pool, a.log)
	linkCache := memcachedLinkRepository.New(a.cache)
	cachedLink := linkRepository.NewCachedRepository(linkStorage, linkCache, a.log, metricsStruct.TotalCacheMisses)
	linkUsecase := link.New(cachedLink)
	linkRouter.Register(r, linkUsecase, a.log)
}

func (a *App) initUser(r *chi.Mux) {
	userStorage := userRepository.New(a.pool, a.log)
	userUsecase := user.New(userStorage)
	userRouter.Register(r, userUsecase, a.log)
}

func (a *App) initRouter(metricsStruct *metricsStruct.Metrics) *chi.Mux {
	r := chi.NewRouter()
	mdlw := prometheus_middleware.New(prometheus_middleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})
	r.Use(std.HandlerProvider("", mdlw))
	r.Use(metricsMiddleware.New(metricsStruct.TotalRequests, metricsStruct.ResponseTime))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	return r
}

func (a *App) startServer(r *chi.Mux, reg prometheus.Gatherer) {
	srv := &http.Server{
		Addr:         a.cfg.HTTPServer.Address,
		Handler:      r,
		ReadTimeout:  a.cfg.HTTPServer.Timeout,
		WriteTimeout: a.cfg.HTTPServer.Timeout,
		IdleTimeout:  a.cfg.HTTPServer.Timeout,
	}
	a.srv = srv

	go func() {
		err := http.ListenAndServe(":9102", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
		if err != nil {
			a.log.Error("prom", err)
		}
	}()
	go srv.ListenAndServe()
}

func (a *App) Stop(ctx context.Context) {
	a.srv.Shutdown(ctx)
}
