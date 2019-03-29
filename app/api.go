
package main

import (
    "flag"
    "fmt"
    "io"
    "net/http"
    "os"
    "runtime"
    "time"
    "database/sql"
    "github.com/gorilla/handlers"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    _ "github.com/go-sql-driver/mysql"
)

var bind = ""
var port int = 8000

const driverName = "mysql"
var pool *sql.DB

var dbPort int = 3306
var dbHost, dbName, dbUser, dbPassword string

var version = false
var versionNumber string
var gitCommit string

var readTimeout = 10
var writeTimeout = 300

const envDbHost = "APP_DB_HOST"
const envDbName = "APP_DB_NAME"
const envDbUser = "APP_DB_USER"
const envDbPassword = "APP_DB_PASSWORD"

const countSql = "SELECT COUNT(*) AS count FROM employees"
const countHeader = "X-Total-Count"

const welcome = "Welcome to api server\n"
const messageOk = "OK"
const messageNotSupported = "NotSupported"

func printVersion() {
    fmt.Printf("api %s\n", versionNumber)
    fmt.Printf("  git commit: %s\n", gitCommit)
    fmt.Printf("  go version: %s\n", runtime.Version())
    os.Exit(0)
}

func validate() {
    // check if required arguments are available
    if version == true {
        printVersion()
    }

    // check environment variables if command-line arguments is nil
    if dbHost == "" && os.Getenv(envDbHost) != "" {
         dbHost = os.Getenv(envDbHost)
    }
    if dbName == "" && os.Getenv(envDbName) != "" {
         dbName = os.Getenv(envDbName)
    }
    if dbUser == "" && os.Getenv(envDbUser) != "" {
         dbUser = os.Getenv(envDbUser)
    }
    if dbPassword == "" && os.Getenv(envDbPassword) != "" {
         dbPassword = os.Getenv(envDbPassword)
    }
}

// handle home page
func indexHandler(w http.ResponseWriter, req *http.Request) {
    io.WriteString(w, welcome)
}

// handle healthcheck
func healthcheckHandler(w http.ResponseWriter, req *http.Request) {
    io.WriteString(w, messageOk)
}

// handle count
func countHandler(w http.ResponseWriter, req *http.Request) {
    var count int

    statusCode := http.StatusBadRequest
    output := messageNotSupported

    if pool != nil {
        stmt, err := pool.Prepare(countSql)
        if err != nil {
            output = fmt.Sprintf("%v", err)
        } else {
            err := stmt.QueryRow().Scan(&count)
            if err != nil {
                output = fmt.Sprintf("%v", err)
            } else {
                statusCode = http.StatusOK
                output = messageOk
                // set the count in response header
                w.Header().Set(countHeader, fmt.Sprintf("%d", count))
            }
        }
    }

    w.WriteHeader(statusCode)
    io.WriteString(w, output)
}

func init() {
    flag.BoolVar(&version, "version", version, "print version")
    flag.IntVar(&port, "port", port, "Port on which api server runs")
    flag.StringVar(&bind, "bind", bind, "To bind to a specific address")
    flag.StringVar(&dbHost, "db-host", dbHost, "db host")
    flag.StringVar(&dbName, "db-name", dbHost, "db name")
    flag.StringVar(&dbUser, "db-user", dbUser, "db user")
    flag.StringVar(&dbPassword, "db-password", dbPassword, "db password")
    flag.IntVar(&readTimeout, "read-timeout", readTimeout, "Read timeout for api server")
    flag.IntVar(&writeTimeout, "write-timeout", writeTimeout, "Write timeout for api server")

    // parse and validate input
    flag.Parse()
    validate()
}

func main() {
    addr := fmt.Sprintf("%s:%d", bind, port)
    fmt.Printf("Starting api server address=%s version=%s gitcommit=%s\n", addr, versionNumber, gitCommit)

    // check if db arguments are available
    if dbHost != "" && dbName != "" && dbUser != "" && dbPassword != ""  {
        var err error
        dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
        pool, err = sql.Open(driverName, dsn)
        if err != nil {
            fmt.Printf("Could not open %s - %v\n", dbHost, err)
            os.Exit(1)
        }
        err = pool.Ping()
        if err != nil {
            fmt.Printf("Could not ping %s - %v\n", dbHost, err)
            os.Exit(1)
        }
        fmt.Printf("Connected to %s:%d/%s\n", dbHost, dbPort, dbName)
    }
    defer pool.Close()

    handler := http.NewServeMux()
    handler.HandleFunc("/", indexHandler)
    handler.HandleFunc("/api/employees", countHandler)
    handler.HandleFunc("/healthcheck", healthcheckHandler)
    handler.Handle("/metrics", promhttp.Handler())

    server := http.Server{
        Addr: addr,
        Handler: handlers.CombinedLoggingHandler(os.Stdout, handler),
        ReadTimeout: time.Duration(readTimeout) * time.Second,
        WriteTimeout: time.Duration(writeTimeout) * time.Second,
    }

    server.ListenAndServe()
}
