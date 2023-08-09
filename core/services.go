package core

func setupServices(s *Server) {
	if s.Cfg.Mocking {
		mockedServices(s)
	} else {
		realServices(s)
	}
}

func mockedServices(s *Server) {
	s.Log.Info("mocking services")

	s.Log.Info("services mocked")
}

func realServices(s *Server) {
	s.Log.Info("initializing services")

	s.Log.Info("services initialized")
}
