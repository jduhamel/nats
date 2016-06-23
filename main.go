package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/nats-io/nats"
	"github.com/urfave/cli"
)

var (
	url   string // NATS server url
	reply string // Publish reply-to subject
)

func main() {
	app := cli.NewApp()
	app.Name = "nats"
	app.Usage = "a nats.io CLI"
	app.Version = "0.0.1"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Kevin Sookocheff",
			Email: "kevin.sookocheff@gmail.com",
		},
	}
	app.Copyright = "(c) 2016 Kevin Sookocheff"

	// TODO: Add additional config options as needed
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "server, s",
			Usage:       "NATS server url and port. Separate multiple servers with a comma.",
			Value:       nats.DefaultURL,
			EnvVar:      "NATS_SERVER",
			Destination: &url,
		},
	}

	// Pub/Sub commands
	app.Commands = []cli.Command{
		cli.Command{
			Name:    "pub",
			Aliases: []string{"a"},
			Usage:   "publish messages to a subject",
			Action:  pub,
		},
		cli.Command{
			Name:    "sub",
			Aliases: []string{"s"},
			Usage:   "subscribe to a subject",
			Action:  sub,
		},
	}

	app.Run(os.Args)
}

// splitFirstWord splits the input string at the first word
func splitFirstWord(s string) (string, string) {
	for i := range s {
		// If we encounter a space, reduce the count.
		if s[i] == ' ' {
			return s[0:i], s[i:]
		}
	}
	// Return the entire string.
	return s, ""
}

func pubUsage() string {
	return "Usage: <subject<message>"
}

func pub(c *cli.Context) error {
	nc, err := nats.Connect(url)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Can't connect to %s: %v", url, err), 1)
	}
	defer nc.Close()

	if len(c.Args()) != 1 {
		return cli.NewExitError("Usage: nats pub <subject>", 1)
	}

	subj, i := c.Args()[0], 0
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		i++
		msg := strings.TrimSpace(scanner.Text())
		if msg == "" {
			continue
		}
		msg = strings.TrimSpace(msg)
		nc.Publish(subj, []byte(msg))
		nc.Flush()

		fmt.Printf("[#%d] Published on [%s] : '%s'\n", i, subj, msg)
	}

	return nil
}

func sub(c *cli.Context) error {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Can't connect to %s: %v", nats.DefaultURL, err), 1)
	}
	defer nc.Close()

	if len(c.Args()) < 1 || len(c.Args()) > 2 {
		return cli.NewExitError("Usage: nats sub <subject> [queue group]", 1)
	}

	subj, i := c.Args()[0], 0
	var queue string
	if len(c.Args()) == 2 {
		queue = c.Args()[2]
	}
	nc.QueueSubscribe(subj, queue, func(msg *nats.Msg) {
		i++
		fmt.Printf("[#%d] Received on [%s]: '%s'\n", i, msg.Subject, string(msg.Data))
	})

	fmt.Printf("Listening on [%s]\n", subj)

	// Loop until CTRL-D
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
	}

	return nil
}
