package db


import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"time"
	"regexp"

	"cloud.google.com/go/cloudsqlconn"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/go-sql-driver/mysql" // MySQL driver

	"go-cloud-ssl/models"
)



type Config struct {
	Database struct {
		Driver                 string `yaml:"driver"`                    // "postgres" or "mysql"
		InstanceConnectionName string `yaml:"instance_connection_name"` // project:region:instance
		User                   string `yaml:"user"`                      // IAM user
		Name                   string `yaml:"name"`                      // database name
		Private                string `yaml:"private"`                   // private IP address
	} `yaml:"database"`
}


var DB *sql.DB


func InitDB(cfg Config) {
	dbUser := cfg.Database.User
	dbName := cfg.Database.Name
	instanceConnectionName := cfg.Database.InstanceConnectionName
	usePrivate := cfg.Database.Private != ""

	ctx := context.Background()

	// Create Cloud SQL dialer
	dialer, err := cloudsqlconn.NewDialer(ctx,
		cloudsqlconn.WithIAMAuthN(),
		cloudsqlconn.WithLazyRefresh(),
	)
	if err != nil {
		log.Fatalf("cloudsqlconn.NewDialer: %v", err)
	}

	var opts []cloudsqlconn.DialOption
	if usePrivate {
		opts = append(opts, cloudsqlconn.WithPrivateIP())
	}

	var dbPool *sql.DB

	switch cfg.Database.Driver {
	case "postgres":
		// PostgreSQL DSN
		dsn := fmt.Sprintf("user=%s dbname=%s", dbUser, dbName)

		pgxConfig, err := pgx.ParseConfig(dsn)
		if err != nil {
			log.Fatalf("pgx.ParseConfig: %v", err)
		}

		pgxConfig.DialFunc = func(ctx context.Context, network, _ string) (net.Conn, error) {
			log.Printf("Connecting to PostgreSQL instance: %s", instanceConnectionName)
			return dialer.Dial(ctx, instanceConnectionName, opts...)
		}

		dbURI := stdlib.RegisterConnConfig(pgxConfig)
		dbPool, err = sql.Open("pgx", dbURI)
		if err != nil {
			log.Fatalf("sql.Open (Postgres): %v", err)
		}

	case "mysql":
		// Register MySQL dialer
		mysqlDriverName := "cloudsql-mysql"
		// Only register once
		_ = sql.Drivers() // Avoid import optimization removal
		sql.Register(mysqlDriverName, &mysql.MySQLDriver{})

		mysql.RegisterDialContext("cloudsql+mysql", func(ctx context.Context, addr string) (net.Conn, error) {
			log.Printf("Connecting to MySQL instance: %s", instanceConnectionName)
			return dialer.Dial(ctx, instanceConnectionName, opts...)
		})

		dsn := fmt.Sprintf("%s@cloudsql+mysql(%s)/%s?parseTime=true", dbUser, instanceConnectionName, dbName)

		dbPool, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Fatalf("sql.Open (MySQL): %v", err)
		}

	default:
		log.Fatalf("Unsupported driver: %s", cfg.Database.Driver)
	}

	// Optional: test connection
	dbPool.SetMaxOpenConns(5)
	dbPool.SetConnMaxLifetime(time.Minute * 5)

	if err := dbPool.PingContext(ctx); err != nil {
		log.Fatalf("DB Ping failed: %v", err)
	}

	DB = dbPool
	log.Printf("✅ Connected to Cloud SQL (%s) using IAM auth", cfg.Database.Driver)

	// Auto-migrate
	models.Migrate(DB)
}


// PrepareQuery replaces ? placeholders with $n for Postgres
func PrepareQuery(query, driver string) string {
	if driver != "postgres" {
		return query
	}
	i := 0
	return regexp.MustCompile(`\?`).ReplaceAllStringFunc(query, func(_ string) string {
		i++
		return fmt.Sprintf("$%d", i)
	})
}

func Execute(config Config, query string, args ...interface{}) (sql.Result, error) {
	if DB == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}
	query = PrepareQuery(query, config.Database.Driver)
	log.Printf("Executing query: %s with args: %v", query, args)
	return DB.Exec(query, args...)
}

func ExecuteRow(config Config, query string, args ...interface{}) *sql.Row {
	query = PrepareQuery(query, config.Database.Driver)
	return DB.QueryRow(query, args...)
}
