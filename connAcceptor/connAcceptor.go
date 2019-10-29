package connAcceptor

import (
	"database/sql"
	sdsLib "gopkg.in/loc36_core/sds_lib.v0" // v0.1.0
	"gopkg.in/qamarian-dtp/err" // v0.4.0
	errLib "gopkg.in/qamarian-lib/err" // v0.4.0
	"gopkg.in/qamarian-mmp/rxlib.v0" // v0.2.0
	"net/http"
	"math/big"
	"sync"
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

func StartAndAcceptConn (key rxlib.Key) {
	rexaKey = key
	key.NowRunning ()
	defer key.IndicateShutdown ()	

	//
}

func coordinateServing (req http.Request, res *http.ResponseWriter) {
	defer func () {
		defer func () {recover ()} ()

		// ..1.. {
		panicReason := recover ()
		if panicReason == nil {
			return
		}
		// ..1.. }
		
		errA, okA := panicReason.(*err.Error)

		// ..1.. {
		keyLock.Lock ()
		key, _, errB := rexaKey.NewKey (nextLoggingID.String ())
		nextLoggingID.Add (nextLoggingID, big.NewInt (1))
		keyLock.Unlock ()
		if errB != nil {
			return
		}
		// ..1.. }

		// ..1.. {
		if okA == false {
			errMssg := fmt.Sprintf ("A panic  occured, while coordinating " +
				"the serving of a request. [%v]", panicReason)
			key.Send (errMssg, "logRecorder")
		} else {
			errMssg := fmt.Sprintf ("An error occured, while coordinating " +
				"the serving of a request. [%s]", err.Fup (errA))
			key.Send (errMssg, "logRecorder")
		}
		// ..1.. }

		// ..1.. {
		res.WriteHeader (htpp.StatusInternalServerError)
		res.Write ([]byte ("Fatal error (Error 500): An error occured."))
		// ..1.. }
	} ()

	keyLock.Lock ()
	reqKey, _, errX := rexaKey.NewKey (nextRequestID.String ())
	nextRequestID.Add (nextRequestID, big.NewInt (1))
	keyLock.Unlock ()
	if errX != nil {
		panic (err.New ("Unable to generate a rexa key for a request.", nil, nil,
			errX))
	}
	
	reqServer.ServeReq (req, res, reqKey)
}

var (
	conf *lib.Conf
	sdDB *sql.DB
	rexaKey rxlib.Key
	keyLock &sync.Mutex = &sync.Mutex {}
	nextRequestID *big.Int = big.NewInt (0)
	nextLoggingID *big.Int = big.NewInt (0)
)

func init () {
	if initReport != nil {
		return
	}

	// ..1.. {
	var errY error
	conf, errY = lib.Conf_New ()
	if errY != nil {
		initReport = err.New ("Unable to load service configuration.", nil, nil,
			errY)
		return
	}
	// ..1.. }

	// ..1.. {
	connURLFormat := "%s:%s@tcp(%s:%s)/service_addr?tls=skip-verify&serverPubKey=%s&" +
		"timeout=%ss&writeTimeout=%ss&readTimeout=%ss"
	connURL := fmt.Sprintf (connURLFormat, url.QueryEscape (conf.Get ("sds.username")),
		url.QueryEscape (conf.Get ("sds.pass")),
		url.QueryEscape (conf.Get ("sds.addr")),
		url.QueryEscape (conf.Get ("sds.port")),
		url.QueryEscape (conf.Get ("sds.pub_key_file")),
		url.QueryEscape (conf.Get ("sds.conn_timeout")),
		url.QueryEscape (conf.Get ("sds.write_timeout")),
		url.QueryEscape (conf.Get ("sds.read_timeout")))
	
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
