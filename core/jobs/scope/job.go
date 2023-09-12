package scope

import (
	"database/sql"
	"net/http"

	"github.com/roscrl/light/config"
	"github.com/roscrl/light/db/sqlc"
)

type Job struct {
	Cfg *config.App

	DB  *sql.DB
	Qry *sqlc.Queries

	Client *http.Client
}
