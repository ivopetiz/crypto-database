[![Codacy Badge](https://api.codacy.com/project/badge/Grade/9243d193bbe34717978b72b0477df4d2)](https://app.codacy.com/app/ivopetiz/crypto-database?utm_source=github.com&utm_medium=referral&utm_content=ivopetiz/crypto-database&utm_campaign=Badge_Grade_Dashboard)

# Crypto-database

Database to store all data from crypto exchanges, currently working with Binance, Bittrex, Cryptopia and Poloniex. 

Can be used for technical analysis, bots, backtest, realtime trading, etc.

## Installation


-   [BD Instalation](#bd-installation)
    -   [Docker](##docker)

    -   [Native](##native)
        -   [BD Configuration](#BD-Configuration)

-   [Market Prices To DB](#market-prices-to-db)
    -   [Data Interval](##data-interval)

-   [Balance To DB](#balance-to-db)

-   [Using Chronograf](#using-chronograf)

-   [TODO](#todo)

This install guide was made for **Ubuntu 16.04+**. Will need some adjustments to work with other distros.

## BD Installation

### Docker

Docker directory has a default configuration that allows users to implement an pre configured database, ready to receive data from Exchanges and use it.
In order to use Influxdb Docker container is only necessary to follow the steps bellow.

```bash
git clone https://github.com/ivopetiz/crypto-database.git
cd crypto-database
. docker/build
. docker/start
```
And your Influxdb crypto database throw Docker container should be ready to be used.

### Native

Start by installing Golang, to build the applications responsible for populate Crypto-database.

```bash
sudo apt-get install golang
```

After install Golang, you will need to install InfluxDB. Chronograf is also recomended.

```bash
sudo apt-get install influxdb influxdb-client chronograf
```

You can also get InfluxDB last versions from Influx website. 
This version of Crypto-database where tested with Influxdb 1.5.3 and Chronograf 1.5.0. You can install these packages by running:

```bash
wget https://dl.influxdata.com/influxdb/releases/influxdb_1.5.3_amd64.deb
sudo dpkg -i influxdb_1.5.3_amd64.deb 
wget https://dl.influxdata.com/chronograf/releases/chronograf_1.5.0.0_amd64.deb
sudo dpkg -i chronograf_1.5.0.0_amd64.deb
```

### BD Configuration

Start Influxdb and then run Influxdb prompt.

```bash
sudo systemctl enable influxdb.service
sudo systemctl start influxdb.service
influx
```

Create databases and user with privileges to use databases:

```influx
CREATE DATABASE altcoin
CREATE DATABASE balance
USE altcoin
CREATE USER <username> WITH PASSWORD '<password>' WITH ALL PRIVILEGES
GRANT ALL PRIVILEGES TO <username>
```

## Market Prices To DB

Before build crypto markets and balance applications is necessary to clone this repository.

```bash
git clone https://github.com/ivopetiz/crypto-database.git
```

 After cloning this rep, is necessary to configure Influxdb user, password and server on system. One option is to define variables on system in order to hide it from git. Variables can also be defined on **consts.go**.

```bash
DBUSER=<your-db-user>
DBPASS=<your-db-password>
SERVERDB=<your-db-server>
```

---
### Data Interval

*By default, data is recorded every 10 seconds but this value can be changed. Timeout can be defined on **consts.go**. A timeout too big won't present fast changes on prices. By the other hand, a timeout too small will make your IP address blocked on crypto exchanges, that only allows a certain number of request per minute. All exchanges have different limits and you can consult these values on exchanges API official websites.*

---

Log file will be stored on **/log/altdb_coin.log**. This path can be changed on **main.go**. Make sure you have the right privileges to write in **/log/**.

```bash
sudo mkdir /log
sudo touch /log/altdb_coin.log
sudo chown <user>:<user> /log/altdb_coin.log
sudo chmod 644 /log/altdb_coin.log
```

Compile market data getter executable.

```bash
cd markets
go build -o markets -ldflags="-s -w" main.go consts.go
```

Go lang will return an executable file called **markets**. Now you need to run in order to populate database. One option is to use **markets** as a service, in order to keep tracking of it. To run **markets** as a service is necessary to create **/etc/systemd/system/cryptomarket.service** with the following content:

```bash
[Unit]
Description=Service running a crypto market data getter.
After=network.target

[Service]
Type=simple
User=user
WorkingDirectory=/home/user/crypto-database
ExecStart=/home/user/crypto-database/markets

Restart=always

[Install]
WantedBy=multi-user.target
```

After save **cryptomarket.service** file, run the commands bellow:

```bash
sudo systemctl daemon-reload
sudo systemctl enable cryptomarket.service
sudo systemctl start cryptomarket.service
```

Now **markets** will run as a service, starting when OS initializes and recovers in case of failure.

## Balance To DB

If you want to add your balance to DB, you will need to generate your API key on crypto exchanges, in order to validate your login and get your data. Currently working only with Bittrex exchange.

To generate an API key and secret, you can access crypto exchanges website, on definitions part. Make sure you keep your API key and API secret **secret** and give this key minimum permissions, in order to obtain just balance info. With bad permissions, any person with this key and secret can consult, trade or withdraw your coins.

After get API key and secret from exchanges, you need to add it to **balance.go**. One option is to define system variables:

```bash
BITTREX_API_KEY=<your-bittrex-api-key>
BITTREX_API_SECRET=<your-bittrex-api-secret>
```

After define **API_KEY** and **API_SECRET** in **balance.go** is necessary to build **balance** program:

```bash
cd balance
go build -o balance -ldflags="-s -w" balance.go
```

Once **balance** is not supposed to run every minute it wasn't built with a timeout, so user need to run it everytime. This can be done automatically with a Crontab rule or as a service, like was done with **markets.service**.

To add new balance info to DB every hour, add a crontab rule by running **crontab -e** and add the following line:

```bash
0 * * * *   /home/user/crypto-database/balance
```

The above rule will run **balance** every hour at :00.

## Using Chronograf

Chronograf presents crypto data from Influxdb. Can be particularly useful to plot data or to quick check market prices. Chronograf is easy to use and it's only necessary to configure with your Influxdb definitions.

In order to start service, after installation will need to run:

```bash
sudo systemctl start chronograf.service
```

If you use Influxdb in the same machine as Chronograf service and have used the default configs, Influxdb will be in https://localhost:8086.

You will also need to input your Influxdb user and password in order to give Chronograf access to DB.

If Influxdb and Chronograf services are in diferent machines, you will need to change localhost by Influxdb machine IP address.

Chronograf service uses port **8888**. Port can be changed on **/etc/default/chronograf**.

Anyone can access Chronograf but is possible to block access from other machines. In config file you can change all configs.

By default, Chronograf config file will be looking like this:

```bash
HOST=0.0.0.0
PORT=8888
```

**HOST** will define who can access Chronograf website. It can be blocked to localhost machine or other specific IP address which can be a good option in terms of security.

---
## TODO

-   add more exchanges to Balance
-   makefile
