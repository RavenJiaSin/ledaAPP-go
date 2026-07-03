package config

import (
	"os"
	"path/filepath"
	"testing"
	
	"yolo-go-inference/internal/postprocess"
)

func TestLoadConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	data := []byte(`
server:
  addr: ":9090"

model:
  name: "yolov8"
  task: "detect"
  path: "./models/yolov8.onnx"
  input_size: 640
  num_preds: 8400
  num_classes: 80
  output_layout: "channels_first"
  conf_threshold: 0.25
  iou_threshold: 0.45
`)

	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Server.Addr != ":9090" {
		t.Fatalf("expected addr :9090, got %s", cfg.Server.Addr)
	}

	if cfg.Model.Path != "./models/yolov8.onnx" {
		t.Fatalf("unexpected model path: %s", cfg.Model.Path)
	}

	if cfg.Model.NumClasses != 80 {
		t.Fatalf("expected 80 classes, got %d", cfg.Model.NumClasses)
	}
}

func TestLoadConfigDefaults(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	data := []byte(`
model:
  name: "yolov8"
  task: "detect"
  path: "./models/yolov8.onnx"
  num_classes: 80
`)

	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.Server.Addr != ":8080" {
		t.Fatalf("expected default addr :8080, got %s", cfg.Server.Addr)
	}

	if cfg.Model.InputSize != 640 {
		t.Fatalf("expected default input size 640, got %d", cfg.Model.InputSize)
	}

	if cfg.Model.NumPreds != 8400 {
		t.Fatalf("expected default num preds 8400, got %d", cfg.Model.NumPreds)
	}

	if cfg.Model.OutputLayout != "channels_first" {
		t.Fatalf("expected default output layout channels_first, got %s", cfg.Model.OutputLayout)
	}
}

func TestLoadConfigRequiresModelPath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	data := []byte(`
model:
  num_classes: 80
`)

	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestLoadConfigRejectsInvalidOutputLayout(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	data := []byte(`
model:
  path: "./models/yolov8.onnx"
  num_classes: 80
  output_layout: "wrong"
`)

	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestModelConfigYOLOv8Layout(t *testing.T) {
	cfg := ModelConfig{OutputLayout: "channels_first"}
	if cfg.YOLOv8Layout() != postprocess.YOLOv8LayoutChannelsFirst {
		t.Fatal("expected channels_first layout")
	}

	cfg = ModelConfig{OutputLayout: "preds_first"}
	if cfg.YOLOv8Layout() != postprocess.YOLOv8LayoutPredsFirst {
		t.Fatal("expected preds_first layout")
	}
}
