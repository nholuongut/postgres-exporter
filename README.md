# PostgreSQL Server Exporter

![](https://i.imgur.com/waxVImv.png)
### [View all Roadmaps](https://github.com/nholuongut/all-roadmaps) &nbsp;&middot;&nbsp; [Best Practices](https://github.com/nholuongut/all-roadmaps/blob/main/public/best-practices/) &nbsp;&middot;&nbsp; [Questions](https://www.linkedin.com/in/nholuong/)
<br/>

Prometheus exporter for PostgreSQL server metrics.

CI Tested PostgreSQL versions: `9.4`, `9.5`, `9.6`, `10`, `11`, `12`, `13`, `14`

## Quick Start
This package is available for Docker:
```
# Start an example database
docker run --net=host -it --rm -e POSTGRES_PASSWORD=password postgres
# Connect to it
docker run \
  --net=host \
  -e DATA_SOURCE_NAME="postgresql://postgres:password@localhost:5432/postgres?sslmode=disable" \
  quay.io/prometheuscommunity/postgres-exporter
```

## Multi-Target Support (BETA)
**This Feature is in beta and may require changes in future releases. Feedback is welcome.**

This exporter supports the [multi-target pattern](https://prometheus.io/docs/guides/multi-target-exporter/). This allows running a single instance of this exporter for multiple postgres targets. Using the milti-target funcationality of this exporter is **optional** and meant for users where it is impossible to install the exporter as a sidecar. For example SaaS-managed services.

To use the multi-target functionality, send an http request to the endpoint `/probe?target=foo:5432` where target is set to the DSN of the postgres instance to scrape metrics from.

To avoid putting sensitive information like username and password in the URL, preconfigured auth modules are supported via the [auth_modules](#auth_modules) section of the config file. auth_modules for DSNs can be used with the `/probe` endpoint by specifying the `?auth_module=foo` http parameter.

## Configuration File

The configuration file controls the behavior of the exporter. It can be set using the `--config.file` command line flag and defaults to `postgres_exporter.yml`.

### auth_modules
This section defines preset authentication and connection parameters for use in the [multi-target endpoint](#multi-target-support-beta). `auth_modules` is a map of modules with the key being the identifier which can be used in the `/probe` endpoint.
Currently only the `userpass` type is supported.

Example:
```yaml
auth_modules:
  foo1: # Set this to any name you want
    type: userpass
    userpass:
      username: first
      password: firstpass
    options:
      # options become key=value parameters of the DSN
      sslmode: disable
```

## Building and running

    git clone https://github.com/nholuongut/postgres_exporter.git
    cd postgres_exporter
    make build
    ./postgres_exporter <flags>

To build the Docker image:

    make promu
    promu crossbuild -p linux/amd64 -p linux/armv7 -p linux/amd64 -p linux/ppc64le
    make docker

This will build the docker image as `prometheuscommunity/postgres_exporter:${branch}`.

### Flags

* `help`
  Show context-sensitive help (also try --help-long and --help-man).

* `collector.database`
  Enable the pg_database collector. Default is `enabled`

* `collector.bgwriter`
  Enable the pg_stat_bgwriter collector. Default is `enabled`

* `web.listen-address`
  Address to listen on for web interface and telemetry. Default is `:9187`.

* `web.telemetry-path`
  Path under which to expose metrics. Default is `/metrics`.

* `disable-default-metrics`
  Use only metrics supplied from `queries.yaml` via `--extend.query-path`.  Default is `false`.

* `disable-settings-metrics`
  Use the flag if you don't want to scrape `pg_settings`.  Default is `false`.

* `auto-discover-databases`
  Whether to discover the databases on a server dynamically.  Default is `false`.

* `extend.query-path`
  Path to a YAML file containing custom queries to run. Check out [`queries.yaml`](queries.yaml)
  for examples of the format.

* `dumpmaps`
  Do not run - print the internal representation of the metric maps. Useful when debugging a custom
  queries file.

* `constantLabels`
  Labels to set in all metrics. A list of `label=value` pairs, separated by commas.

* `version`
  Show application version.

* `exclude-databases`
  A list of databases to remove when autoDiscoverDatabases is enabled.

* `include-databases`
  A list of databases to only include when autoDiscoverDatabases is enabled.

* `log.level`
  Set logging level: one of `debug`, `info`, `warn`, `error`.

* `log.format`
  Set the log format: one of `logfmt`, `json`.

* `web.config.file`
  Configuration file to use TLS and/or basic authentication. The format of the
  file is described [in the exporter-toolkit repository](https://github.com/prometheus/exporter-toolkit/blob/master/docs/web-configuration.md).

### Environment Variables

The following environment variables configure the exporter:

* `DATA_SOURCE_NAME`
  the default legacy format. Accepts URI form and key=value form arguments. The
  URI may contain the username and password to connect with.

* `DATA_SOURCE_URI`
   an alternative to `DATA_SOURCE_NAME` which exclusively accepts the hostname
   without a username and password component. For example, `my_pg_hostname` or
   `my_pg_hostname?sslmode=disable`.

* `DATA_SOURCE_URI_FILE`
   The same as above but reads the URI from a file.

* `DATA_SOURCE_USER`
  When using `DATA_SOURCE_URI`, this environment variable is used to specify
  the username.

* `DATA_SOURCE_USER_FILE`
  The same, but reads the username from a file.

* `DATA_SOURCE_PASS`
  When using `DATA_SOURCE_URI`, this environment variable is used to specify
  the password to connect with.

* `DATA_SOURCE_PASS_FILE`
  The same as above but reads the password from a file.

* `PG_EXPORTER_WEB_LISTEN_ADDRESS`
  Address to listen on for web interface and telemetry. Default is `:9187`.

* `PG_EXPORTER_WEB_TELEMETRY_PATH`
  Path under which to expose metrics. Default is `/metrics`.

* `PG_EXPORTER_DISABLE_DEFAULT_METRICS`
  Use only metrics supplied from `queries.yaml`. Value can be `true` or `false`. Default is `false`.

* `PG_EXPORTER_DISABLE_SETTINGS_METRICS`
  Use the flag if you don't want to scrape `pg_settings`. Value can be `true` or `false`. Default is `false`.

* `PG_EXPORTER_AUTO_DISCOVER_DATABASES`
  Whether to discover the databases on a server dynamically. Value can be `true` or `false`. Default is `false`.

* `PG_EXPORTER_EXTEND_QUERY_PATH`
  Path to a YAML file containing custom queries to run. Check out [`queries.yaml`](queries.yaml)
  for examples of the format.

* `PG_EXPORTER_CONSTANT_LABELS`
  Labels to set in all metrics. A list of `label=value` pairs, separated by commas.

* `PG_EXPORTER_EXCLUDE_DATABASES`
  A comma-separated list of databases to remove when autoDiscoverDatabases is enabled. Default is empty string.

* `PG_EXPORTER_INCLUDE_DATABASES`
  A comma-separated list of databases to only include when autoDiscoverDatabases is enabled. Default is empty string,
  means allow all.

* `PG_EXPORTER_METRIC_PREFIX`
  A prefix to use for each of the default metrics exported by postgres-exporter. Default is `pg`

Settings set by environment variables starting with `PG_` will be overwritten by the corresponding CLI flag if given.

### Setting the Postgres server's data source name

The PostgreSQL server's [data source name](http://en.wikipedia.org/wiki/Data_source_name)
must be set via the `DATA_SOURCE_NAME` environment variable.

For running it locally on a default Debian/Ubuntu install, this will work (transpose to init script as appropriate):

    sudo -u postgres DATA_SOURCE_NAME="user=postgres host=/var/run/postgresql/ sslmode=disable" postgres_exporter

Also, you can set a list of sources to scrape different instances from the one exporter setup. Just define a comma separated string.

    sudo -u postgres DATA_SOURCE_NAME="port=5432,port=6432" postgres_exporter

See the [github.com/lib/pq](http://github.com/lib/pq) module for other ways to format the connection string.

### Adding new metrics

The exporter will attempt to dynamically export additional metrics if they are added in the
future, but they will be marked as "untyped". Additional metric maps can be easily created
from Postgres documentation by copying the tables and using the following Python snippet:

```python
x = """tab separated raw text of a documentation table"""
for l in StringIO(x):
    column, ctype, description = l.split('\t')
    print """"{0}" : {{ prometheus.CounterValue, prometheus.NewDesc("pg_stat_database_{0}", "{2}", nil, nil) }}, """.format(column.strip(), ctype, description.strip())
```
Adjust the value of the resultant prometheus value type appropriately. This helps build
rich self-documenting metrics for the exporter.

### Adding new metrics via a config file

The -extend.query-path command-line argument specifies a YAML file containing additional queries to run.
Some examples are provided in [queries.yaml](queries.yaml).

### Disabling default metrics
To work with non-officially-supported postgres versions (e.g. 8.2.15),
or variants of postgres (e.g. Greenplum), you can disable the default metrics with the `--disable-default-metrics`
flag. This removes all built-in metrics, and uses only metrics defined by queries in the `queries.yaml` file you supply
(so you must supply one, otherwise the exporter will return nothing but internal statuses and not your database).

### Automatically discover databases
To scrape metrics from all databases on a database server, the database DSN's can be dynamically discovered via the
`--auto-discover-databases` flag. When true, `SELECT datname FROM pg_database WHERE datallowconn = true AND datistemplate = false and datname != current_database()` is run for all configured DSN's. From the
result a new set of DSN's is created for which the metrics are scraped.

In addition, the option `--exclude-databases` adds the possibily to filter the result from the auto discovery to discard databases you do not need.

If you want to include only subset of databases, you can use option `--include-databases`. Exporter still makes request to
`pg_database` table, but do scrape from only if database is in include list.

### Running as non-superuser

To be able to collect metrics from `pg_stat*` views as non-superuser in PostgreSQL
server versions >= 10 you can grant the `pg_monitor` or `pg_read_all_stats` [built-in roles](https://www.postgresql.org/docs/current/predefined-roles.html) to the user. If
you need to monitor older PostgreSQL servers, you will have to create functions
and views as a superuser, and assign permissions separately to those.

```sql
-- To use IF statements, hence to be able to check if the user exists before
-- attempting creation, we need to switch to procedural SQL (PL/pgSQL)
-- instead of standard SQL.
-- More: https://www.postgresql.org/docs/9.3/plpgsql-overview.html
-- To preserve compatibility with <9.0, DO blocks are not used; instead,
-- a function is created and dropped.
CREATE OR REPLACE FUNCTION __tmp_create_user() returns void as $$
BEGIN
  IF NOT EXISTS (
          SELECT                       -- SELECT list can stay empty for this
          FROM   pg_catalog.pg_user
          WHERE  usename = 'postgres_exporter') THEN
    CREATE USER postgres_exporter;
  END IF;
END;
$$ language plpgsql;

SELECT __tmp_create_user();
DROP FUNCTION __tmp_create_user();

ALTER USER postgres_exporter WITH PASSWORD 'password';
ALTER USER postgres_exporter SET SEARCH_PATH TO postgres_exporter,pg_catalog;

-- If deploying as non-superuser (for example in AWS RDS), uncomment the GRANT
-- line below and replace <MASTER_USER> with your root user.
-- GRANT postgres_exporter TO <MASTER_USER>;

GRANT CONNECT ON DATABASE postgres TO postgres_exporter;
```

Run following command if you use PostgreSQL versions >= 10
```sql
GRANT pg_monitor to postgres_exporter;
```

Run following SQL commands only if you use PostgreSQL versions older than 10.
In PostgreSQL, views run with the permissions of the user that created them so
they can act as security barriers. Functions need to be created to share this
data with the non-superuser. Only creating the views will leave out the most
important bits of data.
```sql
CREATE SCHEMA IF NOT EXISTS postgres_exporter;
GRANT USAGE ON SCHEMA postgres_exporter TO postgres_exporter;

CREATE OR REPLACE FUNCTION get_pg_stat_activity() RETURNS SETOF pg_stat_activity AS
$$ SELECT * FROM pg_catalog.pg_stat_activity; $$
LANGUAGE sql
VOLATILE
SECURITY DEFINER;

CREATE OR REPLACE VIEW postgres_exporter.pg_stat_activity
AS
  SELECT * from get_pg_stat_activity();

GRANT SELECT ON postgres_exporter.pg_stat_activity TO postgres_exporter;

CREATE OR REPLACE FUNCTION get_pg_stat_replication() RETURNS SETOF pg_stat_replication AS
$$ SELECT * FROM pg_catalog.pg_stat_replication; $$
LANGUAGE sql
VOLATILE
SECURITY DEFINER;

CREATE OR REPLACE VIEW postgres_exporter.pg_stat_replication
AS
  SELECT * FROM get_pg_stat_replication();

GRANT SELECT ON postgres_exporter.pg_stat_replication TO postgres_exporter;

CREATE EXTENSION IF NOT EXISTS pg_stat_statements;
CREATE OR REPLACE FUNCTION get_pg_stat_statements() RETURNS SETOF pg_stat_statements AS
$$ SELECT * FROM public.pg_stat_statements; $$
LANGUAGE sql
VOLATILE
SECURITY DEFINER;

CREATE OR REPLACE VIEW postgres_exporter.pg_stat_statements
AS
  SELECT * FROM get_pg_stat_statements();

GRANT SELECT ON postgres_exporter.pg_stat_statements TO postgres_exporter;
```

> **NOTE**
> <br />Remember to use `postgres` database name in the connection string:
> ```
> DATA_SOURCE_NAME=postgresql://postgres_exporter:password@localhost:5432/postgres?sslmode=disable
> ```


## Running the tests
```
# Run the unit tests
make test
# Start the test database with docker
docker run -p 5432:5432 -e POSTGRES_DB=circle_test -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=test -d postgres
# Run the integration tests
DATA_SOURCE_NAME='postgresql://postgres:test@localhost:5432/circle_test?sslmode=disable' GOOPTS='-v -tags integration' make test
```

# 🚀 I'm are always open to your feedback.  Please contact as bellow information:
### [Contact ]
* [Name: Nho Luong]
* [Skype](luongutnho_skype)
* [Github](https://github.com/nholuongut/)
* [Linkedin](https://www.linkedin.com/in/nholuong/)
* [Email Address](luongutnho@hotmail.com)
* [PayPal.me](https://www.paypal.com/paypalme/nholuongut)

![](https://i.imgur.com/waxVImv.png)
![](Donate.png)
[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/nholuong)

# License
