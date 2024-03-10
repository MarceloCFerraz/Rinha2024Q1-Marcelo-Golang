package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	// "runtime/pprof"
	"time"

	models "github.com/MarceloCFerraz/Rinha2024Q1-Marcelo-Golang/Models"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	connectionPool, poolErr := pgxpool.New(context.Background(), getConnString())
	if poolErr != nil {
		fmt.Println("Error creating conn pool:\n", poolErr)
		os.Exit(1)
	}

	app := fiber.New(fiber.Config{
		AppName:       "Rinha 2024 - Marcelo - Go",
		CaseSensitive: true,
		// Prefork:       true,
		// RequestMethods: []string{fiber.MethodPost, fiber.MethodGet},
	})

	// fcpu, _ := os.Create("cpuprofile.out")
	// defer fcpu.Close()

	clients := app.Group("/clientes")

	clients.Get("/:id/extrato", func(c *fiber.Ctx) error {

		// pprof.StartCPUProfile(fcpu)
		// defer pprof.StopCPUProfile()
		id_cliente, noIdErr := c.ParamsInt("id")

		if noIdErr != nil {
			fmt.Println("No id provided", noIdErr)
			return c.SendStatus(422) // bad request, no id provided
		}

		client := models.Client{Id: id_cliente}

		if client.IsInvalid() {
			// fmt.Println("client is INVALID (GET): ", id_cliente)
			return c.SendStatus(404)
			// invalid client id
		}
		//else {
		// fmt.Println("client is valid (GET): ", id_cliente)
		//}

		conn, connErr := connectionPool.Acquire(context.Background())

		if connErr != nil {
			fmt.Println("An error ocurred when getting a conn from pool:\n", connErr)
			return c.SendStatus(500)
		}

		defer conn.Release() // re-add conn back to pool once finished

		var customer_balance, customer_limit int
		var report_date time.Time
		var last_transactions []models.Transaction

		queryErr := conn.QueryRow(
			context.Background(),
			"SELECT customer_balance, customer_limit, report_date, last_transactions FROM get_client_data($1)",
			id_cliente). /*  */
			Scan(&customer_balance, &customer_limit, &report_date, &last_transactions)

		if queryErr != nil {
			fmt.Println("An error ocurred when executing query:\n", queryErr)
			return c.SendStatus(500)
		}

		return c.Status(200).JSON(fiber.Map{
			"saldo": fiber.Map{
				"total":        customer_balance,
				"limite":       customer_limit,
				"data_extrato": report_date.UTC(),
			},
			"ultimas_transacoes": last_transactions,
		})
	})

	clients.Post("/:id/transacoes", func(c *fiber.Ctx) error {
		c.Set("content-type", "application/json")

		// pprof.StartCPUProfile(fcpu)
		// defer pprof.StopCPUProfile()

		id_cliente, noIdErr := c.ParamsInt("id")

		if noIdErr != nil {
			fmt.Println("No id provided", noIdErr)
			return c.SendStatus(422) // bad request, no id provided
		}

		client := models.Client{Id: id_cliente}

		if client.IsInvalid() {
			// fmt.Println("client is INVALID: ", id_cliente)
			return c.SendStatus(404)
			// invalid client id
		}
		//else {
		// fmt.Println("client is valid (POST): ", id_cliente)
		//}

		var transaction models.Transaction

		// invBodyErr := c.BodyParser(&transaction) // not using this to not enforce the header '"Content-Type": "application/json"'
		invBodyErr := json.Unmarshal(c.Request().Body(), &transaction)

		if invBodyErr != nil {
			// fmt.Println(reflect.TypeOf(invBodyErr), invBodyErr)
			return c.SendStatus(422) // bad request, no id provided
		}

		if transaction.IsInvalid() {
			// fmt.Println("transaction body is INVALID: ", transaction)
			return c.SendStatus(422)
		}
		// else {
		// fmt.Println("transaction body is valid: ", transaction)
		//}

		conn, connErr := connectionPool.Acquire(context.Background())

		if connErr != nil {
			fmt.Println("An error ocurred when getting a conn from pool:\n", connErr)
			return c.SendStatus(500)
		}

		defer conn.Release() // re-add conn back to pool once finished

		var success bool
		var client_limit, new_balance int

		command := ("SELECT success, client_limit, new_balance FROM " + transaction.GetDbOperation() + "($1, $2, $3)")

		queryErr := conn.QueryRow(
			context.Background(),
			command,
			id_cliente,
			int(transaction.Value),
			transaction.Description).Scan(&success, &client_limit, &new_balance)

		if queryErr != nil {
			fmt.Println("An error ocurred when executing query:\n", queryErr)
			return c.SendStatus(500)
		}

		if success {
			return c.Status(200).JSON(fiber.Map{
				"saldo":  new_balance,
				"limite": client_limit,
			})
		} else {
			// fmt.Println("Transaction UNSUCCESSFUL: ", transaction)
			return c.SendStatus(422)
		}
	})

	app.Listen(getPortFromEnvOrDefault())
}

func getPortFromEnvOrDefault() string {
	port := getEnvOrDefault("API_PORT", "8080")

	return fmt.Sprintf(":%s", port)
}

func getConnString() string {
	// NOTE: The recommended way to establish a connection with pgx is to use pgx.ParseConfig but i just want to get this running
	// Valid formats:
	// URL Format: postgres://user:password@host:port/database?param1=value1&...
	// DSN Format: user=username password=password host=host port=port dbname=database param1=value1 ...
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbPort := getEnvOrDefault("DB_PORT", "5432")
	dbName := getEnvOrDefault("DB_NAME", "rinha")
	dbUser := getEnvOrDefault("DB_USER", "admin")
	dbPass := getEnvOrDefault("DB_PASS", "mystrongpassword")
	dbCons := getEnvOrDefault("MAX_DB_CONNECTIONS", "100")

	connString := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s pool_max_conns=%s",
		dbUser, dbPass, dbHost, dbPort, dbName, dbCons)

	return connString
}

func getEnvOrDefault(key string, defaultValue string) string {
	value, found := os.LookupEnv(key)
	if found {
		return value
	}
	return defaultValue
}
