package connAcceptor

import (
	"context"
	"database/sql"
	sdsLib "gopkg.in/loc36_core/sdsLib.v0" // v0.1.0
	"gopkg.in/qamarian-dtp/cwg" // v0.1.0
	"gopkg.in/qamarian-dtp/err" // v0.4.0
	errLib "gopkg.in/qamarian-lib/err" // v0.4.0
	"gopkg.in/qamarian-lib/str.v3"
	"gopkg.in/qamarian-mmp/rxlib.v0" // v0.2.0
	"net/http"
	"math/big"
	"os"
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
	defer func () {
		defer func () {recover ()} ()
		key.Send ("")
	} ()
	defer key.IndicateShutdown ()

	// | --
	port, okX := conf.Get ("http_sever.port")
	if okX != nil {
		errMssg := "Start up failed. [Port number seems invalid.]"
		str.PrintEtr (errMssg, "err", "StartAndAcceptConn ()")
		key.StartupFailed (errMssg)
		return
	}
	conn, errY := sdDB.Conn (context.Background ())
	if errY != nil {
		errZ := err.New ("Unable to get a connection to the SDS.", nil,
			nil, errY)
		errMssg := errLib.Fup (errZ)
		str.PrintEtr (errMssg, "err", "StartAndAcceptConn ()")
		key.StartupFailed (errMssg)
		return
	}	
	// -- |

	// --1-- [
	errA := sdsLib.UpdateAddr (conf.Get ("http_sever.addr"), port,
		conf.Get ("sds.service_id"), conf.Get ("sds.update_pass", conn)
	if errA != nil {
		errB := err.New ("Unable to update my address in the SDS.", nil, nil, errA)
		errMssg := errLib.Fup (errB)
		str.PrintEtr (errMssg, "err", "StartAndAcceptConn ()")
		key.StartupFailed (errMssg)
		return
	}
	// --1-- ]

	// | --
		readTimeout, errC       := time.ParseDuration ("4m")
		readHeaderTimeout, errD := time.ParseDuration ("4m")
		writeTimeout, errE      := time.ParseDuration ("4m")
		idleTimeout, errF       := time.ParseDuration ("4m")

		if errC != nil || errD != nil || errE != nil || errF != nil {
			errMssg := "Error creating a timeout duration of the HTTP server."
			str.PrintEtr (errMssg, "err", "StartAndAcceptConn ()")
			key.StartupFailed (errMssg)
			return
		}
	// -- |

	// --1-- [
	server := &http.Server {
		ReadTimeout: readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout: idleTimeout,
		MaxHeaderBytes: 1 << 12,
	}
	server.Handler:= mux.NewRouter ().HandleFunc ("/location", coordinateServing)
	// --1-- ]

	// | --
	pvf, pbf := cwg.New (1)
	stillRunning := true
	// -- |

	// --1-- [
	go func (stillRunning *bool, k rxlib.Key, c *lib.Conf, proceedWG cwg.PbfCWG) {
		time.Sleep (time.Second * 2)
		
		if stillRunning == true {
			keyLock.Lock ()
			key.NowRunning ()
			keyLock.Unlock ()

			notfMssg := fmt.Sprintf ("HTTP interface now running: [%s]:%s " +
				"(HTTPS)", c.Get ("http_server.addr"),
				c.Get ("http_server.port"))

			str.PrintEtr (notfMssg, "std", "StartAndAcceptConn ()")
		}

		pbf.Done ()
	} (stillRunning, key, conf, pbf)
	// --1-- ]
	
	errX := s.ListenAndServeTLS (c.Get ("http_server.tls_crt"),
		c.Get ("http_server.tls_key"))

	stillRunning = false
	pvf.Wait ()
			
	if errX != nil && errX != http.ErrServerClosed {
		errY := err.New ("HTTP interface shutdown, due to an error.", nil, nil,
			errX)
		keyLock.Lock ()
		defer keyLock.Unlock ()
		key.Send (errLib.Fup (errY), "logRecorder")
		str.PrintEtr (errLib.Fup (errY), "err", "HTTPInterface ()")
	}
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
	connURLFormat := "%s:%s@tcp(%s:%s)/service_addr?timeout=%ss&writeTimeout=%ss&" +
		"readTimeout=%ss&tls=skip-verify"
	connURL := fmt.Sprintf (connURLFormat, url.QueryEscape (conf.Get ("sds.username")),
		url.QueryEscape (conf.Get ("sds.pass")),
		url.QueryEscape (conf.Get ("sds.addr")),
		url.QueryEscape (conf.Get ("sds.port")),
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
