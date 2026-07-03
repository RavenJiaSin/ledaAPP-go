package config

import (
	"os"

	"errors"
	"gopkg.in/yaml.v3"

	"yolo-go-inference/internal/postprocess"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	Model  ModelConfig  `yaml:"model"`
}

type ServerConfig struct {
	Addr string `yaml:"addr"`
}

type ModelConfig struct {
	Name          string  `yaml:"name"`
	Task          string  `yaml:"task"`
	Path          string  `yaml:"path"`
	InputSize     int     `yaml:"input_size"`
	NumPreds      int     `yaml:"num_preds"`
	NumClasses    int     `yaml:"num_classes"`
	OutputLayout  string  `yaml:"output_layout"`
	ConfThreshold float32 `yaml:"conf_threshold"`
	IouThreshold  float32 `yaml:"iou_threshold"`
}

func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	if cfg.Server.Addr == "" {
		cfg.Server.Addr = ":8080"
	}
	if cfg.Model.InputSize == 0 {
		cfg.Model.InputSize = 640
	}
	if cfg.Model.NumPreds == 0 {
		cfg.Model.NumPreds = 8400
	}
	if cfg.Model.OutputLayout == "" {
		cfg.Model.OutputLayout = "channels_first"
	}
	if cfg.Model.ConfThreshold == 0 {
		cfg.Model.ConfThreshold = 0.25
	}
	if cfg.Model.IouThreshold == 0 {
		cfg.Model.IouThreshold = 0.45
	}
	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c Config) Validate() error {
	if c.Model.Path == "" {
		return errors.New("model path is required")
	}
	if c.Model.NumClasses <= 0 {
		return errors.New("model num_classes must be greater than 0")
	}
	switch c.Model.OutputLayout {
	case "channels_first", "preds_first":
	default:
		return errors.New("model output_layout must be channels_first or preds_first")
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
