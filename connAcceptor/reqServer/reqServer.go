package reqServer

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"gopkg.in/gorilla/mux.v1"
	"gopkg.in/qamarian-dtp/err.v0" // v0.4.0
	errLib "gopkg.in/qamarian-lib/err.v0" // v0.4.0
	"gopkg.in/qamarian-mmp/rxlib.v0" // v0,2.0
	"net/http"
	"../../lib"
	_ "gopkg.in/go-sql-driver/mysql.v1"
)
func init () {
	if initReport != nil {
		return
	}

	if errX := lib.InitReport (); errX != nil {
		initReport = err.New (`Package "../../lib" init failed.`, nil, nil, errX)
		return
	}
}

func ServeReq (req http.Request, res *http.ResponseWriter, key rxlib.Key) {
	key.NowRunning ()
	defer key.IndicateShutdown ()	
	defer func () {
		// ..1.. {
		panicReason := recover ()
		if panicReason == nil {
			return
		}
		// ..1.. }

		errX, okX := panicReason.(err.Error)

		// ..1.. {
		if okX == false {
			errMssg := fmt.Sprintf ("Panic occured during operation. [%v]", panicReason)
			key.Send (errMssg, "log_recorder")
			output := fmt.Sprintf (responseFormat, 2, "Service error.", "[]")
			res.Write ([]byte (output))
			return
		} else if errX.Class () != nil || errX.Class () != reqDataErr {
			errY := err.New ("An error occured while serving request.", nil, nil, errX)
			errMssg := errLib.Fup (errY)
			key.Send (errMssg, "log_recorder")
			output := fmt.Sprintf (responseFormat, 2, "Service error.", "[]")
			res.Write ([]byte (output))
			return
		} else {
			errZ := err.New ("Error on client side.", nil, nil, errX)
			errMssg := errLib.Fup (errZ)
			output := fmt.Sprintf (responseFormat, 1, errMssg, "[]")
			res.Write ([]byte (output))
			return
		}
		// ..1.. }
	} ()

	// ..1.. {
	if i, _ := mux.Vars (req)["sID"); i != serviceID {
		errX := err.New ("Service requested from the wrong service.", reqDataErr, nil)
		panic (errX)
	}

	if v, _ := mux.Vars (req)["v"); v != serviceVer {
		errY := err.New ("Unsupported service version.", reqDataErr, nil)
		panic (errY)
	}
	// ..1.. }

	// ..1.. {
	errZ := db.Ping ()
	if errZ != nil {
		errA := err.New ("Database unreachable.", operationErr, nil, errZ)
		panic (errA)
	}
	// ..1.. }

	// ..1.. {
	q := `SELECT location_id, name
		FROM location`
	r, errB := db.Query (q)
	if errB != nil {
		errC := err.New ("Unable to fetch locations from database.", operationErr, nil, errB)
		panic (errC)
	}
	// ..1.. }

	// ..1.. {
	locations = []struct {ID string, Name string}{}
	for r.Next () {
		id := ""; name := ""
		errD := r.Scan (&id, &name)
		if errD != nil {
			errE := err.New ("Unable to fetch a record.", operationErr, nil, errD)
			panic (errE)
		}
		locations = append (locations, struct {ID string, Name string}{id, name})
	}
	// ..1.. }

	// ..1.. {
	jsonData, errF := json.Marshal (locations)
	if errF != nil {
		errG := err.New ("Unable to marshal locations data.", operationErr, nil, errF)
		panic (errG)
	}
	output := fmt.Sprinf (responseFormat, 0, "Success!", string (jsonData))
	res.Write ([]byte (output))

	return
}

var (
	serviceID = "0"
	serviceVer = "v0.1"

	db *sql.DB
	responseFormat string = `{
ResponseCode: %d,
Elaboration: "%s",
Data: %s`

	reqDataErr = big.NewInt (0)
	operationErr = big.NewInt (1)

	rexaKey rxlib.Key
)

func init () {
	if initReport != nil {
		return
	}

	// ..1.. {
	connURLFormat := "%s:%s@tcp(%s:%s)/state?tls=skip-verify&serverPubKey=%s&timeout=%ss&writeTimeout=%ss&readTimeout=%ss"
	conf, errY := lib.Conf_New ()
	if errY != nil {
		initReport = err.New ("Unable to load service configuration.", nil, nil, errY)
		return
	}
	connURL := fmt.Sprintf (connURLFormat, url.QueryEscape (conf.Get ("dbms.username")), url.QueryEscape (conf.Get ("dbms.pass")), url.QueryEscape (conf.Get ("dbms.addr")), url.QueryEscape (conf.Get ("dbms.port")), url.QueryEscape (conf.Get ("dbms.pub_key_file")), url.QueryEscape (conf.Get ("dbms.conn_timeout")), url.QueryEscape (conf.Get ("dbms.write_timeout")), url.QueryEscape (conf.Get ("dbms.read_timeout")))
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
