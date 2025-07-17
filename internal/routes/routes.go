package routes

import (
	"EverythingSuckz/fsb/internal/types"
	"EverythingSuckz/fsb/internal/utils"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Route struct {
	Name   string
	Engine *gin.Engine
}

func (r *Route) Init(engine *gin.Engine) {
	r.Engine = engine
}

type allRoutes struct {
	log *zap.Logger
}

// Load is responsible for dynamically loading all methods in allRoutes struct
func Load(log *zap.Logger, r *gin.Engine) {
	log = log.Named("routes")
	defer log.Sugar().Info("Loaded all API Routes")
	route := &Route{Name: "/", Engine: r}
	route.Init(r)
	Type := reflect.TypeOf(&allRoutes{log})
	Value := reflect.ValueOf(&allRoutes{log})
	for i := 0; i < Type.NumMethod(); i++ {
		Type.Method(i).Func.Call([]reflect.Value{Value, reflect.ValueOf(route)})
	}
}

// Root route - basic health check
func (a *allRoutes) Root(r *Route) {
	a.log.Debug("Loading root route")
	r.Engine.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, types.RootResponse{
			Message: "Server is running.",
			Ok:      true,
			Uptime:  utils.TimeFormat(uint64(time.Since(utils.StartTime).Seconds())),
			Version: utils.VersionString,
		})
	})
}

// StreamRoute serves the media file stream (existing route in original bot)
func (a *allRoutes) Stream(r *Route) {
	a.log.Debug("Loading stream route")
	r.Engine.GET("/stream/:messageID", func(c *gin.Context) {
		messageID := c.Param("messageID")
		filePath := filepath.Join("downloads", messageID)
		c.File(filePath)
	})
	a.log.Info("Loaded stream route")
}

// StreamPlayer renders HTML player for the given messageID
func (a *allRoutes) StreamPlayer(r *Route) {
	a.log.Debug("Loading player template route")

	// Ruta nueva, sin conflicto
	r.Engine.GET("/player/:messageID", func(c *gin.Context) {
		messageID := c.Param("messageID")

		tmpl, err := template.ParseFiles("templates/player.html")
		if err != nil {
			a.log.Error("Template parsing failed", zap.Error(err))
			c.String(http.StatusInternalServerError, "Template error")
			return
		}

		// Datos que se pasan al HTML
		data := gin.H{
			"StreamURL": fmt.Sprintf("/stream/%s", messageID),
			"MessageID": messageID,
		}

		err = tmpl.Execute(c.Writer, data)
		if err != nil {
			a.log.Error("Template execution failed", zap.Error(err))
		}
	})
	a.log.Info("Loaded player route")
}
