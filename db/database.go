package db

import (
	"regexp"
	"strings"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"os"
    "go-cloud-ssl/models"

	mysql "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var DB *sql.DB

type Config struct {
	Database struct {
		Driver   string `yaml:"driver"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
	SSL struct {
		Enabled    bool   `yaml:"enabled"`
		CACert     string `yaml:"ca_cert"`
		ClientCert string `yaml:"client_cert"`
		ClientKey  string `yaml:"client_key"`
	} `yaml:"ssl"`
}

func InitDB(cfg Config) {
	var dsn string
	var err error

	if cfg.Database.Driver == "mysql" {
		if cfg.SSL.Enabled {
			rootCertPool := x509.NewCertPool()
			pem, err := os.ReadFile(cfg.SSL.CACert)
			if err != nil {
				panic(err)
			}
			rootCertPool.AppendCertsFromPEM(pem)

			certs, err := tls.LoadX509KeyPair(cfg.SSL.ClientCert, cfg.SSL.ClientKey)
			if err != nil {
				panic(err)
			}

			mysql.RegisterTLSConfig("custom", &tls.Config{
				RootCAs:      rootCertPool,
				Certificates: []tls.Certificate{certs},
			})
		}

		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?tls=custom",
			cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
	} else {
		sslMode := "disable"
		if cfg.SSL.Enabled {
			sslMode = "verify-ca"
		}
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s sslcert=%s sslkey=%s sslrootcert=%s",
			cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name, sslMode,
			cfg.SSL.ClientCert, cfg.SSL.ClientKey, cfg.SSL.CACert)
	}

	DB, err = sql.Open(cfg.Database.Driver, dsn)
	if err != nil {
		panic(err)
	}

	models.Migrate(DB)
}

// PrepareQuery replaces ? placeholders to $n if driver is postgres
func PrepareQuery(query, driver string) string {
	if driver != "postgres" {
		return query // MySQL and others use "?"
	}

	// Replace each ? with $1, $2, ...
	var i int
	return regexp.MustCompile(`\?`).ReplaceAllStringFunc(query, func(_ string) string {
		i++
		return "$" + strings.TrimSpace((string(rune(i + '0'))))
	})
}



func Execute(config Config, query string, args ...interface{}) (sql.Result, error) {
	if DB == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}


	if config.Database.Driver == "mysql" {
		// For MySQL, use placeholders like ?
		query = query
	} else {
		// For PostgreSQL, use placeholders like $1, $2, etc.
		query = PrepareQuery(query, config.Database.Driver)
	}	

	println(config.Database.Driver, query, args)


	return DB.Exec(query, args...)
}

func ExecuteRow(config Config, query string, args ...interface{}) (*sql.Row) {

	if config.Database.Driver == "mysql" {
		// For MySQL, use placeholders like ?
		query = query
	} else {
		// For PostgreSQL, use placeholders like $1, $2, etc.
		query = PrepareQuery(query, config.Database.Driver)
	}	

	return DB.QueryRow(query, args...)
}