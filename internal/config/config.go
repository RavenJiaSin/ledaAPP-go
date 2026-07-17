package config

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"

	"yolo-go-inference/internal/postprocess"
)


type Config struct {

	Server ServerConfig `yaml:"server"`

	Models []ModelConfig `yaml:"models"`

	Cameras []CameraConfig `yaml:"cameras"`
}



type ServerConfig struct {

	Addr string `yaml:"addr"`
}



type ModelConfig struct {

	Name string `yaml:"name"`

	Task string `yaml:"task"`

	Path string `yaml:"path"`

	InputSize int `yaml:"input_size"`

	NumPreds int `yaml:"num_preds"`

	NumClasses int `yaml:"num_classes"`

	OutputLayout string `yaml:"output_layout"`

	ConfThreshold float32 `yaml:"conf_threshold"`

	IouThreshold float32 `yaml:"iou_threshold"`
}



type CameraConfig struct {

	Name string `yaml:"name"`

	Source interface{} `yaml:"source"`
}



func Load(path string) (Config,error) {


	data,err :=
		os.ReadFile(path)


	if err != nil {
		return Config{},err
	}



	var cfg Config


	if err :=
		yaml.Unmarshal(
			data,
			&cfg,
		);
		err != nil {

		return Config{},err
	}



	if cfg.Server.Addr == "" {
		cfg.Server.Addr=":8080"
	}



	for i:=range cfg.Models {


		model :=
			&cfg.Models[i]


		if model.InputSize == 0 {
			model.InputSize=640
		}


		if model.NumPreds == 0 {
			model.NumPreds=8400
		}


		if model.OutputLayout=="" {
			model.OutputLayout="channels_first"
		}


		if model.ConfThreshold==0 {
			model.ConfThreshold=0.25
		}


		if model.IouThreshold==0 {
			model.IouThreshold=0.45
		}

	}



	if err:=cfg.Validate();err!=nil {
		return Config{},err
	}


	return cfg,nil
}



func (c Config) Validate() error {


	if len(c.Models)==0 {
		return errors.New("no models configured")
	}



	for _,model:=range c.Models {


		if model.Name=="" {
			return errors.New("model name is required")
		}


		if model.Path=="" {
			return errors.New(
				"model path is required: "+model.Name,
			)
		}


		if model.NumClasses<=0 {

			// 允許之後 cls/seg 擴充
			// detect 再檢查

			if model.Task=="detect" {

				return errors.New(
					"model num_classes must be greater than 0: "+
						model.Name,
				)

			}
		}



		switch model.OutputLayout {

		case "channels_first",
			"preds_first":

		default:

			return errors.New(
				"invalid output_layout: "+
					model.Name,
			)
		}

	}



	return nil
}



func (m ModelConfig) YOLOv8Layout() postprocess.YOLOv8Layout {

	switch m.OutputLayout {

	case "preds_first":

		return postprocess.YOLOv8LayoutPredsFirst


	default:

		return postprocess.YOLOv8LayoutChannelsFirst
	}

}
