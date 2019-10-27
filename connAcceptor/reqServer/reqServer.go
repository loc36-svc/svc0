package reqServer

import (
	"database/sql"
	"gopkg.in/qamarian-mmp/rxlib.v0" // v0,2.0
	"../../lib"
)

var (
	rexaKey rxlib.Key
	db *sql.DB
)

func init () {
	if initReport != nil {
		return
	}

	connURLFormat := "%s:%s@tcp(%s:%s)/service_addr?tls=skip-verify&serverPubKey=%s&timeout=%ss&writeTimeout=%ss&readTimeout=%ss"