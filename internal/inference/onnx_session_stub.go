// internal/inference/onnx_session_stub.go
//go:build !onnxruntime

package inference

import "errors"

type ONNXSession struct{}

func NewONNXSession(modelPath string) (*ONNXSession, error) {
    return nil, errors.New("onnxruntime build tag is not enabled")
}

func (o *ONNXSession) Run(input []float32, shape []int64) ([]float32, error) {
    return nil, errors.New("onnxruntime build tag is not enabled")
}

func (o *ONNXSession) Close() {}
