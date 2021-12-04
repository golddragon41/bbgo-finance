# bbgo

A trading bot framework written in Go. The name bbgo comes from the BB8 bot in the Star Wars movie. aka Buy BitCoin Go!

## Current Status

[![Build Status](https://travis-ci.org/c9s/bbgo.svg?branch=main)](https://travis-ci.org/c9s/bbgo)

## Features

- Exchange abstraction interface
- Stream integration (user data websocket)
- PnL calculation
- Slack notification
- KLine-based backtest
- Built-in strategies
- Multi-session support
- Standard indicators (SMA, EMA, BOLL)
- React-powered Web Dashboard

## Supported Exchanges

- MAX Spot Exchange (located in Taiwan)
- Binance Spot Exchange
- FTX Spot Exchange

## Requirements

Get your exchange API key and secret after you register the accounts (you can choose one or more exchanges):

- For MAX: <https://max.maicoin.com/signup?r=c7982718>
- For Binance: <https://www.binancezh.com/en/register?ref=VGDGLT80>
- For FTX: <https://ftx.com/#a=7710474>

Since the exchange implementation and support are done by a small team, if you like the work they've done for you, It
would be great if you can use their referral code as your support to them. :-D

## Installation

### Install from binary

The following script will help you set up a config file, dotenv file:

```sh
bash <(curl -s https://raw.githubusercontent.com/c9s/bbgo/main/scripts/setup-grid.sh)
```

### Install and Run from the One-click Linode StackScript:

- BBGO USDT/TWD Market Grid Trading <https://cloud.linode.com/stackscripts/793380>
- BBGO USDC/TWD Market Grid Trading <https://cloud.linode.com/stackscripts/797776>
- BBGO LINK/TWD Market Grid Trading <https://cloud.linode.com/stackscripts/797774>
- BBGO USDC/USDT Market Grid Trading <https://cloud.linode.com/stackscripts/797777>
- BBGO Standard Grid Trading <https://cloud.linode.com/stackscripts/795788>

### Install from source

If you need to use go-sqlite, you will need to enable CGO first:

```
CGO_ENABLED=1 go get github.com/mattn/go-sqlite3
```

Install the bbgo command:

```sh
go get -u github.com/c9s/bbgo/cmd/bbgo
```

Add your dotenv file:

```sh
# if you have one
BINANCE_API_KEY=
BINANCE_API_SECRET=

# if you have one
MAX_API_KEY=
MAX_API_SECRET=

# if you have one
FTX_API_KEY=
FTX_API_SECRET=
# specify it if credentials are for subaccount
FTX_SUBACCOUNT=
```

Prepare your dotenv file `.env.local` and BBGO yaml config file `bbgo.yaml`.

The minimal bbgo.yaml could be generated by:

```sh
curl -o bbgo.yaml https://raw.githubusercontent.com/c9s/bbgo/main/config/minimal.yaml
```

To sync your own trade data:

```sh
bbgo sync --session max
bbgo sync --session binance
```

If you want to switch to other dotenv file, you can add an `--dotenv` option or `--config`:

```sh
bbgo sync --dotenv .env.dev --config config/grid.yaml --session binance
```

To sync remote exchange klines data for backtesting:

```sh
bbgo backtest --exchange binance -v --sync --sync-only --sync-from 2020-01-01
```

To run backtest:

```sh
bbgo backtest --exchange binance --base-asset-baseline
```

To query transfer history:

```sh
bbgo transfer-history --session max --asset USDT --since "2019-01-01"
```

To calculate pnl:

```sh
bbgo pnl --exchange binance --asset BTC --since "2019-01-01"
```

To run strategy:

```sh
bbgo run
```

To start bbgo with the frontend dashboard:

```sh
bbgo run --enable-webserver
```

## Advanced Setup

### Setting up Telegram Bot Notification

Open your Telegram app, and chat with @botFather

Enter `/newbot` to create a new bot

Enter the bot display name. ex. `your_bbgo_bot`

Enter the bot username. This should be global unique. e.g., `bbgo_bot_711222333`

Botfather will response your a bot token. *Keep bot token safe*

Set `TELEGRAM_BOT_TOKEN` in the `.env.local` file, e.g.,

```sh
TELEGRAM_BOT_TOKEN=347374838:ABFTjfiweajfiawoejfiaojfeijoaef
```

For the telegram chat authentication (your bot needs to verify it's you), if you only need a fixed authentication token,
you can set `TELEGRAM_AUTH_TOKEN` in the `.env.local` file, e.g.,

```sh
TELEGRAM_BOT_AUTH_TOKEN=itsme55667788
```

Run your bbgo,

Open your Telegram app, search your bot `bbgo_bot_711222333`

Enter `/start` and `/auth {code}`

Done! your notifications will be routed to the telegram chat.

### Setting up Slack Notification

Put your slack bot token in the .env.local file:

```sh
SLACK_TOKEN=xxoox
```

### Synchronizing Trading Data

By default, BBGO does not sync your trading data from the exchange sessions, so it's hard to calculate your profit and
loss correctly.

By synchronizing trades and orders to the local database, you can earn some benefits like PnL calculations, backtesting
and asset calculation.

#### Configure MySQL Database

To use MySQL database for data syncing, first you need to install your mysql server:

```sh
# For Ubuntu Linux
sudo apt-get install -y mysql-server
```

Or [run it in docker](https://hub.docker.com/_/mysql)

Create your mysql database:

```sh
mysql -uroot -e "CREATE DATABASE bbgo CHARSET utf8"
```

Then put these environment variables in your `.env.local` file:

```sh
DB_DRIVER=mysql
DB_DSN="user:password@tcp(127.0.0.1:3306)/bbgo"
```

#### Configure Sqlite3 Database

Just put these environment variables in your `.env.local` file:

```sh
DB_DRIVER=sqlite3
DB_DSN=bbgo.sqlite3
```

## Built-in Strategies

Check out the strategy directory [strategy](pkg/strategy) for all built-in strategies:

- `pricealert` strategy demonstrates how to use the notification system [pricealert](pkg/strategy/pricealert)
- `xpuremaker` strategy demonstrates how to maintain the orderbook and submit maker
  orders [xpuremaker](pkg/strategy/xpuremaker)
- `buyandhold` strategy demonstrates how to subscribe kline events and submit market
  order [buyandhold](pkg/strategy/pricedrop)
- `bollgrid` strategy implements a basic grid strategy with the built-in bollinger
  indicator [bollgrid](pkg/strategy/bollgrid)
- `grid` strategy implements the fixed price band grid strategy [grid](pkg/strategy/grid)
- `flashcrash` strategy implements a strategy that catches the flashcrash [flashcrash](pkg/strategy/flashcrash)

To run these built-in strategies, just modify the config file to make the configuration suitable for you, for example if
you want to run
`buyandhold` strategy:

```sh
vim config/buyandhold.yaml

# run bbgo with the config
bbgo run --config config/buyandhold.yaml
```

## Adding New Built-in Strategy

Fork and clone this repository, Create a directory under `pkg/strategy/newstrategy`, write your strategy
at `pkg/strategy/newstrategy/strategy.go`.

Define a strategy struct:

```go
package newstrategy

import (
	"github.com/c9s/bbgo/pkg/fixedpoint"
)

type Strategy struct {
	Symbol string           `json:"symbol"`
	Param1 int              `json:"param1"`
	Param2 int              `json:"param2"`
	Param3 fixedpoint.Value `json:"param3"`
}
```

Register your strategy:

```go
package newstrategy

const ID = "newstrategy"

const stateKey = "state-v1"

var log = logrus.WithField("strategy", ID)

func init() {
	bbgo.RegisterStrategy(ID, &Strategy{})
}
```

Implement the strategy methods:

```go
package newstrategy

func (s *Strategy) Subscribe(session *bbgo.ExchangeSession) {
	session.Subscribe(types.KLineChannel, s.Symbol, types.SubscribeOptions{Interval: "2m"})
}

func (s *Strategy) Run(ctx context.Context, orderExecutor bbgo.OrderExecutor, session *bbgo.ExchangeSession) error {
	// ....
	return nil
}
```

Edit `pkg/cmd/builtin.go`, and import the package, like this:

```go
package cmd

// import built-in strategies
import (
	_ "github.com/c9s/bbgo/pkg/strategy/bollgrid"
	_ "github.com/c9s/bbgo/pkg/strategy/buyandhold"
	_ "github.com/c9s/bbgo/pkg/strategy/flashcrash"
	_ "github.com/c9s/bbgo/pkg/strategy/grid"
	_ "github.com/c9s/bbgo/pkg/strategy/pricealert"
	_ "github.com/c9s/bbgo/pkg/strategy/support"
	_ "github.com/c9s/bbgo/pkg/strategy/swing"
	_ "github.com/c9s/bbgo/pkg/strategy/trailingstop"
	_ "github.com/c9s/bbgo/pkg/strategy/xmaker"
	_ "github.com/c9s/bbgo/pkg/strategy/xpuremaker"
)
```

## Write your own strategy

Create your go package, and initialize the repository with `go mod` and add bbgo as a dependency:

```sh
go mod init
go get github.com/c9s/bbgo@main
```

Write your own strategy in the strategy file:

```sh
vim strategy.go
```

You can grab the skeleton strategy from <https://github.com/c9s/bbgo/blob/main/pkg/strategy/skeleton/strategy.go>

Now add your config:

```sh
mkdir config
(cd config && curl -o bbgo.yaml https://raw.githubusercontent.com/c9s/bbgo/main/config/minimal.yaml)
```

Add your strategy package path to the config file `config/bbgo.yaml`

```yaml
---
build:
  dir: build
  imports:
  - github.com/your_id/your_swing
  targets:
  - name: swing-amd64-linux
    os: linux
    arch: amd64
  - name: swing-amd64-darwin
    os: darwin
    arch: amd64
```

Run `bbgo run` command, bbgo will compile a wrapper binary that imports your strategy:

```sh
dotenv -f .env.local -- bbgo run --config config/bbgo.yaml
```

Or you can build your own wrapper binary via:

```shell
bbgo build --config config/bbgo.yaml
```

## Command Usages

### Submitting Orders to a specific exchagne session

```shell
bbgo submit-order --session=okex --symbol=OKBUSDT --side=buy --price=10.0 --quantity=1
```

### Listing Open Orders of a specific exchange session

```sh
bbgo list-orders open --session=okex --symbol=OKBUSDT
bbgo list-orders open --session=ftx --symbol=FTTUSDT
bbgo list-orders open --session=max --symbol=MAXUSDT
bbgo list-orders open --session=binance --symbol=BNBUSDT
```

### Canceling an open order

```shell
# both order id and symbol is required for okex
bbgo cancel-order --session=okex --order-id=318223238325248000 --symbol=OKBUSDT

# for max, you can just give your order id
bbgo cancel-order --session=max --order-id=1234566
```

### Debugging user data stream

```shell
bbgo userdatastream --session okex
bbgo userdatastream --session max
bbgo userdatastream --session binance
```

## Dynamic Injection

In order to minimize the strategy code, bbgo supports dynamic dependency injection.

Before executing your strategy, bbgo injects the components into your strategy object if it found the embedded field
that is using bbgo component. for example:

```go
type Strategy struct {
*bbgo.Notifiability
}
```

And then, in your code, you can call the methods of Notifiability.

Supported components (single exchange strategy only for now):

- `*bbgo.Notifiability`
- `bbgo.OrderExecutor`

If you have `Symbol string` field in your strategy, your strategy will be detected as a symbol-based strategy, then the
following types could be injected automatically:

- `*bbgo.ExchangeSession`
- `types.Market`

## Strategy Execution Phases

1. Load config from the config file.
2. Allocate and initialize exchange sessions.
3. Add exchange sessions to the environment (the data layer).
4. Use the given environment to initialize the trader object (the logic layer).
5. The trader initializes the environment and start the exchange connections.
6. Call strategy.Run() method sequentially.

## Exchange API Examples

Please check out the example directory: [examples](examples)

Initialize MAX API:

```go
key := os.Getenv("MAX_API_KEY")
secret := os.Getenv("MAX_API_SECRET")

maxRest := maxapi.NewRestClient(maxapi.ProductionAPIURL)
maxRest.Auth(key, secret)
```

Creating user data stream to get the orderbook (depth):

```go
stream := max.NewStream(key, secret)
stream.Subscribe(types.BookChannel, symbol, types.SubscribeOptions{})

streambook := types.NewStreamBook(symbol)
streambook.BindStream(stream)
```

## How To Add A New Exchange

(TBD)

## Helm Chart

If you need redis:

```sh
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install redis bitnami/redis
```

To get the dynamically generated redis password, you can use the following command:

```sh
export REDIS_PASSWORD=$(kubectl get secret --namespace bbgo redis -o jsonpath="{.data.redis-password}" | base64 --decode)
```

Prepare your docker image locally (you can also use the docker image from docker hub):

```sh
make docker DOCKER_TAG=1.16.0
```

The docker tag version number is from the file [Chart.yaml](charts/bbgo/Chart.yaml)

Choose your instance name:

```sh
export INSTANCE=grid
```

Prepare your secret:

```sh
kubectl create secret generic bbgo-$INSTANCE --from-env-file .env.local
```

Configure your config file, the chart defaults to read config/bbgo.yaml to create a configmap:

```sh
cp config/grid.yaml bbgo-$INSTANCE.yaml
vim bbgo-$INSTANCE.yaml
```

Prepare your configmap:

```sh
kubectl create configmap bbgo-$INSTANCE --from-file=bbgo.yaml=bbgo-$INSTANCE.yaml
```

Install chart with the preferred release name, the release name maps to the previous secret we just created, that
is, `bbgo-grid`:

```sh
helm install --set existingConfigmap=bbgo-$INSTANCE bbgo-$INSTANCE ./charts/bbgo
```

To use the latest version:

```sh
helm install --set existingConfigmap=bbgo-$INSTANCE --set image.tag=latest bbgo-$INSTANCE ./charts/bbgo
```

To upgrade:

```sh
helm upgrade bbgo-$INSTANCE ./charts/bbgo
helm upgrade --set image.tag=1.15.2 bbgo-$INSTANCE ./charts/bbgo
```

Delete chart:

```sh
helm delete bbgo-$INSTANCE
```

## Development

The overview function flow at bbgo
![image info](./assets/overview.svg)

### Setting up your local repository

1. Click the "Fork" button from the GitHub repository.
2. Clone your forked repository into `$GOPATH/github.com/c9s/bbgo`.
3. Change directory into `$GOPATH/github.com/c9s/bbgo`.
4. Create a branch and start your development.
5. Test your changes.
6. Push your changes to your fork.
7. Send a pull request.

### Adding new migration

1. The project used rockerhopper for db migration. 
https://github.com/c9s/rockhopper


2. Create migration files

```sh
rockhopper --config rockhopper_sqlite.yaml create --type sql add_pnl_column
rockhopper --config rockhopper_mysql.yaml create --type sql add_pnl_column
```

or

```
bash utils/generate-new-migration.sh add_pnl_column
```

Be sure to edit both sqlite3 and mysql migration files. ( [Sample] (migrations/mysql/20210531234123_add_kline_taker_buy_columns.sql) )


To test the drivers, you have to update the rockhopper_mysql.yaml file to connect your database,
then do:

```sh
rockhopper --config rockhopper_sqlite.yaml up
rockhopper --config rockhopper_mysql.yaml up
```

Then run the following command to compile the migration files into go files:

```shell
make migrations
```

### Setup frontend development environment

```sh
cd frontend
yarn install
```

### Testing Desktop App

for webview

```sh
make embed && go run -tags web ./cmd/bbgo-webview
```

for lorca

```sh
make embed && go run -tags web ./cmd/bbgo-lorca
```


## Contributing

See [Contributing](./CONTRIBUTING.md)

## Community

You can join our telegram channels:

- BBGO International <https://t.me/bbgo_intl>
- BBGO Taiwan <https://t.me/bbgocrypto>

## License

MIT License
