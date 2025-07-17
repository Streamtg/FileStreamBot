package routes

import (
	"EverythingSuckz/fsb/internal/types"
	"EverythingSuckz/fsb/internal/utils"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"reflect"
	"strings"
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

// Cargar rutas automáticamente
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

// Ruta principal "/"
func (a *allRoutes) Root(route *Route) {
	startTime := time.Now()

	route.Engine.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, types.RootResponse{
			Message: "Server is running.",
			Ok:      true,
			Uptime:  utils.TimeFormat(uint64(time.Since(startTime).Seconds())),
			Version: "3.1.0",
		})
	})
}

// Ruta para renderizar contenido con template HTML y reproductor
func (a *allRoutes) Watch(route *Route) {
	route.Engine.GET("/watch/:id", func(ctx *gin.Context) {
		fileID := ctx.Param("id")
		if fileID == "" {
			ctx.String(http.StatusBadRequest, "Missing file ID")
			return
		}

		// Simulación de obtención de archivo (reemplaza por tu lógica real)
		dummyFile := struct {
			FileName string
			MimeType string
			StreamURL string
		}{
			FileName:  "sample_video.mp4",
			MimeType:  "video/mp4",
			StreamURL: fmt.Sprintf("https://your-ngrok-url.ngrok-free.app/stream/%s", fileID),
		}

		// Detectar si es video o audio
		tag := "video"
		if strings.HasPrefix(dummyFile.MimeType, "audio") {
			tag = "audio"
		}

		// Cargar plantilla
		tmplPath := filepath.Join("templates", "watch.html")
		tmpl, err := template.ParseFiles(tmplPath)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Template error: %v", err)
			return
		}

		ctx.Header("Content-Type", "text/html; charset=utf-8")
		tmpl.Execute(ctx.Writer, gin.H{
			"Tag":      tag,
			"Title":    "Watch " + dummyFile.FileName,
			"FileName": dummyFile.FileName,
			"Src":      dummyFile.StreamURL,
		})
	})
}
