package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/warthog618/modem/at"
)

type phonebook struct {
	modem   *at.AT
	entries int
	pos     int
	cur     *contact
	err     error
}

func newPhonebook(modem *at.AT) *phonebook {
	return &phonebook{
		modem:   modem,
		entries: -1,
	}
}

func (p *phonebook) Next() bool {
	if p.err != nil {
		return false
	}

	if p.entries == -1 {
		count, err := p.countEntries()
		if err != nil {
			p.err = err
			return false
		}

		p.entries = count
	}

	if p.pos < p.entries {
		contact, err := p.readEntry()
		if err != nil {
			p.err = err
			return false
		}
		p.cur = contact
		p.pos++

		return true
	}

	return false
}

func (p *phonebook) countEntries() (int, error) {
	output, err := p.modem.Command("+CAPBR=?")
	if err != nil {
		return 0, err
	}

	// Output is in the format
	// +CAPBR: <n>
	// where <n> is the number of entries in the phonebook
	const header = "+CAPBR: "

	if len(output) < 1 || len(output[0]) < len(header) {
		return 0, fmt.Errorf("short output reading phonebook entries")
	}

	return strconv.Atoi(output[0][len(header):])
}

func (p *phonebook) readEntry() (*contact, error) {
	output, err := p.modem.Command(fmt.Sprintf("+CAPBR=%d", p.pos))
	if err != nil {
		return nil, err
	}

	// Output is in the format
	// +CAPBR: <name>,<home_number>,<work_number>,<mobile_number> ,<other_number>,<email>,<notes>
	// where each field is a hex string of the big-endian UCS-2 encoded string

	const header = "+CAPBR: "

	if len(output) < 1 || len(output[0]) < len(header) {
		return nil, fmt.Errorf("short output reading phonebook entries")
	}

	var c contact

	fields := strings.Split(output[0][len(header):], ",")
	for i, f := range fields {
		s, err := decodeUCS2Hex(f)
		if err != nil {
			return nil, fmt.Errorf("unable to decode %q: %w", f, err)
		}

		switch i {
		case 0:
			c.name = s
		case 1:
			c.homeNumber = s
		case 2:
			c.workNumber = s
		case 3:
			c.mobileNumber = s
		case 4:
			c.otherNumber = s
		case 5:
			c.email = s
		case 6:
			c.notes = s
		}
	}

	return &c, nil
}

func (p *phonebook) Contact() *contact {
	return p.cur
}

func (p *phonebook) Err() error {
	return p.err
}

func (p *phonebook) WriteContact(c *contact) error {
	fields := []string{
		c.name,
		c.homeNumber,
		c.workNumber,
		c.mobileNumber,
		c.otherNumber,
		c.email,
		c.notes,
	}
	ucs2Fields := make([]string, len(fields))
	for i := range fields {
		ucs2Hex, err := encodeUCS2Hex(fields[i])
		if err != nil {
			return err
		}
		ucs2Fields[i] = ucs2Hex
	}

	// This command has no output
	command := "+CAPBW=" + strings.Join(ucs2Fields, ",")
	_, err := p.modem.Command(command, at.WithTimeout(10*time.Second))

	return err
}

func (p *phonebook) DeleteContact(pos int) error {
	// This command has no output
	_, err := p.modem.Command(fmt.Sprintf("+CAPBD=%d", pos))
	return err
}

func (p *phonebook) DeleteAllContacts() error {
	// This command has no output
	_, err := p.modem.Command("+CAPBD=ALL")
	return err
}
