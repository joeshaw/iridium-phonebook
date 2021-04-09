package main

import (
	"fmt"
	"strings"
)

type contact struct {
	name         string
	homeNumber   string
	workNumber   string
	mobileNumber string
	otherNumber  string
	email        string
	notes        string
}

func (c *contact) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "name:%s", c.name)
	if c.homeNumber != "" {
		fmt.Fprintf(&sb, " home:%s", c.homeNumber)
	}
	if c.workNumber != "" {
		fmt.Fprintf(&sb, " work:%s", c.workNumber)
	}
	if c.mobileNumber != "" {
		fmt.Fprintf(&sb, " mobile:%s", c.mobileNumber)
	}
	if c.otherNumber != "" {
		fmt.Fprintf(&sb, " other:%s", c.otherNumber)
	}
	if c.email != "" {
		fmt.Fprintf(&sb, " email:%s", c.email)
	}
	return sb.String()
}

func (c *contact) CSV() string {
	return strings.Join([]string{
		c.name,
		c.homeNumber,
		c.workNumber,
		c.mobileNumber,
		c.otherNumber,
		c.email,
		c.notes,
	}, ",")
}

func contactFromCSV(csv string) *contact {
	var c contact
	fields := strings.Split(csv, ",")
	for i := range fields {
		switch i {
		case 0:
			c.name = fields[i]
		case 1:
			c.homeNumber = fields[i]
		case 2:
			c.workNumber = fields[i]
		case 3:
			c.mobileNumber = fields[i]
		case 4:
			c.otherNumber = fields[i]
		case 5:
			c.email = fields[i]
		case 6:
			c.notes = fields[i]
		}
	}
	return &c
}
