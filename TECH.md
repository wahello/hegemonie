# Hegemonie

A Web interface manages the authentication of the players, displays the
status of the country managed by each player and proposes actions (as HTML forms)
to mpact the world.

Meanwhile, the game engine is managed by a standalone daemon that makes
the world evolve with external triggers: long term actions progress a bit
toward their completion, the movements are executed, attacks started, resources
produced, etc etc.

1. Written in 100% in Golang: for the sake of Simplicity and Portability. The
   code mostly depends on [Go-Macaron](https://go-macaron.com) and the Golang
   [standard library](https://golang.org/pkg). At the moment no special
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

A game instance for a small community is lightweight enough to run on a small
ARM-based board.

## Architecture

As of today a single tool called ``hegemonie`` provides all the elements
* **hegemonie front** serves HTTP pages for the human beings
* **hegemonie region** serves a portion of the game's world restricted to a single map region,
  through a simple HTTP/Golang RPC interface.
* **hegemonie round** triggers the rounds in the game's world, and is destined to be triggered
  by ``cron``

## Scalability

This is not the topic yet. However there are already a few opportunities.

The *front* service is stateless, you might deploy many of them.

The *region* service is stateful and it manages all the game entities. Distinct world
services (i.e. processes) will host distinct datasets. A region service is not replicated
and the population pf users requires to be shareded among the regions to grow.

## Deploy with Docker

This is still a work in progress and here can be only one region in the world, because the
``front`` doesn't manage discovery or directory.

Deploy a ``front`` service:
```
docker run --network host jfsmig/hegemonie:latest -- front --north 127.0.0.1:8080 --region 127.0.0.1:8081
```

Deploy a ``region`` service:
```
docker run --network host jfsmig/hegemonie:latest -- region --north 127.0.0.1:8081
```

## Deploy with Snapcraft

TODO

## Try it from scratch

Starting from the sources, if you have the go environment and the ``make`` installed, then simply run:
 
```
make try
```

It will expose a couple of services, bond to ``localhost`` and, respectively, the TCP port ``8080``
for the front and ``8081`` for the only region. Then try the [sandbox](http://127.0.0.1:8080).

As a hint, try to log-in with the user ``admin@hegemonie.be`` and the password ``plop`` ;)
