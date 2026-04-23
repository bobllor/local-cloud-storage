package main

import (
	"log"
	"os"

	"github.com/bobllor/cloud-project/src/api"
	"github.com/bobllor/cloud-project/src/config"
	dbgateway "github.com/bobllor/cloud-project/src/db_gateway"
	"github.com/bobllor/cloud-project/src/server"
	"github.com/bobllor/cloud-project/src/utils"
	"github.com/bobllor/gologger"
)

func main() {
	// TODO: add a function to look for case-insensitive and fuzzy yaml files
	scfg, err := config.NewServerConfig("config.yml")
	if err != nil {
		log.Fatal(err)
	}

	// TODO: add this to the config file
	// TODO: add real logging for prod
	logger := gologger.NewLogger(log.New(os.Stdout, "", log.Ltime|log.Ldate), gologger.Lsilent)

	err = scfg.LoadEnv()
	if err != nil {
		log.Fatal(err)
	}

	filePw := os.Getenv(config.EnvFilePwKey)
	userPw := os.Getenv(config.EnvUserPwKey)

	dbName := scfg.Database.Name
	network := scfg.Database.NetProtocol
	dbAddr := scfg.Database.Address

	fileConfig := dbgateway.NewConfig(
		scfg.Database.FileUser.User,
		filePw,
		network,
		dbAddr,
		dbName,
	)
	userConfig := dbgateway.NewConfig(
		scfg.Database.AccountUser.User,
		userPw,
		network,
		dbAddr,
		dbName,
	)

	fdb, err := dbgateway.NewDatabase(fileConfig)
	if err != nil {
		log.Fatal(err)
	}
	udb, err := dbgateway.NewDatabase(userConfig)
	if err != nil {
		log.Fatal(err)
	}

	deps := utils.NewDeps(logger)

	fg := dbgateway.NewFileGateway(fdb, deps)
	ug := dbgateway.NewUserGateway(udb, deps)
	sg := dbgateway.NewSessionGateway(udb, deps)

	gw := &dbgateway.Gateway{
		File:    fg,
		User:    ug,
		Session: sg,
	}

	ap := api.NewApiHandler(gw, logger)
	serv, err := server.NewServer(scfg.ServerAddress)
	if err != nil {
		log.Fatal(err)
	}

	serv.RegisterHandler(api.SessionGetValidateSessionRoute, ap.CreateRequestMiddleware(ap.SessionHandler.GetValidateSession))

	serv.RegisterHandler(api.UserPostRegisterRoute, ap.CreateRequestMiddleware(ap.UserHandler.PostRegisterUser))
	serv.RegisterHandler(api.UserPostLoginRoute, ap.CreateRequestMiddleware(ap.UserHandler.PostLogin))
	serv.RegisterHandler(api.UserGetUserRoute, ap.CreateAuthMiddleware(ap.UserHandler.GetUserBySessionID))

	// handles both dynamic and root based access
	serv.RegisterHandler(api.FileGetFileRootRoute, ap.CreateAuthMiddleware(ap.FileHandler.GetFiles))
	serv.RegisterHandler(api.FileGetFileParentRoute, ap.CreateAuthMiddleware(ap.FileHandler.GetFiles))

	logger.Info("Starting server")
	log.Fatal(serv.Start())
}
