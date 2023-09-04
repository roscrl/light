package app

//nolint:revive,staticcheck,gocritic
func (app *App) services() {
	if app.Cfg.Mocking {
		// mockedServices
	} else {
		// realServices
	}
}
