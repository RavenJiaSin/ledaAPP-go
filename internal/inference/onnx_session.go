//go:build onnxruntime

package inference

/*
#cgo windows CFLAGS: -I${SRCDIR}/../../onnxruntime-win-x64-1.26.0/include
#cgo windows LDFLAGS: -L${SRCDIR}/../../onnxruntime-win-x64-1.26.0/lib -lonnxruntime
#cgo !windows CFLAGS: -I/usr/local/include
#cgo !windows LDFLAGS: -L/usr/local/lib -lonnxruntime

#include <stdlib.h>
#include "onnxruntime_c_api.h"

#ifdef _WIN32
#include <windows.h>
#endif

const OrtApi* ort_api;

int InitOrtAPI() {
    ort_api = OrtGetApiBase()->GetApi(ORT_API_VERSION);
    return ort_api != NULL;
}

OrtStatus* OrtCreateEnvWrapper(OrtLoggingLevel level, const char* logid, OrtEnv** out) {
    return ort_api->CreateEnv(level, logid, out);
}

OrtStatus* OrtCreateSessionOptionsWrapper(OrtSessionOptions** options) {
    return ort_api->CreateSessionOptions(options);
}

OrtStatus* OrtSetIntraOpNumThreadsWrapper(OrtSessionOptions* options, int threads) {
    return ort_api->SetIntraOpNumThreads(options, threads);
}

OrtStatus* OrtCreateSessionUtf8Wrapper(const OrtEnv* env, const char* model_path, const OrtSessionOptions* options, OrtSession** out) {
#ifdef _WIN32
    int len = MultiByteToWideChar(CP_UTF8, 0, model_path, -1, NULL, 0);
    if (len == 0) {
        return NULL;
    }

    wchar_t* wide_path = (wchar_t*)calloc(len, sizeof(wchar_t));
    if (wide_path == NULL) {
        return NULL;
    }

    MultiByteToWideChar(CP_UTF8, 0, model_path, -1, wide_path, len);
    OrtStatus* status = ort_api->CreateSession(env, wide_path, options, out);
    free(wide_path);
    return status;
#else
    return ort_api->CreateSession(env, model_path, options, out);
#endif
}

OrtStatus* OrtCreateCpuMemoryInfoWrapper(OrtMemoryInfo** out) {
    return ort_api->CreateCpuMemoryInfo(OrtArenaAllocator, OrtMemTypeDefault, out);
}

OrtStatus* OrtGetDefaultAllocatorWrapper(OrtAllocator** out) {
    return ort_api->GetAllocatorWithDefaultOptions(out);
}

OrtStatus* OrtSessionGetInputNameWrapper(const OrtSession* session, size_t index, OrtAllocator* allocator, char** value) {
    return ort_api->SessionGetInputName(session, index, allocator, value);
}

OrtStatus* OrtSessionGetOutputNameWrapper(const OrtSession* session, size_t index, OrtAllocator* allocator, char** value) {
    return ort_api->SessionGetOutputName(session, index, allocator, value);
}

OrtStatus* OrtCreateTensorWithDataAsOrtValueWrapper(
    const OrtMemoryInfo* info,
    void* data,
    size_t data_len,
    const int64_t* shape,
    size_t shape_len,
    OrtValue** out
) {
    return ort_api->CreateTensorWithDataAsOrtValue(
        info,
        data,
        data_len,
        shape,
        shape_len,
        ONNX_TENSOR_ELEMENT_DATA_TYPE_FLOAT,
        out
    );
}

OrtStatus* OrtRunWrapper(
    OrtSession* session,
    const char* const* input_names,
    const OrtValue* const* inputs,
    size_t input_count,
    const char* const* output_names,
    size_t output_count,
    OrtValue** outputs
) {
    return ort_api->Run(
        session,
        NULL,
        input_names,
        inputs,
        input_count,
        output_names,
        output_count,
        outputs
    );
}

OrtStatus* OrtGetTensorMutableDataWrapper(OrtValue* value, void** out) {
    return ort_api->GetTensorMutableData(value, out);
}

OrtStatus* OrtGetTensorTypeAndShapeWrapper(const OrtValue* value, OrtTensorTypeAndShapeInfo** out) {
    return ort_api->GetTensorTypeAndShape(value, out);
}

OrtStatus* OrtGetTensorShapeElementCountWrapper(const OrtTensorTypeAndShapeInfo* info, size_t* out) {
    return ort_api->GetTensorShapeElementCount(info, out);
}

const char* OrtGetErrorMessageWrapper(const OrtStatus* status) {
    return ort_api->GetErrorMessage(status);
}

void OrtReleaseStatusWrapper(OrtStatus* status) {
    ort_api->ReleaseStatus(status);
}

void OrtReleaseValueWrapper(OrtValue* value) {
    ort_api->ReleaseValue(value);
}

void OrtReleaseTensorTypeAndShapeInfoWrapper(OrtTensorTypeAndShapeInfo* info) {
    ort_api->ReleaseTensorTypeAndShapeInfo(info);
}

void OrtReleaseSessionWrapper(OrtSession* session) {
    ort_api->ReleaseSession(session);
}

void OrtReleaseEnvWrapper(OrtEnv* env) {
    ort_api->ReleaseEnv(env);
}

void OrtReleaseMemoryInfoWrapper(OrtMemoryInfo* mem_info) {
    ort_api->ReleaseMemoryInfo(mem_info);
}

void OrtReleaseSessionOptionsWrapper(OrtSessionOptions* options) {
    ort_api->ReleaseSessionOptions(options);
}

void OrtAllocatorFreeWrapper(OrtAllocator* allocator, void* p) {
    allocator->Free(allocator, p);
}
*/
import "C"

