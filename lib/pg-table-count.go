package mppgtablecount

import (
	"flag"
	"os"

	"fmt"
	"github.com/jmoiron/sqlx"
	// PostgreSQL Driver
	_ "github.com/lib/pq"
	mp "github.com/mackerelio/go-mackerel-plugin"
	"strings"
)

// PgTableCountPlugin mackerel plugin
type PgTableCountPlugin struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	Table    string
	Column   string
	Option   string
	SSLmode  string
	Timeout  int
	Prefix   string
}

// PgCount postgres table count struct
type PgCount struct {
	Count uint64 `db:"count"`
}

// MetricKeyPrefix interface for PluginWithPrefix
func (p *PgTableCountPlugin) MetricKeyPrefix() string {
	if p.Prefix == "" {
		p.Prefix = "postgres"
	}
	return p.Prefix
}

// GraphDefinition interface for mackerelplugin
func (p *PgTableCountPlugin) GraphDefinition() map[string]mp.Graphs {
	labelPrefix := strings.Title(p.Prefix)
	return map[string]mp.Graphs{
		"table.#": {
			Label: labelPrefix + " Table Counts",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "counts", Label: "Counts"},
			},
		},
	}
}

// FetchMetrics interface for mackerelplugin
func (p *PgTableCountPlugin) FetchMetrics() (map[string]float64, error) {
	db, err := sqlx.Connect("postgres",
		fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s connect_timeout=%d", p.User, p.Password, p.Host, p.Port, p.Database, p.SSLmode, p.Timeout))
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var query string
	if p.Option == "" {
		query = fmt.Sprintf("SELECT count(%s) FROM %s", p.Column, p.Table)
	} else {
		query = fmt.Sprintf("SELECT count(%s) FROM %s %s", p.Column, p.Table, p.Option)
	}

	db = db.Unsafe()
	rows, err := db.Queryx(query)
	if err != nil {
		return nil, err
	}

	metrics := make(map[string]float64)
	for rows.Next() {
		var result PgCount
		if err := rows.StructScan(&result); err != nil {
			return nil, err
		}
		metrics["table."+p.Table+".counts"] = float64(result.Count)
	}
	return metrics, nil
}

// Do the plugin
func Do() {
	optHost := flag.String("host", "localhost", "Hostname")
	optPort := flag.String("port", "5432", "Port")
	optUser := flag.String("user", "postgres", "Username")
	optPassword := flag.String("password", os.Getenv("PGPASSEORD"), "Password")
	optDatabase := flag.String("database", "", "Database")
	optTable := flag.String("tabel", "", "Table")
	optColumn := flag.String("column", "*", "Count column")
	optOption := flag.String("option", "", "Query option")
	optSSLmode := flag.String("sslmode", "disable", "Whether or not to use SSL")
	optConnectTimeout := flag.Int("connect-timeout", 5, "Maximum wait for connection, in seconds.")
	optPrefix := flag.String("metric-key-prefix", "postgres", "Metric key prefix")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	plugin := mp.NewMackerelPlugin(&PgTableCountPlugin{
		Host:     *optHost,
		Port:     *optPort,
		User:     *optUser,
		Password: *optPassword,
		Database: *optDatabase,
		Table:    *optTable,
		Column:   *optColumn,
		Option:   *optOption,
		SSLmode:  *optSSLmode,
		Timeout:  *optConnectTimeout,
		Prefix:   *optPrefix,
	})
	plugin.Tempfile = *optTempfile
	plugin.Run()
}
