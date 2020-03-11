package main

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/syariatifaris/anymind/app"
	"github.com/syariatifaris/anymind/app/rest"
	"github.com/syariatifaris/anymind/internal/repo"
	"github.com/syariatifaris/anymind/internal/service"
	"github.com/syariatifaris/anymind/pkg/database"
)
var(
	useCase *app.UseCase
)

func init(){
	db, err := database.NewSimpleMongoClient(context.Background(), "mongodb", 27017)
	if err != nil{
		logrus.Panicln("mongo connection failed", err.Error())
	}

	anyMindDB := db.Database("anymind")
	if anyMindDB == nil{
		logrus.Panicln("unable to get database")
	}

	defLocation, err := time.LoadLocation("Asia/Jakarta")
	if err != nil{
		logrus.Panicln("unable to load default location", err.Error())
	}

	trxDAL := repo.NewTransactionRepository(anyMindDB)
	trxService := service.NewTransactionService(trxDAL, defLocation)

	useCase = &app.UseCase{TransactionService:trxService}
}
func main(){
	restApplication := rest.NewRestApplication(useCase, 9091)
	logrus.Println("starting http application..")
	if err := restApplication.Start(); err != nil{
		logrus.Fatalln("error while serving http rest", err.Error())
	}
}