import (
	"errors"
	"unsafe"
)

type ONNXSession struct {
	env        *C.OrtEnv
	session    *C.OrtSession
	memInfo    *C.OrtMemoryInfo
	inputName  *C.char
	outputName *C.char
}

func NewONNXSession(modelPath string) (*ONNXSession, error) {
	if C.InitOrtAPI() == 0 {
		return nil, errors.New("failed to initialize ONNX Runtime API: loaded DLL version does not match headers")
	}

	logID := C.CString("yolo")
	defer C.free(unsafe.Pointer(logID))

	var env *C.OrtEnv
	status := C.OrtCreateEnvWrapper(C.ORT_LOGGING_LEVEL_WARNING, logID, &env)
	if status != nil {
		return nil, ortError(status)
	}

	var sessionOptions *C.OrtSessionOptions
	status = C.OrtCreateSessionOptionsWrapper(&sessionOptions)
	if status != nil {
		C.OrtReleaseEnvWrapper(env)
		return nil, ortError(status)
	}
	defer C.OrtReleaseSessionOptionsWrapper(sessionOptions)

	status = C.OrtSetIntraOpNumThreadsWrapper(sessionOptions, 1)
	if status != nil {
		C.OrtReleaseEnvWrapper(env)
		return nil, ortError(status)
	}

	modelPathC := C.CString(modelPath)
	defer C.free(unsafe.Pointer(modelPathC))

	var session *C.OrtSession
	status = C.OrtCreateSessionUtf8Wrapper(env, modelPathC, sessionOptions, &session)
	if status != nil {
		C.OrtReleaseEnvWrapper(env)
		return nil, ortError(status)
	}

	var memInfo *C.OrtMemoryInfo
	status = C.OrtCreateCpuMemoryInfoWrapper(&memInfo)
	if status != nil {
		C.OrtReleaseSessionWrapper(session)
		C.OrtReleaseEnvWrapper(env)
		return nil, ortError(status)
	}

	var allocator *C.OrtAllocator
	status = C.OrtGetDefaultAllocatorWrapper(&allocator)
	if status != nil {
		C.OrtReleaseMemoryInfoWrapper(memInfo)
		C.OrtReleaseSessionWrapper(session)
		C.OrtReleaseEnvWrapper(env)
		return nil, ortError(status)
	}

	var inputName *C.char
	status = C.OrtSessionGetInputNameWrapper(session, 0, allocator, &inputName)
	if status != nil {
		C.OrtReleaseMemoryInfoWrapper(memInfo)
		C.OrtReleaseSessionWrapper(session)
		C.OrtReleaseEnvWrapper(env)
		return nil, ortError(status)
	}

	var outputName *C.char
	status = C.OrtSessionGetOutputNameWrapper(session, 0, allocator, &outputName)
	if status != nil {
		C.OrtAllocatorFreeWrapper(allocator, unsafe.Pointer(inputName))
		C.OrtReleaseMemoryInfoWrapper(memInfo)
		C.OrtReleaseSessionWrapper(session)
		C.OrtReleaseEnvWrapper(env)
		return nil, ortError(status)
	}

	return &ONNXSession{
		env:        env,
		session:    session,
		memInfo:    memInfo,
		inputName:  inputName,
		outputName: outputName,
	}, nil
}

