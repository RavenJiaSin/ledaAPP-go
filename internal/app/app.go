package app

import (
	"yolo-go-inference/internal/config"
	"yolo-go-inference/internal/server"
	"yolo-go-inference/internal/stream"
)


type App struct {

	Config config.Config

	Runtime *stream.Runtime

	Server *server.Server
}



func New(cfg config.Config) (*App,error){

	runtime := stream.NewRuntime()


	srv := server.New(
		runtime,
	)


	app := &App{

		Config: cfg,

		Runtime: runtime,

		Server: srv,
	}


	if err := app.init(); err != nil {
		runtime.Stop()
		return nil,err
	}


	return app,nil
}



func (a *App) init() error {

	if err := loadModels(
		a.Runtime,
		a.Config,
	); err != nil {
		return err
	}

	return nil
}



func (a *App) Run() error {

	return a.Server.Router.Run(
		a.Config.Server.Addr,
	)

}



func (a *App) Close(){

	a.Runtime.Stop()

}
