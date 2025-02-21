package cmd

type Configuration struct {
	DSName string `short:"n" long:"ds" env:"DATASTORE" description:"DataStore name (format: mongo/null)" required:"false" default:"mongo"`
	DSDB   string `short:"d" long:"ds-db" env:"DATASTORE_DB" description:"DataStore database name (format: go-rest-template)" required:"false" default:"go-rest-template"`
	DSURL  string `short:"u" long:"ds-url" env:"DATASTORE_URL" description:"DataStore URL (format: mongodb://localhost:27017)" required:"false" default:"mongodb://localhost:27017"`

	ListenAddr string `short:"l" long:"listen" env:"LISTEN" description:"Listen Address (format: :8080|127.0.0.1:8080)" required:"false" default:":8080"`
	BasePath   string `long:"base-path" env:"BASE_PATH" description:"base path of the host" required:"false" default:"sso"`
	CertFile   string `short:"c" long:"cert" env:"CERT_FILE" description:"Location of the SSL/TLS cert file" required:"false" default:""`
	KeyFile    string `short:"k" long:"key" env:"KEY_FILE" description:"Location of the SSL/TLS key file" required:"false" default:""`

	Dbg       bool `long:"dbg" env:"DEBUG" description:"debug mode"`
	IsTesting bool `long:"testing" env:"APP_TESTING" description:"testing mode"`
}
