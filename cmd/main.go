package main

import (
	"fmt"
	"net/http"

	_ "net/http/pprof" // profiler

	"eadownloader/internal/bot"
	"eadownloader/internal/config"
	"eadownloader/internal/database"
	"eadownloader/internal/localization"
	"eadownloader/internal/logger"
	"eadownloader/internal/util"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	logger.Init()
	defer logger.L.Sync()

	localization.Init()
	config.Load()
	logger.SetLevel(config.Env.LogLevel)

	if !util.CheckFFmpeg() {
		logger.L.Fatal("ffmpeg binary not found in PATH")
	}

	if len(config.Env.Admins) > 0 {
		logger.L.Infof("admins: %v", config.Env.Admins)
	}

	if len(config.Env.Whitelist) > 0 {
		config.Env.Whitelist = append(config.Env.Whitelist, config.Env.Admins...)
		logger.L.Infof("whitelist is enabled: %v", config.Env.Whitelist)
	}

	if config.Env.ProfilerPort > 0 {
		go func() {
			port := config.Env.ProfilerPort
			logger.L.Infof("starting profiler on port %d", port)
			if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil); err != nil {
				logger.L.Fatalf("failed to start profiler: %v", err)
			}
		}()
	}

	if config.Env.MetricsPort > 0 {
		go func() {
			port := config.Env.MetricsPort
			logger.L.Infof("starting prometheus metrics on port %d", port)
			http.Handle("/metrics", promhttp.Handler())
			if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil); err != nil {
				logger.L.Fatalf("failed to start metrics server: %v", err)
			}
		}()
	}

	database.Init()
	util.CleanupDownloadsJob()

	go bot.Start()

	select {}
}
