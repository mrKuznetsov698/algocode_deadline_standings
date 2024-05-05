package server

import (
	"algocode_deadline_standings/configs"
	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
	"log/slog"
	"slices"
	"time"
)

func RunServer(config *configs.Config) {
	updPrd := time.Duration(config.CacheTime) * time.Second
	slices.SortFunc(config.UnsolvedBorders, func(a, b *configs.UnsolvedBorder) int {
		return a.Count - b.Count
	})

	// data
	data := NewDataAccess(config)

	// cache
	store := persist.NewMemoryStore(updPrd)

	// release mode?
	if config.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	// router
	router := gin.Default()

	// I think now we don't need this... https://pkg.go.dev/github.com/gin-gonic/gin#Engine.SetTrustedProxies
	err := router.SetTrustedProxies(nil)
	if err != nil {
		slog.Error("Cant set trusted proxies: %s", err)
		panic(err)
	}

	// templates
	router.LoadHTMLGlob("templates/*.gohtml")

	// routes (static)
	router.Static("/static", "./static")

	// table routes
	router.GET("/", cache.CacheByRequestURI(store, updPrd), addHandler(mainPage, data))

	router.GET("/search/:name", cache.CacheByRequestURI(store, updPrd), addHandler(studentStats, data))

	router.GET("/stats", cache.CacheByRequestURI(store, updPrd), addHandler(allStats, data))

	err = router.Run(config.ServerAddressPort)
	if err != nil {
		slog.Error("Server down with ", "error", err)
		panic(err) // app shouldn't be able to start in this case
	}
}
