package main

import (
	"log"
	"os"

	"github.com/SawitProRecruitment/EstateService/core/usecase"
	"github.com/SawitProRecruitment/EstateService/generated"
	"github.com/SawitProRecruitment/EstateService/handler"
	"github.com/SawitProRecruitment/EstateService/storage/postgres"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoMiddleware "github.com/oapi-codegen/echo-middleware"
)

func main() {
	e := echo.New()
	swagger, err := generated.GetSwagger()
	if err != nil {
		log.Fatalf("Error loading swagger spec: %s", err)
	}

	repo := postgres.NewRepository(os.Getenv("DATABASE_URL"))
	usecase := usecase.NewEstateUsecase(repo)

	generated.RegisterHandlers(e, handler.NewServer(usecase))

	e.Use(echoMiddleware.OapiRequestValidator(swagger))
	e.Use(middleware.Logger())
	e.Logger.Fatal(e.Start(":1323"))
}
