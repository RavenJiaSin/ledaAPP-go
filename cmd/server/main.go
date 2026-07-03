package main

import (
	"fmt"
	"log"
	"net/http"

	"yolo-go-inference/internal/config"
	"yolo-go-inference/internal/inference"
	"yolo-go-inference/internal/pipeline"
	"yolo-go-inference/internal/server"
)


// -----------------------------
// main
// -----------------------------
func main() {
	cfg, err := config.Load("./configs/config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	session, err := inference.NewONNXSession(cfg.Model.Path)
	if err != nil {
		log.Fatalf("failed to load model: %v", err)
	}

	yolo := pipeline.NewPipeline(
		session,
		cfg.Model.InputSize,
		cfg.Model.NumPreds,
		cfg.Model.NumClasses,
		cfg.Model.YOLOv8Layout(),
		cfg.Model.ConfThreshold,
		cfg.Model.IouThreshold,
	)

	srv := server.New(yolo)

	fmt.Printf("YOLO server started at %s\n", cfg.Server.Addr)

	http.HandleFunc("/health", srv.HealthHandler)
	http.HandleFunc("/infer", srv.InferHandler)

	err = http.ListenAndServe(cfg.Server.Addr, nil)
	if err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
