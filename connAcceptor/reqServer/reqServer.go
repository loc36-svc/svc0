package reqServer

import (
	"database/sql"
	"fmt"
	"gopkg.in/qamarian-dtp/err.v0" // v0.4.0
	"gopkg.in/qamarian-mmp/rxlib.v0" // v0,2.0
	"../../lib"
)

var (
	rexaKey rxlib.Key
	db *sql.DB
)

func init () {
	// ..1.. {
	if initReport != nil {
		return
	}

	if errX := lib.InitReport (); errX != nil {
		initReport = err.New (`Package "../../lib" init failed.`, nil, nil, errX)
		return
	}
	// ..1.. }

	// ..1.. {
	connURLFormat := "%s:%s@tcp(%s:%s)/service_addr?tls=skip-verify&serverPubKey=%s&timeout=%ss&writeTimeout=%ss&readTimeout=%ss"
	conf, errY := lib.Conf_New ()
	if errY != nil {
		initReport = err.New ("Unable to load service configuration.", nil, nil, errY)
		return
	}
	connURL := fmt.Sprintf (connURLFormat, url.QueryEscape ((*conf) ["dbms.username"]), url.QueryEscape ((*conf) ["dbms.pass"]), url.QueryEscape ((*conf) ["dbms.addr"]), url.QueryEscape ((*conf) ["dbms.port"]), url.QueryEscape ((*conf) ["dbms.pub_key_file"]), url.QueryEscape ((*conf) ["dbms.conn_timeout"]), url.QueryEscape ((*conf) ["dbms.write_timeout"]), url.QueryEscape ((*conf) ["dbms.read_timeout"]))
	// ..1.. }

	// ..1.. {
	var errZ error
	db, errZ = sql.Open ("mysql", connURL)
	if errZ != nil {
		initReport = err.New ("Database unreachable.", nil, nil, errZ)
		return
	}
	errA := db.Ping ()
	if errA != nil {
		initReport = err.New ("Database unreachable.", nil, nil, errA)
		return
	}
	// ..1.. }
}
