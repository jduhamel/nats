# nats

A command-line interface to [gnatsd](nats.io).

nats is a command-line utility for sending and receiving messages via
i a gnatsd cluster. It is a thin wrapper around the [NATS Go
client](https://github.com/nats-io/nats).

## Install

```bash
$ go get github.com/soofaloofa/nats
```

## Usage

Publishing to a subject

```bash
$ nats pub subject
test
[#1] Published on [subject] : 'test'
sending
[#2] Published on [subject] : 'sending'
messages
[#3] Published on [subject] : 'messages'
```

Subscribing on a subject

```bash
$ nats sub subject
Listening on [subject]
[#1] Received on [subject]: 'test'
[#2] Received on [subject]: 'sending'
[#3] Received on [subject]: 'messages'
```

## Configuration

Use the `-s` flag to publish or subscribe to a different server.

```bash
$ nats -s tls://192.168.1.45:4222 pub subject
```

comma separate multiple servers

```bash
$ nats -s nats://192.168.1.45:4222,nats://192.168.1.46:4222 pub subject
```
