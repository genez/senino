package main

import (
	_ "expvar"
	"flag"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"log"
	"mime"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path"
	myMiddleware "senino/middleware"
	"time"
)

func main() {
	cwd, _ := os.Getwd()
	var basePath = flag.String("path", cwd, "base path from where XML files are taken")
	flag.Parse()

	e := echo.New()

	e.GET("/login", func(c echo.Context) error {
		return c.HTML(http.StatusOK, LoginPage)
	})
	e.POST("/login", func(c echo.Context) error {
		login := c.FormValue("login")
		log.Println("login:", login)
		password := c.FormValue("password")
		log.Println("password:", password)

		if login == "antares" && password == "vision" {
			log.Println("autorizzato")
			cookie := new(echo.Cookie)
			cookie.SetName("sessionId")
			cookie.SetValue("anafestico")
			cookie.SetExpires(time.Now().Add(2 * time.Minute))
			c.SetCookie(cookie)
			return c.String(http.StatusOK, "Autorizzato")
		} else {
			return c.NoContent(http.StatusUnauthorized)
		}
	})
	e.GET("/css/style.css", func(c echo.Context) (err error) {
		c.Response().Header().Set(echo.HeaderContentType, mime.TypeByExtension(".css"))
		c.Response().WriteHeader(http.StatusOK)
		_, err = c.Response().Write([]byte(StyleSheet))
		return
	})

	file := e.Group("/file")
	file.Use(myMiddleware.CookieAuth(func(remoteAddress string, sessionId string) bool {
		if sessionId != "" {
			return true
		}
		return false
	}))
	file.GET("/:fileName", func(c echo.Context) error {
		fileName := c.Param("fileName")
		fullName := path.Join(*basePath, fileName) + ".xml"
		log.Println("Serving XML file:", fullName, "to client:", c.Request().RemoteAddress())
		return c.File(fullName)
	})

	// Group, Middleware and Routes for /debug/* from Go's stdlib
	// GET handlers (or POST if it needs)
	dbg := e.Group("/debug")

	// expvar
	dbg.Get("/vars", func(c echo.Context) error {
		w := c.Response().(*standard.Response).ResponseWriter
		r := c.Request().(*standard.Request).Request
		if h, p := http.DefaultServeMux.Handler(r); p != "" {
			h.ServeHTTP(w, r)
			return nil
		}
		return echo.NewHTTPError(http.StatusNotFound)
	})

	dbg.Get("/pprof/*", func(c echo.Context) error {
		w := c.Response().(*standard.Response).ResponseWriter
		r := c.Request().(*standard.Request).Request
		if h, p := http.DefaultServeMux.Handler(r); p != "" {
			h.ServeHTTP(w, r)
			return nil
		}
		return echo.NewHTTPError(http.StatusNotFound)
	})

	log.Println("Starting HTTP REST server on port 8043 (forms authentication). Base path is", *basePath)
	e.Run(standard.New(":8043"))
}
