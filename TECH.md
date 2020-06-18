# Hegemonie

## Architecture

The Hegemonie platform consist in a set of microservices.
  * **web server** is the front service that does the web UI: it authenticates the
    users and serves simple HTML/CSS pages to depict the status of the countries
    manager by the players, as well as it propsoes HTTP forms to make it evolve.
  * **auth server** is a ``grpc`` service responsible of the user management
    (mostly C.R.U.D. operations plus password checks). The service itself is not
    authenticated. Please refer for the
    [Auth API](https://github.com/jfsmig/hegemonie/blob/master/pkg/auth/service.proto)
    API for more information.
  * **region server** is a tuple of 3 ``grpc`` services responsible for the game
    logic within a single region. The service itself is not authenticated. Please refer to
    the ``grpc`` description of the
    [City API](https://github.com/jfsmig/hegemonie/blob/master/pkg/region/city.proto),
    [Army API](https://github.com/jfsmig/hegemonie/blob/master/pkg/region/army.proto) and
    [Admin API](https://github.com/jfsmig/hegemonie/blob/master/pkg/region/admin.proto)
    for more information.
  * **events server** is a ``grpc`` service proposing to subscribe to events related
    to specific topics (e.g. a country, a region). The service itself is not authenticated.
  * **api server** exposes an authenticated API that agregates and consolidates all the
    other API. The users are authenticated and the usage authorized with OAuth2. 


![Hegemonie Architecture](https://raw.githubusercontent.com/jfsmig/hegemonie/master/docs/system-architecture.png)

Everything is controlled by a single CLI tool, ``hegemonie``. That single CLI carries 
the 5 servers and their respective 5 clients. 

 1. Written in 100% in Golang: for the sake of Simplicity and Portability. The
    code mostly depends on [Go-Macaron](https://go-macaron.com) (for ``web server``),
    [Golang grpc](https://github.com/grpc/grpc-go) (for all the other internal services)
    and the Golang [standard library](https://golang.org/pkg). At the moment no special
    attention has been paid to the performance of the whole thing: this will
    happen after the release of a very first MVP.
 2. No database required: the system has all its components in RAM while it is
    alive, it periodically persist its state and restore it at the startup.
    The status is written in [JSON](https://json.org) to ease the daily
    administration.
 3. Notifications will be emitted upon special events in the game.
    No technical solution has been chosen yet.
    It is likely to be split into a collect by either [Redis](https://redis.io),
    [Kafka](https://kafka.apache.org) or [Beanstalkd](https://beanstalkd.github.io),
    and then forwarded to any IM (instant messenging) application like
    [Discord](https://discord.io/), [Slack](https://slack.com),
    [RocketChat](https://rocket.chat), [Riot](https://riot.im) or whatever.


## Scalability

This is not the topic yet. However there are already a few opportunities to let
the game scale:
  * ``web server`` is stateless because it relies on sessions in a side ``memcached``.
    It can be scaled as necessary. 
  * ``auth server`` is currently stateful because it relies on a local storage. Further
    scaling plans exist, either based on a sharding of the users or on a scalable storage
    backend. This is still to be discussed and is not a topic yet.
  * ``event server`` is currently stateful because it relies on a local storage. Further
    scaling plans exist, either based on a sharding of the users or on a scalable storage
    backend. This is still to be discussed and is not a topic yet.
  * ``region server`` is stateful and it manages all the entities in-game. Distinct world
    services (i.e. processes) will host distinct datasets. A region service represents an
    opportunity to shard the users.
  * ``api server`` is currently a vaporware so it scales without any limit (lol).
    By design the API server will be stateless. It will scale seamlessly.

Whatever the solution in place, only the ``web server`` will require an external load
balancing, at the ingress. It is likely that hegemonie will soon rely on a service mesh
to route the grpc messages to the appropriate targets.


## Reliability

This is not a topic yet.


## Performance

Further than the design choice, in the region server, to keep each region "live in RAM",
the performance is not a topic yet.

We roughly target a system that can manage a game instance for a small community of less than
50 players, that would be lightweight enough to run on a ARM-based single board computer (e.g.
a RaspberryPi v3).


## Deploy with Docker

This is still a work in progress and here can be only one region in the world, because
the ``web server`` doesn't manage the discovery or a directory of regions (it only relies
on CLI options).

With the help of the subsequent alias:
```
alias HEGE='docker run --network host jfsmig/hegemonie:latest --'
```

Deploy a ``front`` service:
```
HEGE web server --endpoint 127.0.0.1:8080 --region 127.0.0.1:8081
```

Deploy a ``region`` service:
```
HEGE region server --endpoint 127.0.0.1:8081
```

## Deploy with Snapcraft

TODO


## Try it from scratch

Starting from the sources, if you have the go environment and the ``make`` installed,
then simply run:

```
set -e
set -x
BASE=github.com/jfsmig/hegemonie
go get "${BASE}"
go mod download "${BASE}"
cd "${BASE}"
make try
```

It will expose a couple of services, bond to ``localhost`` and, respectively, the TCP port ``8080``
for the front and ``8081`` for the only region. Then try the [sandbox](http://127.0.0.1:8080).

As a hint, try to log-in with the user ``admin@hegemonie.be`` and the password ``plop`` ;)
