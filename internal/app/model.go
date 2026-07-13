package app

import (
	"yolo-go-inference/internal/config"
	"yolo-go-inference/internal/inference"
	"yolo-go-inference/internal/pipeline"
	"yolo-go-inference/internal/stream"
)


func loadModels(
	runtime *stream.Runtime,
	cfg config.Config,
) error {


	session,err :=
		inference.NewONNXSession(
			cfg.Model.Path,
		)


	if err!=nil{
		return err
	}


	pipe :=
		pipeline.NewPipeline(
			session,
			cfg.Model.InputSize,
			cfg.Model.NumPreds,
			cfg.Model.NumClasses,
			cfg.Model.YOLOv8Layout(),
			cfg.Model.ConfThreshold,
			cfg.Model.IouThreshold,
		)


	runtime.RegisterPipeline(
		cfg.Model.Name,
		pipe,
	)


	return nil
}