func (o *ONNXSession) Run(input []float32, shape []int64) ([]float32, error) {
	if len(input) == 0 {
		return nil, errors.New("empty input")
	}
	if len(shape) == 0 {
		return nil, errors.New("empty input shape")
	}

	inputDataPtr := unsafe.Pointer(&input[0])

	var inputTensor *C.OrtValue
	status := C.OrtCreateTensorWithDataAsOrtValueWrapper(
		o.memInfo,
		inputDataPtr,
		C.size_t(len(input)*4),
		(*C.int64_t)(unsafe.Pointer(&shape[0])),
		C.size_t(len(shape)),
		&inputTensor,
	)
	if status != nil {
		return nil, ortError(status)
	}
	defer C.OrtReleaseValueWrapper(inputTensor)

	inputNames := []*C.char{o.inputName}
	outputNames := []*C.char{o.outputName}

	var outputTensor *C.OrtValue
	status = C.OrtRunWrapper(
		o.session,
		(**C.char)(unsafe.Pointer(&inputNames[0])),
		(**C.OrtValue)(unsafe.Pointer(&inputTensor)),
		1,
		(**C.char)(unsafe.Pointer(&outputNames[0])),
		1,
		&outputTensor,
	)
	if status != nil {
		return nil, ortError(status)
	}
	defer C.OrtReleaseValueWrapper(outputTensor)

	var outputData unsafe.Pointer
	status = C.OrtGetTensorMutableDataWrapper(outputTensor, &outputData)
	if status != nil {
		return nil, ortError(status)
	}

	var shapeInfo *C.OrtTensorTypeAndShapeInfo
	status = C.OrtGetTensorTypeAndShapeWrapper(outputTensor, &shapeInfo)
	if status != nil {
		return nil, ortError(status)
	}
	defer C.OrtReleaseTensorTypeAndShapeInfoWrapper(shapeInfo)

	var outputSize C.size_t
	status = C.OrtGetTensorShapeElementCountWrapper(shapeInfo, &outputSize)
	if status != nil {
		return nil, ortError(status)
	}

	size := int(outputSize)
	outputSlice := (*[1 << 30]float32)(outputData)[:size:size]

	result := make([]float32, size)
	copy(result, outputSlice)

	return result, nil
}

func (o *ONNXSession) Close() {
	var allocator *C.OrtAllocator
	hasAllocator := C.OrtGetDefaultAllocatorWrapper(&allocator) == nil

	if o.inputName != nil && hasAllocator {
		C.OrtAllocatorFreeWrapper(allocator, unsafe.Pointer(o.inputName))
		o.inputName = nil
	}
	if o.outputName != nil && hasAllocator {
		C.OrtAllocatorFreeWrapper(allocator, unsafe.Pointer(o.outputName))
		o.outputName = nil
	}
	if o.memInfo != nil {
		C.OrtReleaseMemoryInfoWrapper(o.memInfo)
		o.memInfo = nil
	}
	if o.session != nil {
		C.OrtReleaseSessionWrapper(o.session)
		o.session = nil
	}
	if o.env != nil {
		C.OrtReleaseEnvWrapper(o.env)
		o.env = nil
	}
}

func ortError(status *C.OrtStatus) error {
	msg := C.GoString(C.OrtGetErrorMessageWrapper(status))
	C.OrtReleaseStatusWrapper(status)
	return errors.New(msg)
}
