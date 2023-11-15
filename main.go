package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/fazarmitrais/atm-simulation-stage-3/config/postgre"
	Controller "github.com/fazarmitrais/atm-simulation-stage-3/delivery/controller"
	"github.com/fazarmitrais/atm-simulation-stage-3/delivery/menu"
	accountCsv "github.com/fazarmitrais/atm-simulation-stage-3/domain/account/repository/csv"
	accountPostgreRepo "github.com/fazarmitrais/atm-simulation-stage-3/domain/account/repository/postgre"
	transactionMapdatastore "github.com/fazarmitrais/atm-simulation-stage-3/domain/transaction/repository/mapDatastore"
	"github.com/fazarmitrais/atm-simulation-stage-3/service"
	templateRender "github.com/fazarmitrais/atm-simulation-stage-3/template"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	envInit()
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{http.MethodGet, http.MethodPost},
	}))
	context := e.AcquireContext()
	fmt.Println("Welcome to the atm simulation")
	db := postgre.Connection()

	e.Static("/", "public")
	t := &templateRender.Template{
		Templates: template.Must(template.ParseGlob("./public/views/*.html")),
	}
	e.Renderer = t

	var path string
	flag.StringVar(&path, "path", "", "Path to directory")
	flag.Parse()
	acctRepo := accountPostgreRepo.NewPostgre(db)
	trxMap := transactionMapdatastore.NewMapDatastore()
	acctCsv := accountCsv.NewCSV()
	svc := service.NewService(acctRepo, acctCsv, trxMap)

	if strings.TrimSpace(path) != "" {
		importDataFromCSV(context, path, svc)

		m := menu.NewMenu(svc)
		m.Start(context)
	} else {
		Controller.New(svc).Register(e)
		e.Logger.Fatal(e.Start(":8080"))
	}
}

func importDataFromCSV(ctx echo.Context, path string, svc *service.Service) {
	fmt.Println("Importing data from csv file...")
	err := svc.Import(ctx, path)
	if err != nil {
		log.Fatalf("Error importing: %v", err)
	}
	fmt.Println("Successfully imported data from csv: ", path)
}

func envInit() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Fatalln("No .env file found")
	}
}
