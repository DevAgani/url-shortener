package main

import (
	"fmt"
	h "github.com/DevAgani/url-shortener/api"
	mr "github.com/DevAgani/url-shortener/repository/mongo"
	"github.com/DevAgani/url-shortener/shortener"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

const (
	MONGO_URL ="<replace with your url>"
	MONGO_TIMEOUT = "30"
	MONGO_DB="shortener"
)

func main() {
	repo := chooseRepo("mongo")
	service := shortener.NewRedirectService(repo)
	handler := h.NewHandler(service)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/{code}",handler.Get)
	r.Post("/",handler.Post)

	errs := make(chan error, 2)
	go func() {
		fmt.Println("Listening on port :8000 ")
		errs <- http.ListenAndServe(httpPort(),r)
	}()

	go func() {
		c := make(chan os.Signal,1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s",<-c)
	}()

	fmt.Printf("Terminated %s",<- errs)
}

func httpPort() string  {
	port := "8000"
	if os.Getenv("PORT") != ""{
		port = os.Getenv("PORT")
	}
	return fmt.Sprintf(":%s",port)
}

func chooseRepo(env string) shortener.RedirectRepository  {
	switch env {
	case "mongo":
		mongoURL := MONGO_URL
		mongodb := MONGO_DB
		mongoTimeout, _ := strconv.Atoi(MONGO_TIMEOUT)
		repo,err := mr.NewMongoRepository(mongoURL,mongodb,mongoTimeout)
		if err != nil{
			log.Fatal(err)
		}
		return repo
	}
	return nil
}
