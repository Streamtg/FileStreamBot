package routes

import (
	"fmt"
	"net/http"
	"path/filepath"
	"reflect"
	"strings"

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

// Ruta raíz
func (r *allRoutes) Root(route *Route) {
	route.Engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "FileStreamBot is running",
			"ok":      true,
		})
	})
	r.log.Info("Loaded root route")
}

// Ruta que sirve archivos desde ./downloads/
func (r *allRoutes) Stream(route *Route) {
	route.Engine.GET("/file/:id", func(c *gin.Context) {
		id := c.Param("id")
		// Directorio base de archivos
		filePath := filepath.Join("downloads", id)

		c.FileAttachment(filePath, id)
	})
	r.log.Info("Loaded file stream route")
}

// Ruta que muestra el template del reproductor
func (r *allRoutes) StreamPlayer(route *Route) {
	route.Engine.GET("/stream/:id", func(c *gin.Context) {
		fileID := c.Param("id")
		isVideo := strings.HasSuffix(fileID, ".mp4") || strings.HasSuffix(fileID, ".webm") || strings.HasSuffix(fileID, ".mov")
		isAudio := strings.HasSuffix(fileID, ".mp3") || strings.HasSuffix(fileID, ".ogg") || strings.HasSuffix(fileID, ".wav")

		fileURL := fmt.Sprintf("/file/%s", fileID)

		c.HTML(http.StatusOK, "player.html", gin.H{
			"FileID":  fileID,
			"FileURL": fileURL,
			"IsVideo": isVideo,
			"IsAudio": isAudio,
		})
	})
	r.log.Info("Loaded stream player route")
}
