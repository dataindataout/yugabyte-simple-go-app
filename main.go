//
//  Copyright 2022 Yugabyte
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.
//

package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
)

const (
	host        = "127.0.0.1"
	port        = 5433
	dbName      = "yugabyte"
	dbUser      = "yugabyte"
	dbPassword  = "yugabyte"
	sslMode     = "disable"
	sslRootCert = "/tmp/asdf"
)

func checkIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func createDatabase(conn *pgx.Conn) {
	stmt := `DROP TABLE IF EXISTS DemoAccount`
	_, err := conn.Exec(context.Background(), stmt)
	checkIfError(err)

	stmt = `CREATE TABLE DemoAccount (
                        id int PRIMARY KEY,
                        name varchar,
                        age int,
                        country varchar,
                        balance int)`

	_, err = conn.Exec(context.Background(), stmt)
	checkIfError(err)

	stmt = `INSERT INTO DemoAccount VALUES
                (1, 'Jessica', 28, 'USA', 10000),
                (2, 'John', 28, 'Canada', 9000)`

	_, err = conn.Exec(context.Background(), stmt)
	checkIfError(err)

	fmt.Println(">>>> Successfully created table DemoAccount.")
}

func selectAccounts(conn *pgx.Conn) {
	fmt.Println(">>>> Selecting accounts:")

	rows, err := conn.Query(context.Background(), "SELECT name, age, country, balance FROM DemoAccount")
	checkIfError(err)

	defer rows.Close()

	var name, country string
	var age, balance int

	for rows.Next() {
		err = rows.Scan(&name, &age, &country, &balance)
		checkIfError(err)

		fmt.Printf("name = %s, age = %v, country = %s, balance = %v\n",
			name, age, country, balance)
	}
}

func transferMoneyBetweenAccount(conn *pgx.Conn, amount int) {

	tx, err := conn.Begin(context.Background())
	checkIfError(err)

	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), `UPDATE DemoAccount SET balance = balance - $1 WHERE name = 'Jessica'`, amount)
	if checkIfTxAborted(err) {
		return
	}
	_, err = tx.Exec(context.Background(), `UPDATE DemoAccount SET balance = balance + $1 WHERE name = 'John'`, amount)
	if checkIfTxAborted(err) {
		return
	}

	err = tx.Commit(context.Background())
	if checkIfTxAborted(err) {
		return
	}

	fmt.Printf(">>>> Transferred %v between accounts.\n", amount)
}

func checkIfTxAborted(err error) bool {

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			fmt.Println(pgErr.Message) // => syntax error at end of input
			fmt.Println(pgErr.Code)    // => 42601
			if pgErr.Code == `40001` {
				fmt.Println(
					`The operation is aborted due to a concurrent transaction that is modifying the same set of rows.
				Consider adding retry logic or using pessimistic locking.`)
			}
		}

		log.Fatal(err)
	}

	return false
}

func main() {
	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		dbUser, dbPassword, host, port, dbName)

	if sslMode != "" {
		url += fmt.Sprintf("?sslmode=%s", sslMode)

		if sslRootCert != "" {
			url += fmt.Sprintf("&sslrootcert=%s", sslRootCert)
		}
	}

	conn, err := pgx.Connect(context.Background(), url)
	checkIfError(err)
	defer conn.Close(context.Background())

	fmt.Println(">>>> Successfully connected to YugabyteDB!")

	createDatabase(conn)
	selectAccounts(conn)
	transferMoneyBetweenAccount(conn, 800)
	selectAccounts(conn)
}
