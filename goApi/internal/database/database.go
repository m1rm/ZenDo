package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"zendo/models"

	_ "github.com/go-sql-driver/mysql"
)

// Service represents the service that interacts with the database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error

	GetTodos() ([]models.Todo, error)

	GetTodo(id int64) (models.Todo, error)

	DeleteTodo(id int64) (int64, error)

	InsertTodo(todo models.Todo) (int64, error)
}

type service struct {
	db *sql.DB
}

var (
	dbname     = "zendo"
	password   = "football"
	username   = "goApi"
	port       = "3306"
	host       = "db:3306"
	dbInstance *service
)

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	// Opening a driver typically will not attempt to connect to the database.
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbname))
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Fatal(err)
	}
	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(50)

	dbInstance = &service{
		db: db,
	}
	return dbInstance
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf(fmt.Sprintf("db down: %v", err)) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}
	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

// Close the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", dbname)
	return s.db.Close()
}

func (s *service) GetTodos() ([]models.Todo, error) {
	var todos []models.Todo

	rows, err := s.db.Query("SELECT * FROM todos")
	if err != nil {
		return nil, fmt.Errorf("queryTodos: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var todo models.Todo
		if err := rows.Scan(&todo.Id, &todo.Description, &todo.Status); err != nil {
			return nil, fmt.Errorf("queryTodos: %v", err)
		}
		todos = append(todos, todo)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("queryTodos: %v", err)
	}
	return todos, nil
}

func (s *service) GetTodo(id int64) (models.Todo, error) {
	var todo models.Todo

	row := s.db.QueryRow("SELECT * FROM todos WHERE id = ?", id)
	if err := row.Scan(&todo.Id, &todo.Description, &todo.Status); err != nil {
		if err == sql.ErrNoRows {
			return todo, fmt.Errorf("queryTodo %d: no such todo", id)
		}
		return todo, fmt.Errorf("queryTodo %d: %v", id, err)
	}
	return todo, nil
}

func (s *service) InsertTodo(todo models.Todo) (int64, error) {
	query := "INSERT INTO `todos` (`description`, `status`) VALUES (?, ?);"
	insertResult, err := s.db.ExecContext(context.Background(), query, &todo.Description, &todo.Status)
	if err != nil {
		return 0, err
	}
	id, err := insertResult.LastInsertId()
	if err != nil {
		return id, err
	}
	return id, nil
}

func (s *service) DeleteTodo(id int64) (int64, error) {
	_, err := s.GetTodo(id)
	if err != nil {
		return id, err
	}

	query := "DELETE FROM `todos` WHERE `id` = ?;"

	_, execErr := s.db.ExecContext(context.Background(), query, id)
	if execErr != nil {
		return id, err
	}

	return id, nil
}
