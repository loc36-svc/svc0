package connAcceptor

import (
	"database/sql"
	"gopkg.in/qamarian-mmp/rxlib.v0"
	"./reqServer"
	"../lib"
	_ "gopkg.in/go-sql-driver/mysql.v1"
)
func init () {
	if initReport != nil {
		return
	}

	if errX :=  reqServer.InitReport (); errX != nil {
		initReport = err.New (`Package "./reqServer" init failed.`, nil, nil, errX)
		return
	} else if errY := lib.InitReport (); errY != nil {
		initReport = err.New (`Package "../lib" init failed.`, nil, nil, errY)
		return
	}
}

func AcceptConn (key rxlib.Key) () {
	rexaKey = key
	//
}

var (
	sdDB *sql.DB
	rexaKey rxlib.Key
)

func init () {
	if initReport != nil {
		return
	}

	// ..1.. {
	connURLFormat := "%s:%s@tcp(%s:%s)/service_addr?tls=skip-verify&serverPubKey=%s&timeout=%ss&writeTimeout=%ss&readTimeout=%ss"
	conf, errY := lib.Conf_New ()
	if errY != nil {
		initReport = err.New ("Unable to load service configuration.", nil, nil, errY)
		return
	}
	connURL := fmt.Sprintf (connURLFormat, url.QueryEscape ((*conf) ["sds.username"]), url.QueryEscape ((*conf) ["sds.pass"]), url.QueryEscape ((*conf) ["sds.addr"]), url.QueryEscape ((*conf) ["sds.port"]), url.QueryEscape ((*conf) ["sds.pub_key_file"]), url.QueryEscape ((*conf) ["sds.conn_timeout"]), url.QueryEscape ((*conf) ["sds.write_timeout"]), url.QueryEscape ((*conf) ["sds.read_timeout"]))
	// ..1.. }

	// ..1.. {
	var errZ error
	sdDB, errZ = sql.Open ("mysql", connURL)
	if errZ != nil {
		initReport = err.New ("SDS database unreachable.", nil, nil, errZ)
		return
	}
	errA := sdDB.Ping ()
	if errA != nil {
		initReport = err.New ("SDS database unreachable.", nil, nil, errA)
		return
	}
	// ..1.. }
}
