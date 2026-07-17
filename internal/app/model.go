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


	for _, modelCfg := range cfg.Models {


		session, err :=
			inference.NewONNXSession(
				modelCfg.Path,
			)


		if err != nil {
			return err
		}



		pipe :=
			pipeline.NewPipeline(
				session,
				modelCfg.InputSize,
				modelCfg.NumPreds,
				modelCfg.NumClasses,
				modelCfg.YOLOv8Layout(),
				modelCfg.ConfThreshold,
				modelCfg.IouThreshold,
			)



		runtime.RegisterPipeline(
			modelCfg.Name,
			pipe,
		)

	}


	return nil
}
