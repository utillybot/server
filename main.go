package main

import (
	"encoding/gob"
	"flag"
	"github.com/antonlindstrom/pgstore"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/utillybot/server/controllers"
	"github.com/utillybot/server/discord"
	"github.com/utillybot/server/middlewares"
	"github.com/utillybot/server/redisClient"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	gob.Register(discord.TokenRequestResult{})
	gob.Register(discord.User{})
	staticPathPtr := flag.String("base-web-url", "./", "A relative or absolute path to the static web files")
	flag.Parse()

	store, err := pgstore.NewPGStore(os.Getenv("DATABASE_URL"), []byte(os.Getenv("SESSION_KEY")))
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer store.Close()
	defer store.StopCleanup(store.Cleanup(time.Minute * 5))

	go redisClient.StartRedis()

	router := chi.NewRouter()
	router.Use(middlewares.RemoveTrailingSlash)
	router.Use(middlewares.Sessions(store))
	router.Use(middlewares.ExcludeSourceMaps)

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/api", controllers.APIController)
	router.Get("/*", controllers.ReactController(*staticPathPtr))

	port := "3006"
	if envPort, empty := os.LookupEnv("PORT"); empty == true {
		port = envPort
	}
	server := &http.Server{Handler: router, Addr: ":" + port}

	log.Fatal(server.ListenAndServe())
}

func connectToDatabase() *gorm.DB {
	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})

	if err != nil {
		panic("Could not connect to database")
	}
	return db
}
