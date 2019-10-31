package lib

import (
	"gopkg.in/asaskevich/govalidator.v9"
	"gopkg.in/qamarian-dtp/err.v0" // v0.4.0
	viperLib "gopkg.in/qamarian-lib/viper.v0" // v0.1.0
	"gopkg.in/spf13/afero.v1"
	"net"
	"regexp"
	"strconv"
)

func Conf_New () (output *Conf, somrErr error) {
	conf, errX := viperLib.NewFileViper (confFileName, "yaml")
	if errX != nil {
		return nil, err.New ("Unable to load conf file.", nil, nil, errX)
	}

	output = &Conf {}


// Section X

	// Processing conf data 'http_server.addr'.  ..1.. {
	if ! conf.IsSet ("http_server.addr") {
		return nil, err.New ("Conf data 'http_server.addr': Data not set.", nil,
			nil)
	}
	(*output) ["http_server.addr"] = conf.GetString ("http_server.addr")
	if ! govalidator.IsIP ((*output) ["http_server.addr"]) {
		return nil, err.New ("Conf data 'http_server.addr': Invalid IP v4 or " +
			"v6 address.", nil, nil)
	}
	// .. }

	// Processing conf data 'http_server.port'.  ..1.. {
	if ! conf.IsSet ("http_server.port") {
		return nil, err.New ("Conf data 'http_server.port': Data not set.", nil,
			nil)
	}
	(*output) ["http_server.port"] = conf.GetString ("http_server.port")
	portA, okA := strconv.Atoi ((*output) ["http_server.port"])
	if okA != nil || portA > 65535 {
		return nil, err.New ("Conf data 'http_server.port': Invalid port number.",
			nil, nil)
	}
	if portA == 0 {
		return nil, err.New ("Conf data 'http_server.port': Port 0 can not be " +
			"used.", nil, nil)
	}
	// .. }

	// Processing conf data 'http_server.tls_key'.  ..1.. {	
	if ! conf.IsSet ("http_server.tls_key") ||
		conf.GetString ("http_server.tls_key") == "" {

		return nil, err.New ("Conf data 'http_server.tls_key': Data not set.",
			nil, nil)
	}
	(*output) ["http_server.tls_key"] = conf.GetString ("http_server.tls_key")
	okB, errB := afero.Exists (afero.NewOsFs (), (*output) ["http_server.tls_key"])
	if errB != nil {
		return nil, err.New ("Conf data 'http_server.tls_key': Unable to " +
			"confirm existence of file.", nil, nil)
	} else if okB == false {
		return nil, err.New ("Conf data 'http_server.tls_key': File not found.",
			nil, nil)
	}
	// .. }

	// Processing conf data 'http_server.tls_crt'.  ..1.. {	
	if ! conf.IsSet ("http_server.tls_crt") ||
		conf.GetString ("http_server.tls_crt") == "" {

		return nil, err.New ("Conf data 'http_server.tls_crt': Data not set.",
			nil, nil)
	}
	(*output) ["http_server.tls_crt"] = conf.GetString ("http_server.tls_crt")
	okC, errC := afero.Exists (afero.NewOsFs (), (*output) ["http_server.tls_crt"])
	if errC != nil {
		return nil, err.New ("Conf data 'http_server.tls_crt': Unable to " +
			"confirm existence of file.", nil, nil)
	} else if okC == false {
		return nil, err.New ("Conf data 'http_server.tls_crt': File not found.",
			nil, nil)
	}
	// .. }


// Section Y

	// Processing conf data 'sds.addr'.  ..1.. {
	if ! conf.IsSet ("sds.addr") {
		return nil, err.New ("Conf data 'sds.addr': Data not set.", nil, nil)
	}
	(*output) ["sds.addr"] = conf.GetString ("sds.addr")
	if ! govalidator.IsHost ((*output) ["sds.addr"]) {
		return nil, err.New ("Conf data 'sds.addr': Invalid hostname.", nil, nil)
	}
	// .. }

	// Processing conf data 'sds.port'.  ..1.. {
	if ! conf.IsSet ("sds.port") {
		return nil, err.New ("Conf data 'sds.port': Data not set.", nil, nil)
	}
	(*output) ["sds.port"] = conf.GetString ("sds.port")
	portD, okD := strconv.Atoi ((*output) ["sds.port"])
	if okD != nil || portD > 65535 {
		return nil, err.New ("Conf data 'sds.port': Invalid port number.", nil,
			nil)
	}
	if portD == 0 {
		return nil, err.New ("Conf data 'sds.port': Port 0 can not be used.",
			nil, nil)
	}
	// .. }

	// Processing conf data 'sds.pub_key_file'.  ..1.. {	
	if ! conf.IsSet ("sds.pub_key_file") || conf.GetString ("sds.pub_key_file") == "" {
		return nil, err.New ("Conf data 'sds.pub_key_file': Data not set.",
			nil, nil)
	}
	(*output) ["sds.pub_key_file"] = conf.GetString ("sds.pub_key_file")
	okE, errE := afero.Exists (afero.NewOsFs (), (*output) ["sds.pub_key_file"])
	if errE != nil {
		return nil, err.New ("Conf data 'sds.pub_key_file': Unable to confirm " +
		"existence of file.", nil, nil)
	} else if okE == false {
		return nil, err.New ("Conf data 'sds.pub_key_file': File not found.", nil,
			nil)
	}
	// .. }

	// Processing conf data 'sds.username'.  ..1.. {
	if ! conf.IsSet ("sds.username") || conf.GetString ("sds.username") == "" {
		return nil, err.New ("Conf data 'sds.username': Data not set.", nil, nil)
	}
	(*output) ["sds.username"] = conf.GetString ("sds.username")
	// .. }

	// Processing conf data 'sds.pass'.  ..1.. {
	if ! conf.IsSet ("sds.pass") || conf.GetString ("sds.pass") == "" {
		return nil, err.New ("Conf data 'sds.pass': Data not set.", nil, nil)
	}
	(*output) ["sds.pass"] = conf.GetString ("sds.pass")
	// .. }

	// Processing conf data 'sds.conn_timeout'.  ..1.. {
	if ! conf.IsSet ("sds.conn_timeout") {
		return nil, err.New ("Conf data 'sds.conn_timeout': Data not set.", nil,
			nil)
	}
	(*output) ["sds.conn_timeout"] = conf.GetString ("sds.conn_timeout")
	timeoutF, errF := strconv.Atoi ((*output) ["sds.conn_timeout"])
	if errF != nil {
		return nil, err.New ("Conf data 'sds.conn_timeout': Value seems invalid.",
			nil, nil, errF)
	}
	if timeoutF == 0 {
		return nil, err.New ("Conf data 'sds.conn_timeout': Timeout can not be " +
			"zero.", nil, nil)
	}
	// .. }

	// Processing conf data 'sds.write_timeout'.  ..1.. {
	if ! conf.IsSet ("sds.write_timeout") {
		return nil, err.New ("Conf data 'sds.write_timeout': Data not set.", nil,
			nil)
	}
	(*output) ["sds.write_timeout"] = conf.GetString ("sds.write_timeout")
	timeoutG, errG := strconv.Atoi ((*output) ["sds.write_timeout"])
	if errG != nil {
		return nil, err.New ("Conf data 'sds.write_timeout': Value seems invalid.",
			nil, nil, errG)
	}
	if timeoutG == 0 {
		return nil, err.New ("Conf data 'sds.write_timeout': Timeout can not be " +
			"zero.", nil, nil)
	}
	// .. }

	// Processing conf data 'sds.read_timeout'.  ..1.. {
	if ! conf.IsSet ("sds.read_timeout") {
		return nil, err.New ("Conf data 'sds.read_timeout': Data not set.", nil,
			nil)
	}
	(*output) ["sds.read_timeout"] = conf.GetString ("sds.read_timeout")
	timeoutH, errH := strconv.Atoi ((*output) ["sds.read_timeout"])
	if errH != nil {
		return nil, err.New ("Conf data 'sds.read_timeout': Value seems invalid.",
			nil, nil, errH)
	}
	if timeoutH == 0 {
		return nil, err.New ("Conf data 'sds.read_timeout': Timeout can not be " +
			"zero.", nil, nil)
	}
	// .. }

	// Processing conf data 'sds.service_id'.  ..1.. {
	if ! conf.IsSet ("sds.service_id") || conf.GetString ("sds.service_id") == "" {
		return nil, err.New ("Conf data 'sds.service_id': Data not set.", nil, nil)
	}
	(*output) ["sds.service_id"] = conf.GetString ("sds.service_id")
	if ! serviceIDPattern.Match ([]byte ((*output) ["sds.service_id"])) {
		return nil, err.New (`Conf data 'sds.service_id': Invalid value, " +
			"according to the "SDS v0.3.0".`, nil, nil)
	}
	// .. }

	// Processing conf data 'sds.update_pass'.  ..1.. {
	if ! conf.IsSet ("sds.update_pass") || conf.GetString ("sds.update_pass") == "" {
		return nil, err.New ("Conf data 'sds.update_pass': Data not set.", nil,
			nil)
	}
	(*output) ["sds.update_pass"] = conf.GetString ("sds.update_pass")
	if ! updatePassPattern.Match ([]byte ((*output) ["sds.update_pass"])) {
		return nil, err.New (`Conf data 'sds.update_pass': Invalid value, " +
			"according to the "SDS v0.3.0".`, nil, nil)
	}
	// .. }


// Section Z

	// Processing conf data 'dbms.addr'.  ..1.. {
	if ! conf.IsSet ("dbms.addr") {
		return nil, err.New ("Conf data 'dbms.addr': Data not set.", nil, nil)
	}
	(*output) ["dbms.addr"] = conf.GetString ("dbms.addr")
	if ! govalidator.IsHost ((*output) ["dbms.addr"]) {
		return nil, err.New ("Conf data 'dbms.addr': Invalid hostname.", nil, nil)
	}
	// .. }

	// Processing conf data 'dbms.port'.  ..1.. {
	if ! conf.IsSet ("dbms.port") {
		return nil, err.New ("Conf data 'dbms.port': Data not set.", nil, nil)
	}
	(*output) ["dbms.port"] = conf.GetString ("dbms.port")
	portI, okI := strconv.Atoi ((*output) ["dbms.port"])
	if okI != nil || portI > 65535 {
		return nil, err.New ("Conf data 'dbms.port': Invalid port number.", nil,
			nil)
	}
	if portI == 0 {
		return nil, err.New ("Conf data 'dbms.port': Port 0 can not be used.", nil,
			nil)
	}
	// .. }

	// Processing conf data 'dbms.pub_key_file'.  ..1.. {	
	if ! conf.IsSet ("dbms.pub_key_file") || conf.GetString ("dbms.pub_key_file") == "" {
		return nil, err.New ("Conf data 'dbms.pub_key_file': Data not set.", nil,
			nil)
	}
	(*output) ["dbms.pub_key_file"] = conf.GetString ("dbms.pub_key_file")
	okJ, errJ := afero.Exists (afero.NewOsFs (), (*output) ["dbms.pub_key_file"])
	if errJ != nil {
		return nil, err.New ("Conf data 'dbms.pub_key_file': Unable to confirm " +
			"existence of file.", nil, nil)
	} else if okJ == false {
		return nil, err.New ("Conf data 'dbms.pub_key_file': File not found.", nil,
			nil)
	}
	// .. }

	// Processing conf data 'dbms.username'.  ..1.. {
	if ! conf.IsSet ("dbms.username") || conf.GetString ("dbms.username") == "" {
		return nil, err.New ("Conf data 'dbms.username': Data not set.", nil, nil)
	}
	(*output) ["dbms.username"] = conf.GetString ("dbms.username")
	// .. }

	// Processing conf data 'dbms.pass'.  ..1.. {
	if ! conf.IsSet ("dbms.pass") || conf.GetString ("dbms.pass") == "" {
		return nil, err.New ("Conf data 'dbms.pass': Data not set.", nil, nil)
	}
	(*output) ["dbms.pass"] = conf.GetString ("dbms.pass")
	// .. }

	// Processing conf data 'dbms.conn_timeout'.  ..1.. {
	if ! conf.IsSet ("dbms.conn_timeout") {
		return nil, err.New ("Conf data 'dbms.conn_timeout': Data not set.", nil,
			nil)
	}
	(*output) ["dbms.conn_timeout"] = conf.GetString ("dbms.conn_timeout")
	timeoutK, errK := strconv.Atoi ((*output) ["dbms.conn_timeout"])
	if errK != nil {
		return nil, err.New ("Conf data 'dbms.conn_timeout': Value seems invalid.",
			nil, nil, errK)
	}
	if timeoutK == 0 {
		return nil, err.New ("Conf data 'dbms.conn_timeout': Timeout can not be " +
			"zero.", nil, nil)
	}
	// .. }

	// Processing conf data 'dbms.write_timeout'.  ..1.. {
	if ! conf.IsSet ("dbms.write_timeout") {
		return nil, err.New ("Conf data 'dbms.write_timeout': Data not set.", nil,
			nil)
	}
	(*output) ["dbms.write_timeout"] = conf.GetString ("dbms.write_timeout")
	timeoutL, errL := strconv.Atoi ((*output) ["dbms.write_timeout"])
	if errL != nil {
		return nil, err.New ("Conf data 'dbms.write_timeout': Value seems " +
			"invalid.", nil, nil, errL)
	}
	if timeoutL == 0 {
		return nil, err.New ("Conf data 'dbms.write_timeout': Timeout can not " +
			"be zero.", nil, nil)
	}
	// .. }

	// Processing conf data 'dbms.read_timeout'.  ..1.. {
	if ! conf.IsSet ("dbms.read_timeout") {
		return nil, err.New ("Conf data 'dbms.read_timeout': Data not set.", nil,
			nil)
	}
	(*output) ["dbms.read_timeout"] = conf.GetString ("dbms.read_timeout")
	timeoutM, errM := strconv.Atoi ((*output) ["dbms.read_timeout"])
	if errM != nil {
		return nil, err.New ("Conf data 'dbms.read_timeout': Value seems invalid.",
			nil, nil, errM)
	}
	if timeoutM == 0 {
		return nil, err.New ("Conf data 'dbms.read_timeout': Timeout can not be " +
			"zero.", nil, nil)
	}
	// .. }

	return output, nil
}

type Conf map[string]string

func (c *Conf) Get (name string) (string) {
	output, _ := (*c) [name]
	return output
}

var (
	serviceIDPattern *regexp.Regexp
	updatePassPattern *regexp.Regexp
	confFileName string = "./conf.yml"
)

func init () {
	if initReport != nil {
		return
	}

	var errX error
	serviceIDPattern, errX = regexp.Compile ("^[a-z0-9]{1,2}$")
	if errX != nil {
		initReport = err.New ("Service ID regexp compilation failed.", nil, nil,
			serviceIDPattern)
	}

	var errY error
	updatePassPattern, errY = regexp.Compile ("^[a-z0-9]{32,32}$")
		if errX != nil {
		initReport = err.New (`"Record-update password" regexp compilation " +
			"failed.`, nil, nil, updatePassPattern)
	}
}
