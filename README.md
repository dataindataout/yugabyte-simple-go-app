# Simple Go Application for YugabyteDB

The application connects to your YugabyteDB instance via
[Go PostgreSQL driver pgx](https://docs.yugabyte.com/preview/reference/drivers/go/pgx-reference/) and performs basic SQL
operations. The instructions below are provided for local deployments [Yugabyte Docs for single node](https://docs.yugabyte.com/preview/quick-start/) or [Yugabyte Docs for local multi-node](https://docs.yugabyte.com/preview/explore/multi-region-deployments/synchronous-replication-ysql/#create-a-multi-zone-cluster-in-us-west). If you use a different type of deployment, then update the `sample-app.go` file with proper connection parameters.

This example was adapted from the [pq example](https://github.com/YugabyteDB-Samples/yugabyte-simple-go-app).

## Prerequisites

* Go version 1.19.4 or later is preferred. Earlier versions should work as well.
* Command line tool or your favorite IDE, such as VSCode.

## Clone Application Repository

Clone the application to your machine:

```bash
git clone https://github.com/dataindataout/yugabyte-simple-go-app.git && cd yugabyte-simple-go-app
```

## Provide Cluster Connection Parameters

Open the `main.go` file and specify the following configuration parameters:

* `host` - the hostname of your instance.
* `port` - the port number of your instance (the default is `5433`).
* `dbUser` - the username for your instance.
* `dbPassword` - the database password.
* `sslMode` - the SSL mode. Set to `verify-full` for YugabyteDB Managed deployments.
* `sslRootCert` - a full path to your CA root cert if used (for example, `/Users/ybme/certificates/root.crt`)

## Build and Run Application

1. Import the required packages:

    ```bash
    go mod tidy
    ```

2. Run the application:

    ```bash
    go run main.go
    ```

Upon successful execution, you will see output similar to the following:

```bash
>>>> Successfully connected to YugabyteDB!
>>>> Successfully created table DemoAccount.
>>>> Selecting accounts:
name = Jessica, age = 28, country = USA, balance = 10000
name = John, age = 28, country = Canada, balance = 9000
>>>> Transferred 800 between accounts.
>>>> Selecting accounts:
name = Jessica, age = 28, country = USA, balance = 9200
name = John, age = 28, country = Canada, balance = 9800
```

## Explore App Logic

Congrats! You've successfully executed a simple Go app that works with YugabyteDB.

Now, explore the source code of `main.go`:

1. `main` function - establishes a connection with your cloud instance via the Go PostgreSQL driver pgx.
2. `createDatabase` function - creates a table and populates it with sample data.
3. `selectAccounts` function - queries the data with SQL `SELECT` statements.
4. `transferMoneyBetweenAccounts` function - updates records consistently with distributed transactions.

## Questions or Issues?

Having issues running this application or want to learn more from Yugabyte experts?

Join [our Slack channel](https://communityinviter.com/apps/yugabyte-db/register),
or raise a question on StackOverflow and tag the question with `yugabytedb`!
