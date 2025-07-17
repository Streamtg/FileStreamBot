package routes

import (
	"fmt"
	"net/http"
	"reflect"

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

// Función principal que carga todas las rutas declaradas como métodos del struct allRoutes
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

// Ruta raíz que responde con estado del servidor
func (a *allRoutes) Root(route *Route) {
	route.Engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "FileStreamBot API is active.",
			"status":  true,
		})
	})
}

// Nueva ruta /watch/:id con template HTML y reproductor
func (a *allRoutes) Watch(route *Route) {
	route.Engine.GET("/watch/:id", func(c *gin.Context) {
		id := c.Param("id")
		hash := c.Query("hash")

		if id == "" || hash == "" {
			c.String(http.StatusBadRequest, "Missing ID or hash.")
			return
		}

		// Este es el enlace que se renderiza como fuente del reproductor
		streamURL := fmt.Sprintf("/stream/%s%s", hash, id)

		c.HTML(http.StatusOK, "watch.html", gin.H{
			"Title":     "Reproductor Multimedia",
			"FileName":  id,
			"StreamURL": streamURL,
		})
	})
}
