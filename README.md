<picture>
 <source media="(prefers-color-scheme: dark)" srcset="docs/images/logo_dark.svg">
 <img alt="light" src="docs/images/logo_light.svg">
</picture>

This package provides access to [dxFeed market data](https://dxfeed.com/market-data/).
The library is built as a language-specific wrapper over
the [dxFeed Graal Native](https://dxfeed.jfrog.io/artifactory/maven-open/com/dxfeed/graal-native-api/) library,
which was compiled with [GraalVM Native Image](https://www.graalvm.org/latest/reference-manual/native-image/)
and [dxFeed Java API](https://docs.dxfeed.com/dxfeed/api/overview-summary.html) (our flagman API).

:warning: It’s a **beta** version and still under active development.

![Build](https://github.com/dxFeed/dxfeed-graal-go-api/actions/workflows/build.yml/badge.svg)
![Platform](https://img.shields.io/badge/platform-win--x64%20%7C%20linux--x64%20%7C%20osx--x64%20%7C%20osx--arm64-lightgrey)
[![Go](https://img.shields.io/badge/go_version-1.20-blueviolet)](https://dotnet.microsoft.com/en-us/)
[![License](https://img.shields.io/badge/license-MPL--2.0-orange)](https://github.com/dxFeed/dxfeed-graal-net-api/blob/master/LICENSE)
[![Release](https://img.shields.io/github/v/release/dxFeed/dxfeed-graal-go-api)](https://github.com/dxFeed/dxfeed-graal-net-api/releases/latest)

## Table of Contents

- [Overview](#overview)
    * [Milestones](#milestones)
    * [Future Development](#future-development)
    * [Implementation Details](#implementation-details)
- [Documentation](#documentation)
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
    * [How to connect to QD endpoint](#how-to-connect-to-QD-endpoint)
    * [How to connect to dxLink](#how-to-connect-to-dxlink)
- [Tools](#tools)
- [Samples](#samples)
- [Current State](#current-state)

## Overview

dxFeed Graal Go API allows developers to create efficient applications in Go language. This enables developers to
leverage all the benefits of native app development, resulting in maximum performance and usability for end users.

### Milestones

We are working diligently to ensure that this new repository meets all of our standards for performance, security, and
scalability. We will be providing regular updates throughout the development process.

If you have any questions, please contact us via
our [customer portal](https://jira.in.devexperts.com/servicedesk/customer/portal/1/create/122).

### Future Development

Features planned with **high priority**:

* refactroing public API
* deployment and publication artifacts in JFrog (to simplify the integration of the package into client applications)
* add all market events
* provide more tools and samples
* generate documentation

---

### Implementation Details

We use [GraalVM Native Image](https://www.graalvm.org/latest/reference-manual/native-image/) technology and specially
written code that *wraps* Java methods into native ones
to get dynamically linked libraries for different platforms (Linux, macOS, and Windows) based on
the [latest Java API package](https://dxfeed.jfrog.io/artifactory/maven-open/com/devexperts/qd/dxfeed-api/).

Then, the resulting dynamic link library (dxFeed Graal-native) is used through
C [ABI](https://en.wikipedia.org/wiki/Application_binary_interface) (application binary interface),
and we write programming interfaces that describe our business model (similar to Java API).

As a result, we get a full-featured, similar performance as with Java API.
Regardless of the language, writing the final application logic using API calls will be very similar (only the syntax
will be amended, *"best practices"*, specific language restrictions)

Below is a scheme of this process:

<picture>
 <source media="(prefers-color-scheme: dark)" srcset="docs/images/scheme_dark.svg">
 <img alt="light" src="docs/images/scheme_light.svg">
</picture>

## Documentation

Find useful information in our self-service dxFeed Knowledge Base or Go API documentation:

- [dxFeed Knowledge Base](https://kb.dxfeed.com/index.html?lang=en)
    * [Getting Started](https://kb.dxfeed.com/en/getting-started.html)
    * [Troubleshooting](https://kb.dxfeed.com/en/troubleshooting-guidelines.html)
    * [Market Events](https://kb.dxfeed.com/en/data-model/dxfeed-api-market-events.html)
    * [Event Delivery contracts](https://kb.dxfeed.com/en/data-model/model-of-event-publishing.html#event-delivery-contracts)
    * [dxFeed API Event classes](https://kb.dxfeed.com/en/data-model/model-of-event-publishing.html#dxfeed-api-event-classes)
    * [Exchange Codes](https://kb.dxfeed.com/en/data-model/exchange-codes.html)
    * [Order Sources](https://kb.dxfeed.com/en/data-model/qd-model-of-market-events.html#order-x)
    * [Order Book reconstruction](https://kb.dxfeed.com/en/data-model/dxfeed-order-book/order-book-reconstruction.html)
    * [Symbology Guide](https://kb.dxfeed.com/en/data-model/symbology-guide.html)

## Requirements

### Windows

Only x64 versions are supported.

| OS                                    | Version        | Architectures |
|---------------------------------------|----------------|---------------|
| [Windows][Windows-client]             | 8, 8.1         | x64           |
| [Windows 10][Windows-client]          | Version 1607+  | x64           |
| [Windows 11][Windows-client]          | Version 22000+ | x64           |
| [Windows Server][Windows-Server]      | 2012+          | x64           |
| [Windows Server Core][Windows-Server] | 2012+          | x64           |
| [Nano Server][Nano-Server]            | Version 1809+  | x64           |

#### Requirements

* [MinGW][mingw]

[Windows-client]: https://www.microsoft.com/windows/

[Windows-Server]: https://learn.microsoft.com/windows-server/

[Nano-Server]: https://learn.microsoft.com/windows-server/get-started/getting-started-with-nano-server

[mingw]: https://www.msys2.org/

### Linux

Only x64 versions are supported.

#### Requirements

* [GCC][gcc]

[gcc]: https://gcc.gnu.org/

#### Libc compatibility

- [glibc][glibc]: 2.35+ (from Ubuntu 22.04)
- [musl][musl]: temporarily unsupported

[glibc]: https://www.gnu.org/software/libc/

[musl]: https://musl.libc.org/

### macOS

| OS             | Version | Architectures |
|----------------|---------|---------------|
| [macOS][macOS] | 10.13+  | x64           |
| [macOS][macOS] | 11+     | Arm64         |

Is supported in the Rosetta 2 x64 emulator.

[macOS]: https://support.apple.com/macos

#### Requirements

* [Clang][clang]

[clang]:https://opensource.apple.com/projects/llvm-clang/

## Installation

### Linux

```bash
sudo apt-get update
sudo apt-get install build-essential
go get github.com/dxFeed/dxFeed-graal-go-api
cd ~/go/pkg/mod/github.com/dxfeed/dxfeed-graal-go-api@<version>
go run .
```

## Usage

### How to connect to QD endpoint

```go
package main

import (
	"fmt"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/quote"
	"math"
	"time"
)

type PrintEvents func(events []interface{})

func (pr PrintEvents) Update(events []any) {
	pr(events)
}

func main() {
	// For token-based authorization, use the following address format:
	// "demo.dxfeed.com:7300[login=entitle:token]"
	endpoint, err := api.NewEndpoint(api.Feed)
	if err != nil {
		panic(err)
	}
	defer func(endpoint *api.DXEndpoint) {
		_ = endpoint.Close()
	}(endpoint)

	err = endpoint.Connect("demo.dxfeed.com:7300")
	if err != nil {
		panic(err)
	}

	feed, err := endpoint.GetFeed()
	if err != nil {
		panic(err)
	}

	subscription, err := feed.CreateSubscription(eventcodes.Quote)
	if err != nil {
		panic(err)
	}
	defer subscription.Close()

	err = subscription.AddListener(PrintEvents(func(events []interface{}) {
		for _, event := range events {
			switch v := event.(type) {
			case *quote.Quote:
				fmt.Printf("%s\n", v.String())
			}
		}
	}))

	err = subscription.AddSymbol("AAPL")
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Duration(math.MaxInt64))
}

```

<details>
<summary>Output</summary>
<br>

```
I 231130 124734.411 [main] QD - Using QDS-3.325+file-UNKNOWN, (C) Devexperts
I 231130 124734.415 [main] QD - Using scheme com.dxfeed.api.impl.DXFeedScheme slfwemJduh1J7ibvy9oo8DABTNhNALFQfw0KmE40CMI
I 231130 124734.418 [main] MARS - Started time synchronization tracker using multicast 239.192.51.45:5145 with dPyAu
I 231130 124734.422 [main] MARS - Started JVM self-monitoring
I 231130 124734.423 [main] QD - monitoring with collectors [Ticker, Stream, History]
I 231130 124734.424 [main] QD - monitoring DXEndpoint with dxfeed.address=demo.dxfeed.com:7300
I 231130 124734.425 [main] ClientSocket-Distributor - Starting ClientSocketConnector to demo.dxfeed.com:7300
I 231130 124734.425 [demo.dxfeed.com:7300-Reader] ClientSocketConnector - Resolving IPs for demo.dxfeed.com
I 231130 124734.427 [demo.dxfeed.com:7300-Reader] ClientSocketConnector - Connecting to 208.93.103.170:7300
I 231130 124734.530 [demo.dxfeed.com:7300-Reader] ClientSocketConnector - Connected to 208.93.103.170:7300
D 231130 124734.634 [demo.dxfeed.com:7300-Reader] QD - Distributor received protocol descriptor multiplexor@fFLro [type=qtp, version=QDS-3.319, opt=hs, mars.root=mdd.demo-amazon.multiplexor-demo1] sending [TICKER, STREAM, HISTORY, DATA] from 208.93.103.170
Quote{AAPL, eventTime=0, time=20231130-123206.000, timeNanoPart=0, sequence=0, bidTime=20231130-123206.000, bidExchange=P, bidPrice=189.36, bidSize=3.0, askTime=20231130-123129.000, askExchange=P, askPrice=189.53, askSize=10.0}
```

</details>

### How to connect to dxLink

```go
package main

import (
	"fmt"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/api"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/eventcodes"
	"github.com/dxfeed/dxfeed-graal-go-api/pkg/events/quote"
	"math"
	"time"
)

type PrintEvents func(events []interface{})

func (pr PrintEvents) Update(events []any) {
	pr(events)
}

func main() {
	// The experimental property must be enabled.
	api.SetSystemProperty("dxfeed.experimental.dxlink.enable", "true")
	// Set scheme for dxLink.
	api.SetSystemProperty("scheme", "ext:opt:sysprops,resource:dxlink.xml")

	// For token-based authorization, use the following address format:
	// "dxlink:wss://demo.dxfeed.com/dxlink-ws[login=dxlink:token]"
	endpoint, err := api.NewEndpoint(api.Feed)
	if err != nil {
		panic(err)
	}
	defer func(endpoint *api.DXEndpoint) {
		_ = endpoint.Close()
	}(endpoint)

	err = endpoint.Connect("dxlink:wss://demo.dxfeed.com/dxlink-ws")
	if err != nil {
		panic(err)
	}

	feed, err := endpoint.GetFeed()
	if err != nil {
		panic(err)
	}

	subscription, err := feed.CreateSubscription(eventcodes.Quote)
	if err != nil {
		panic(err)
	}
	defer subscription.Close()

	err = subscription.AddListener(PrintEvents(func(events []interface{}) {
		for _, event := range events {
			switch v := event.(type) {
			case *quote.Quote:
				fmt.Printf("%s\n", v.String())
			}
		}
	}))

	err = subscription.AddSymbol("AAPL")
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Duration(math.MaxInt64))
}

```

<details>
<summary>Output</summary>
<br>

```
I 231130 124929.817 [main] QD - Using QDS-3.325+file-UNKNOWN, (C) Devexperts
I 231130 124929.821 [main] QD - Using scheme com.dxfeed.api.impl.DXFeedScheme slfwemJduh1J7ibvy9oo8DABTNhNALFQfw0KmE40CMI
I 231130 124929.824 [main] MARS - Started time synchronization tracker using multicast 239.192.51.45:5145 with sWipb
I 231130 124929.828 [main] MARS - Started JVM self-monitoring
I 231130 124929.828 [main] QD - monitoring with collectors [Ticker, Stream, History]
I 231130 124929.829 [main] QD - monitoring DXEndpoint with dxfeed.address=dxlink:wss://demo.dxfeed.com/dxlink-ws
I 231130 124929.831 [main] DxLinkClientWebSocket-Distributor - Starting DxLinkClientWebSocketConnector to wss://demo.dxfeed.com/dxlink-ws
SLF4J: Failed to load class "org.slf4j.impl.StaticLoggerBinder".
SLF4J: Defaulting to no-operation (NOP) logger implementation
SLF4J: See http://www.slf4j.org/codes.html#StaticLoggerBinder for further details.
I 231130 124929.831 [wss://demo.dxfeed.com/dxlink-ws-Writer] DxLinkClientWebSocket-Distributor - Connecting to wss://demo.dxfeed.com/dxlink-ws
I 231130 124930.153 [wss://demo.dxfeed.com/dxlink-ws-Writer] DxLinkClientWebSocket-Distributor - Connected to wss://demo.dxfeed.com/dxlink-ws
D 231130 124931.269 [oioEventLoopGroup-2-1] QD - Distributor received protocol descriptor [type=dxlink, version=0.1-0.18-20231017-133150, keepaliveTimeout=120, acceptKeepaliveTimeout=5] sending [] from wss://demo.dxfeed.com/dxlink-ws
D 231130 124931.271 [oioEventLoopGroup-2-1] QD - Distributor received protocol descriptor [type=dxlink, version=0.1-0.18-20231017-133150, keepaliveTimeout=120, acceptKeepaliveTimeout=5, authentication=] sending [] from wss://demo.dxfeed.com/dxlink-ws
Quote{AAPL, eventTime=0, time=20231130-123421.000, timeNanoPart=0, sequence=0, bidTime=20231130-123421.000, bidExchange=Q, bidPrice=189.47, bidSize=4.0, askTime=20231130-123421.000, askExchange=P, askPrice=189.53, askSize=10.0}
```

</details>

To familiarize with the dxLink protocol, please click [here](https://demo.dxfeed.com/dxlink-ws/debug/#/protocol).

## Tools

[Tools](https://github.com/dxFeed/dxfeed-graal-go-api/)
is a collection of utilities that allow you to subscribe to various market events for the specified symbols. The tools
can
be
downloaded
from [Release](https://github.com/dxFeed/dxfeed-graal-go-api/releases)

* [Connect](https://github.com/dxFeed/dxfeed-graal-go-api/blob/main/cmd/tools/connect.go)
  connects to the specified address(es) and subscribes to the specified events with the specified symbol
* [Dump](https://github.com/dxFeed/dxfeed-graal-go-api/blob/main/cmd/tools/dump.go)
  dumps all events received from address. This was designed to retrieve data from a file
* [PerfTest](https://github.com/dxFeed/dxfeed-graal-go-api/blob/main/cmd/tools/perftest.go)
  connects to the specified address(es) and calculates performance counters (events per second, memory usage, CPU usage,
  etc.)
* [LatencyTest](https://github.com/dxFeed/dxfeed-graal-go-api/blob/main/cmd/tools/latencytest.go)
connects to the specified address(es) and calculates latency

To run tools on macOS, it may be necessary to unquarantine them:

```
sudo /usr/bin/xattr -r -d com.apple.quarantine <directory_with_tools>
```

## Samples

- [ ] ConvertTapeFile demonstrates how to convert one tape file to another tape file with optional intermediate
  processing or filtering
- [ ] DxFeedFileParser is a simple demonstration of how events are read form a tape file
- [ ] DxFeedSample is a simple demonstration of how to create multiple event listeners and subscribe to `Quote`
  and `Trade` events
- [x] [PrintQuoteEvents](https://github.com/dxFeed/dxfeed-graal-go-api/blob/main/cmd/samples/print-quote-events/main.go)
  is a simple demonstration of how to subscribe to the `Quote` event, using a `DxFeed` instance singleton
  and `dxfeed.properties` file
- [ ] WriteTapeFile is a simple demonstration of how to write events to a tape file
- [ ] DxFeedIpfConnect is a simple demonstration of how to get Instrument Profiles
- [ ] DXFeedLiveIpfSample is a simple demonstration of how to get live updates for Instrument Profiles
- [ ] DxFeedPublishProfiles is a simple demonstration of how to publish market events
- [ ] ScheduleSample is a simple demonstration of how to get various scheduling information for instruments
- [ ] DXFeedconnect is a simple demonstration of how to subscribe to different events using TimeSeriesSubscription

## Current State

### Endpoint Roles

- [x] [FEED](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/DXEndpoint.Role.html#FEED)
  connects to the remote data feed provider and is optimized for real-time or delayed data processing,
  **this is a default role**

- [x] [STREAM_FEED](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/DXEndpoint.Role.html#STREAM_FEED)
  is similar to `Feed` and also connects to the remote data feed provider but is designed for bulk data parsing from
  files

- [x] [PUBLISHER](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/DXEndpoint.Role.html#PUBLISHER)
  connects to the remote publisher hub (also known as multiplexor) or creates a publisher on the local host

- [x] [STREAM_PUBLISHER](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/DXEndpoint.Role.html#STREAM_PUBLISHER)
  is similar to `Publisher` and also connects to the remote publisher hub, but is designed for bulk data publishing


- [x] [LOCAL_HUB](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/DXEndpoint.Role.html#LOCAL_HUB)
  is a local hub without the ability to establish network connections. Events published via `Publisher` are delivered to
  local `Feed` only

- [ ] [ON_DEMAND_FEED](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/DXEndpoint.Role.html#ON_DEMAND_FEED)
  is similar to `Feed`, but it is designed to be used
  with  [OnDemandService](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/ondemand/OnDemandService.html) for historical
  data replay only

### Event Types

- [x] [Order](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/market/Order.html)
  is a snapshot of the full available market depth for a symbol

- [x] [SpreadOrder](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/market/SpreadOrder.html)
  is a snapshot of the full available market depth for all spreads

- [x] [AnalyticOrder](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/market/AnalyticOrder.html)
  is an `Order` extension that introduces analytic information, such as adding iceberg-related information to a given
  order

- [x] [Trade](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/market/Trade.html)
  is a snapshot of the price and size of the last trade during regular trading hours and an overall day volume and day
  turnover

- [x] [TradeETH](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/market/TradeETH.html)
  is a snapshot of the price and size of the last trade during extended trading hours and the extended trading hours day
  volume and day turnover

- [x] [Candle](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/candle/Candle.html)
  event with open, high, low, and close prices and other information for a specific period

- [x] [Quote](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/market/Quote.html)
  is a snapshot of the best bid and ask prices and other fields that change with each quote

- [x] [Profile](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/market/Profile.html)
  is a snapshot that contains the security instrument description

- [ ] [Summary](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/market/Summary.html)
  is a snapshot of the trading session, including session highs, lows, etc.

- [x] [TimeAndSale](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/market/TimeAndSale.html)
  represents a trade or other market event with price, such as the open/close price of a market, etc.

- [x] [Greeks](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/option/Greeks.html)
  is a snapshot of the option price, Black-Scholes volatility, and greeks

- [ ] [Series](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/option/Series.html)
  is a snapshot of computed values available for all options series for a given underlying symbol based on options
  market prices

- [ ] [TheoPrice](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/option/TheoPrice.html)
  is a snapshot of the theoretical option price computation that is periodically performed
  by [dxPrice](http://www.devexperts.com/en/products/price.html) model-free computation

- [ ] [Underlying](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/option/Underlying.html)
  is a snapshot of computed values available for an option underlying symbol based on the market’s option prices

- [ ] [OptionSale](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/market/OptionSale.html)
  represents a trade or another market event with the price (for example, market open/close price, etc.) for each option
  symbol listed under the specified `Underlying`

- [ ] [Configuration](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/misc/Configuration.html)
  is an event with an application-specific attachment

- [ ] [Message](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/misc/Message.html)
  is an event with an application-specific attachment

### Subscription Symbols

- [x] [String](https://pkg.go.dev/strings)
  is a string representation of the symbol

- [x] [TimeSeriesSubscriptionSymbol](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/osub/TimeSeriesSubscriptionSymbol.html)
  represents subscription to time-series events

- [x] [IndexedEventSubscriptionSymbol](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/osub/IndexedEventSubscriptionSymbol.html)
  represents subscription to a specific source of indexed events

- [x] [WildcardSymbol.ALL](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/osub/WildcardSymbol.html)
  represents a  *wildcard* subscription to all events of the specific event type

- [ ] [CandleSymbol](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/candle/CandleSymbol.html)
  is a symbol used with [DXFeedSubscription](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/DXFeedSubscription.html)
  class to subscribe for [Candle](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/event/candle/Candle.html) events

### Subscriptions & Models

- [x] [DXFeedSubscription](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/DXFeedSubscription.html)
  is a subscription for a set of symbols and event types

- [ ] [DXFeedTimeSeriesSubscription](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/DXFeedTimeSeriesSubscription.html)
  extends `DXFeedSubscription` to conveniently subscribe to time series events for a set of symbols and event types

- [ ] [ObservableSubscription](https://github.com/devexperts/QD/blob/master/dxfeed-api/src/main/java/com/dxfeed/api/osub/ObservableSubscription.java)
  is an observable set of subscription symbols for the specific event
  type ([Java API sample](https://github.com/devexperts/QD/blob/master/dxfeed-samples/src/main/java/com/dxfeed/sample/_simple_/PublishProfiles.java))

- [ ] [GetLastEvent](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/DXFeed.html#getLastEvent-E-)
  returns the last event for the specified event instance
  ([Java API sample](https://github.com/devexperts/QD/blob/master/dxfeed-samples/src/main/java/com/dxfeed/sample/api/DXFeedSample.java))

- [ ] [GetLastEvents](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/DXFeed.html#getLastEvents-java.util.Collection-)
  returns the last events for the specified event instances list

- [ ] [GetLastEventPromise](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/DXFeed.html#getLastEventPromise-java.lang.Class-java.lang.Object-)
  requests the last event for the specified event type and symbol
  ([Java API sample](https://github.com/devexperts/QD/blob/master/dxfeed-samples/src/main/java/com/dxfeed/sample/console/LastEventsConsole.java))

- [ ] [GetLastEventsPromises](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/DXFeed.html#getLastEventsPromises-java.lang.Class-java.util.Collection-)
  requests the last events for the specified event type and symbol collection

- [ ] [GetLastEventIfSubscribed](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/DXFeed.html#getLastEventIfSubscribed-java.lang.Class-java.lang.Object-)
  returns the last event for the specified event type and symbol if there’s a subscription for it

- [ ] [GetIndexedEventsPromise](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/DXFeed.html#getIndexedEventsPromise-java.lang.Class-java.lang.Object-com.dxfeed.event.IndexedEventSource-)
  requests an indexed events list for the specified event type, symbol, and source

- [ ] [GetIndexedEventsIfSubscribed](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/DXFeed.html#getIndexedEventsIfSubscribed-java.lang.Class-java.lang.Object-com.dxfeed.event.IndexedEventSource-)
  returns a list of indexed events for the specified event type, symbol, and source, if there’s a subscription for it

- [ ] [GetTimeSeriesPromise](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/DXFeed.html#getTimeSeriesPromise-java.lang.Class-java.lang.Object-long-long-)
  requests time series events for the specified event type, symbol, and time range
  ([Java API sample](https://github.com/devexperts/QD/blob/master/dxfeed-samples/src/main/java/com/dxfeed/sample/_simple_/FetchDailyCandles.java))

- [ ] [GetTimeSeriesIfSubscribed](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/api/DXFeed.html#getTimeSeriesIfSubscribed-java.lang.Class-java.lang.Object-long-long-)
  requests time series events for the specified event type, symbol, and time range if there’s a subscription for it

- [ ] [TimeSeriesEventModel](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/model/TimeSeriesEventModel.html)
  is a model of a list of time series events
  ([Java API sample](https://github.com/devexperts/QD/blob/master/dxfeed-samples/src/main/java/com/dxfeed/sample/ui/swing/DXFeedCandleChart.java))

- [ ] [IndexedEventModel](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/model/IndexedEventModel.html)
  is an indexed event list model
  ([Java API sample](https://github.com/devexperts/QD/blob/master/dxfeed-samples/src/main/java/com/dxfeed/sample/ui/swing/DXFeedTimeAndSales.java))

- [ ] [OrderBookModel](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/model/market/OrderBookModel.html)
  is a model of convenient Order Book management
  ([Java API sample](https://github.com/devexperts/QD/blob/master/dxfeed-samples/src/main/java/com/dxfeed/sample/ui/swing/DXFeedMarketDepth.java))

### IPF & Schedule

- [x] [InstrumentProfile](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/ipf/InstrumentProfile.html)
  represents basic profile information about a market instrument
  ([Java API sample](https://github.com/devexperts/QD/blob/master/dxfeed-samples/src/main/java/com/dxfeed/sample/ipf/DXFeedIpfConnect.java))

- [x] [InstrumentProfileReader](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/ipf/InstrumentProfileReader.html) reads
  instrument profiles from the stream using Instrument Profile Format (IPF)

- [ ] [InstrumentProfileCollector](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/ipf/live/InstrumentProfileCollector.html)
  collects instrument profile updates and provides the live instrument profiles list
  ([Java API sample](https://github.com/devexperts/QD/blob/master/dxfeed-samples/src/main/java/com/dxfeed/sample/ipf/DXFeedLiveIpfSample.java))

- [ ] [InstrumentProfileConnection](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/ipf/live/InstrumentProfileConnection.html)
  connects to an instrument profile URL and reads instrument profiles with support of streaming live updates

- [ ] [Schedule](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/schedule/Schedule.html)
  provides API to retrieve and explore various exchanges’ trading schedules and different financial instrument classes
  ([Java API sample](https://github.com/devexperts/QD/blob/master/dxfeed-samples/src/main/java/com/dxfeed/sample/schedule/ScheduleSample.java))

- [ ] [Option Series](https://github.com/devexperts/QD/blob/master/dxfeed-api/src/main/java/com/dxfeed/ipf/option/OptionSeries.java)
  is a series of call and put options with different strike sharing the same attributes of expiration, last trading day,
  spc, multiplies,
  etc. ([Java API sample](https://github.com/devexperts/QD/blob/master/dxfeed-samples/src/main/java/com/dxfeed/sample/ipf/option/DXFeedOptionChain.java))

### Services

- [ ] [OnDemandService](https://docs.dxfeed.com/dxfeed/api/com/dxfeed/ondemand/OnDemandService.html)
  provides on-demand historical tick data replay controls
  ([Java API sample](https://github.com/devexperts/QD/blob/master/dxfeed-samples/src/main/java/com/dxfeed/sample/ondemand/OnDemandSample.java))
