package server


func (s *Server) registerRoutes() {


	// health
	s.Router.GET(
		"/health",
		s.HealthHandler,
	)



	// camera

	s.Router.POST(
		"/api/camera/open",
		s.OpenCamera,
	)

	s.Router.POST(
		"/api/camera/close",
		s.CloseCamera,
	)

	s.Router.GET(
		"/api/camera/check",
		s.CheckCamera,
	)

	s.Router.GET(
		"/api/camera/frame",
		s.CameraFrame,
	)

	s.Router.GET(
		"/api/camera/live",
		s.CameraLive,
	)



	// inference control

	s.Router.POST(
		"/api/inference/start",
		s.StartInference,
	)


	s.Router.POST(
		"/api/inference/stop",
		s.StopInference,
	)



	// object detection

	s.Router.GET(
		"/api/infer_od/live",
		s.InferODLive,
	)


	s.Router.GET(
		"/api/infer_od/live_result",
		s.InferODLiveResult,
	)

}
