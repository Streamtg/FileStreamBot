package routes

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
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

// 🆕 Ruta para ver videos directamente
func (a *allRoutes) WatchRoute(r *Route) {
	r.Engine.GET("/watch/:file_id", func(c *gin.Context) {
		fileID := c.Param("file_id")

		// Simulación de datos del archivo (puedes reemplazar por lógica real)
		fileData := struct {
			FileName string
			FileURL  string
			FileSize string
		}{
			FileName: fmt.Sprintf("%s.mp4", fileID),
			FileURL:  fmt.Sprintf("http://%s/%s/%s.mp4", c.Request.Host, fileID, fileID),
			FileSize: "2.6 MiB",
		}

		tmplPath := filepath.Join("templates", "template.html")
		tmpl, err := template.ParseFiles(tmplPath)
		if err != nil {
			c.String(http.StatusInternalServerError, "Template error: %v", err)
			return
		}

		c.Status(http.StatusOK)
		err = tmpl.Execute(c.Writer, fileData)
		if err != nil {
			a.log.Error("Failed to render HTML template", zap.Error(err))
		}
	})
}
