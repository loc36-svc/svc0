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

func AcceptConn (key rxlib.Key) {
	rexaKey = key
	key.NowRunning ()
	defer key.IndicateShutdown ()	

	//
}

func coordinateServing (req http.Request, res *http.ResponseWriter) {
	//
}

var (
	conf *lib.Conf
	sdDB *sql.DB
	rexaKey rxlib.Key
)

func init () {
	if initReport != nil {
		return
	}

	// ..1.. {
	var errY error
	conf, errY = lib.Conf_New ()
	if errY != nil {
		initReport = err.New ("Unable to load service configuration.", nil, nil, errY)
		return
	}
	// ..1.. }

	// ..1.. {
	connURLFormat := "%s:%s@tcp(%s:%s)/service_addr?tls=skip-verify&serverPubKey=%s&timeout=%ss&writeTimeout=%ss&readTimeout=%ss"
	connURL := fmt.Sprintf (connURLFormat, url.QueryEscape (conf.Get ("sds.username")), url.QueryEscape (conf.Get ("sds.pass")), url.QueryEscape (conf.Get ("sds.addr")), url.QueryEscape (conf.Get ("sds.port")), url.QueryEscape (conf.Get ("sds.pub_key_file")), url.QueryEscape (conf.Get ("sds.conn_timeout")), url.QueryEscape (conf.Get ("sds.write_timeout")), url.QueryEscape (conf.Get ("sds.read_timeout")))
	
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
