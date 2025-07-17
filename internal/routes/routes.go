package routes

import (
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

// Ruta raíz (puedes mantenerla si no genera conflicto con otras)
func (r *allRoutes) Root(route *Route) {
	route.Engine.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "FileStreamBot is running",
			"ok":      true,
		})
	})
	r.log.Info("Loaded root route")
}

// Ejemplo de ruta para servir un archivo desde el sistema (completar según tu implementación)
func (r *allRoutes) Stream(route *Route) {
	route.Engine.GET("/file/:id", func(c *gin.Context) {
		id := c.Param("id")
		// Aquí deberías implementar la lógica para buscar el archivo y servirlo
		// Por ahora, solo mostramos el ID
		c.String(200, "Aquí iría el archivo con ID: %s", id)
	})
	r.log.Info("Loaded stream route")
}

// Nueva ruta con template y reproductor
func (r *allRoutes) StreamPlayer(route *Route) {
	route.Engine.GET("/stream/:id", func(c *gin.Context) {
		fileID := c.Param("id")

		// Determinar si es video por extensión simple
		isVideo := strings.HasSuffix(fileID, ".mp4") || strings.HasSuffix(fileID, ".webm") || strings.HasSuffix(fileID, ".mov")

		c.HTML(200, "player.html", gin.H{
			"FileID":  fileID,
			"IsVideo": isVideo,
		})
	})
	r.log.Info("Loaded stream player route")
}
