package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	. "fs/service/config"
)

// Backend for Postgres DB
type BackendPostgres struct {
	Kind          string
	ConnectString string
	Sqldb         *sql.DB
}

// NewBackendPostgres creates and opens new Postgres DB connection with GORM layer
func NewBackendPostgres(bc BackendCredentialsPostgres) (BackendPostgres, error) {
	cs := bc.ConnectString()
	Log.Debug("Postgres DB connect string: " + cs)

	sqldb, err := sql.Open("postgres", cs)
	if err != nil {
		return BackendPostgres{}, fmt.Errorf("Error opepning Postgres DB connection: %s", err)
	}
	Log.Debug("Connected Postgres DB")

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqldb.SetMaxIdleConns(Setup.SQLMaxIdleConns)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqldb.SetMaxOpenConns(Setup.SQLMaxOpenConns)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqldb.SetConnMaxLifetime(Setup.SQLMaxLifetime)

	return BackendPostgres{
		Kind:          "postgres",
		ConnectString: cs,
		Sqldb:         sqldb,
	}, nil
}

// Version obtains the backend server version: it is highly database dependent
func (b BackendPostgres) Version() (string, error) {
	var version string
	err := b.Sqldb.QueryRow("SELECT VERSION()").Scan(&version)
	if err != nil {
		return "", fmt.Errorf("Error selection version: %s", err)
	}

	return version, nil
}

// Close backend connection
func (b BackendPostgres) Ping() error {
	err := b.Sqldb.Ping()
	if err != nil {
		return fmt.Errorf("Error pinging postgres: %s", err)
	}

	return nil
}

// Credentials
func (b BackendPostgres) Credentials() string {
	return b.ConnectString
}

// Close backend connection
func (b BackendPostgres) Close() {
	b.Sqldb.Close()
}
