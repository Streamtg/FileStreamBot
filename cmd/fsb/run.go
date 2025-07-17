package main

import (
	"EverythingSuckz/fsb/config"
	"EverythingSuckz/fsb/internal/bot"
	"EverythingSuckz/fsb/internal/cache"
	"EverythingSuckz/fsb/internal/routes"
	"EverythingSuckz/fsb/internal/types"
	"EverythingSuckz/fsb/internal/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var runCmd = &cobra.Command{
	Use:                "run",
	Short:              "Run the bot with the given configuration.",
	DisableSuggestions: false,
	Run:                runApp,
}

var startTime time.Time = time.Now()

func runApp(cmd *cobra.Command, args []string) {
	utils.InitLogger(config.ValueOf.Dev)
	log := utils.Logger
	mainLogger := log.Named("Main")
	mainLogger.Info("Starting server")

	// Carga de configuración
	config.Load(log, cmd)

	// Inicializar router con nuevas rutas y template rendering
	router := getRouter(log)

	// Iniciar cliente del bot
	mainBot, err := bot.StartClient(log)
	if err != nil {
		log.Panic("Failed to start main bot", zap.Error(err))
	}

	// Inicializar caché
	cache.InitCache(log)

	// Iniciar workers
	workers, err := bot.StartWorkers(log)
	if err != nil {
		log.Panic("Failed to start workers", zap.Error(err))
		return
	}
	workers.AddDefaultClient(mainBot, mainBot.Self)

	// Iniciar bot de usuario si está configurado
	bot.StartUserBot(log)

	mainLogger.Info("Server started", zap.Int("port", config.ValueOf.Port))
	mainLogger.Info("File Stream Bot", zap.String("version", versionString))
	mainLogger.Sugar().Infof("Server is running at %s", config.ValueOf.Host)

	// Ejecutar servidor HTTP
	err = router.Run(fmt.Sprintf(":%d", config.ValueOf.Port))
	if err != nil {
		mainLogger.Sugar().Fatalln(err)
	}
}

func getRouter(log *zap.Logger) *gin.Engine {
	if config.ValueOf.Dev {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Middleware para log de errores HTTP
	router.Use(gin.ErrorLogger())

	// Cargar templates desde carpeta "templates"
	router.LoadHTMLGlob("templates/*.html")

	// Ruta base
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, types.RootResponse{
			Message: "Server is running.",
			Ok:      true,
			Uptime:  utils.TimeFormat(uint64(time.Since(startTime).Seconds())),
			Version: versionString,
		})
	})

	// Ruta del reproductor multimedia
	router.GET("/watch/:id", func(ctx *gin.Context) {
		fileID := ctx.Param("id")
		hash := ctx.Query("hash")

		if fileID == "" || hash == "" {
			ctx.String(http.StatusBadRequest, "Missing id or hash")
			return
		}

		streamURL := fmt.Sprintf("/stream/%s%s", hash, fileID)
		ctx.HTML(http.StatusOK, "watch.html", gin.H{
			"Title":     "Reproductor multimedia",
			"FileName":  fileID,
			"StreamURL": streamURL,
		})
	})

	// Cargar rutas normales del bot
	routes.Load(log, router)

	return router
}
