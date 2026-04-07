package main

import (
	"log"
	"os"

	"github.com/bobllor/cloud-project/src/api"
	"github.com/bobllor/cloud-project/src/config"
	dbgateway "github.com/bobllor/cloud-project/src/db_gateway"
	"github.com/bobllor/cloud-project/src/server"
	"github.com/bobllor/gologger"
)

func main() {
	// TODO: add this to the config file
	addr := ":8080"
	// TODO: add real logging for prod
	logger := gologger.NewLogger(log.New(os.Stdout, "", log.Ltime|log.Ldate), gologger.Linfo)

	serv, err := server.NewServer(addr)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: add a function to look for case-insensitive and fuzzy yaml files
	scfg, err := config.NewServerConfig("config.yml")
	if err != nil {
		log.Fatal(err)
	}

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

	stdConfig := config.NewConfig(logger)

	fg := dbgateway.NewFileGateway(fdb, stdConfig)
	ug := dbgateway.NewUserGateway(udb, stdConfig)
	sg := dbgateway.NewSessionGateway(udb, stdConfig)

	gw := &dbgateway.Gateway{
		File:    fg,
		User:    ug,
		Session: sg,
	}

	api := api.NewApi(gw)

	// TODO: remove later- temp, not permanent
	serv.RegisterHandler("/", api.User.Post.RegisterUser)

	log.Fatal(serv.Start())
}
