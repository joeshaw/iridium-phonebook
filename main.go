package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/warthog618/modem/at"
	"github.com/warthog618/modem/serial"
	"github.com/warthog618/modem/trace"
)

type usageError string

func (err usageError) Error() string { return string(err) }
func (err usageError) Unwrap() error { return flag.ErrHelp }

func main() {
	root := &ffcli.Command{
		Name:       "iridium-phonebook",
		ShortUsage: "iridium-phonebook <subcommand> [flags]",
		Exec: func(context.Context, []string) error {
			return usageError("no command given")
		},
		Subcommands: []*ffcli.Command{
			dumpCmd(),
			loadCmd(),
			clearCmd(),
		},
	}

	if err := root.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr)
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
}

func dumpCmd() *ffcli.Command {
	var device string

	fs := flag.NewFlagSet("dump", flag.ExitOnError)
	fs.StringVar(&device, "d", "", "Serial device to use (ie, /dev/cu.usbmodem1)")

	return &ffcli.Command{
		Name:       "dump",
		ShortUsage: "iridium-phonebook dump -d <device>",
		ShortHelp:  "Dump the Iridium phonebook in CSV format to stdout",
		FlagSet:    fs,
		Exec: func(ctx context.Context, args []string) error {
			if device == "" {
				return usageError("-d <device> must be provided")
			}

			modem, cleanup, err := getModem(device)
			if err != nil {
				return err
			}
			defer cleanup()

			pb := newPhonebook(modem)
			for pb.Next() {
				fmt.Println(pb.Contact().CSV())
			}

			if err := pb.Err(); err != nil {
				return err
			}

			return nil
		},
	}
}

func loadCmd() *ffcli.Command {
	var device string

	fs := flag.NewFlagSet("load", flag.ExitOnError)
	fs.StringVar(&device, "d", "", "Serial device to use (ie, /dev/cu.usbmodem1)")

	return &ffcli.Command{
		Name:       "load",
		ShortUsage: "iridium-phonebook load -d <device> <file>",
		ShortHelp:  "Load a CSV file info the Iridium phonebook, without replacing existing contacts",
		FlagSet:    fs,
		Exec: func(ctx context.Context, args []string) error {
			if device == "" {
				return usageError("-d <device> must be provided")
			}

			if len(args) < 1 {
				return usageError("missing file to load from")
			}

			var rc io.ReadCloser
			if args[0] == "-" {
				rc = ioutil.NopCloser(os.Stdin)
			} else {
				f, err := os.Open(args[0])
				if err != nil {
					return err
				}
				rc = f
			}
			defer rc.Close()

			modem, cleanup, err := getModem(device)
			if err != nil {
				return err
			}
			defer cleanup()

			pb := newPhonebook(modem)

			s := bufio.NewScanner(rc)
			for s.Scan() {
				c := contactFromCSV(s.Text())
				if c.name == "" {
					return fmt.Errorf("missing name from input")
				}
				if err := pb.WriteContact(c); err != nil {
					return fmt.Errorf("unable to write %q to phonebook: %w", c.name, err)
				}
				fmt.Printf("Loaded %v to phonebook\n", c)
			}
			if err := s.Err(); err != nil {
				return err
			}

			return nil
		},
	}
}

func clearCmd() *ffcli.Command {
	var device string

	fs := flag.NewFlagSet("clear", flag.ExitOnError)
	fs.StringVar(&device, "d", "", "Serial device to use (ie, /dev/cu.usbmodem1)")

	return &ffcli.Command{
		Name:       "clear",
		ShortUsage: "iridium-phonebook clear -d <device>",
		ShortHelp:  "Delete all contacts from the Iridium phonebook",
		FlagSet:    fs,
		Exec: func(ctx context.Context, args []string) error {
			if device == "" {
				return usageError("-d <device> must be provided")
			}

			modem, cleanup, err := getModem(device)
			if err != nil {
				return err
			}
			defer cleanup()

			pb := newPhonebook(modem)
			if err := pb.DeleteAllContacts(); err != nil {
				return err
			}

			return nil
		},
	}
}

func getModem(device string) (*at.AT, func(), error) {
	cleanup := func() {}

	port, err := serial.New(serial.WithPort(device), serial.WithBaud(9600))
	if err != nil {
		return nil, cleanup, err
	}

	var p io.ReadWriter = port
	if os.Getenv("MODEM_DEBUG") != "" {
		p = trace.New(p)
	}

	modem := at.New(p)
	if err := modem.Init(); err != nil {
		port.Close()
		return nil, cleanup, err
	}

	cleanup = func() {
		port.Close()
	}

	return modem, cleanup, err
}
