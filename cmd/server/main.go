package main

import (
	"log"

	"yolo-go-inference/internal/app"
	"yolo-go-inference/internal/config"
)


func main(){

	cfg,err :=
		config.Load(
			"./configs/config.yaml",
		)

	if err!=nil{
		log.Fatal(err)
	}


	application,err :=
		app.New(cfg)

	if err!=nil{
		log.Fatal(err)
	}


	defer application.Close()


	if err:=application.Run();err!=nil{
		log.Fatal(err)
	}
}
