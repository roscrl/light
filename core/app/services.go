package app

func setupServices(s *App) {
	if s.Cfg.Mocking {
		mockedServices(s)
	} else {
		realServices(s)
	}
}

func mockedServices(s *App) {
	s.Log.Info("mocking services")

	s.Log.Info("services mocked")
}

func realServices(s *App) {
	s.Log.Info("initializing services")

	s.Log.Info("services initialized")
}
