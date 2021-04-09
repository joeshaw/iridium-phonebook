# iridium-phonebook

This is a tool for managing your Iridium 9555 satellite phone's
addressbook over a USB connection to your computer.

The phone exposes Hayes-like AT commands over a serial interface for all
sorts of things, including managing the phone's addressbook.

This is known to work on a Mac.  It should work on Linux and Windows
too.

## Installation

To install `iridium-phonebook`, with Go 1.16 or newer:

    go install github.com/joeshaw/iridium-phonebook@latest

With Go 1.15 or older:

    go get -u github.com/joeshaw/iridium-phonebook

## Usage

```
USAGE
  iridium-phonebook <subcommand> [flags]

SUBCOMMANDS
  dump   Dump the Iridium phonebook in CSV format to stdout
  load   Load a CSV file info the Iridium phonebook, without replacing existing contacts
  clear  Delete all contacts from the Iridium phonebook
```

All commands require a `-d <device>` flag, for example:

    iridium-phonebook dump -d /dev/cu.usbmodem1

Note that the `load` command does not do any kind of de-duplication, so
if you dump and load the CSV without running `clear` you will end up
with double contacts.

## Known issues

My phone reliably crashes and reboots after loading about 10-20
contacts.  Still working on tracking down why.

The phone sends contact information as UCS-2.  For decoding UCS-2 is a
strict subset of UTF-16, but not for encoding.  UCS-2 is a fixed-width
encoding, whereas UTF-16 is variable width with surrogate pairs.  I am
encoding as UTF-16 and things will break if you try to use a Unicode
code point that requires a surrogate pair.  So don't try to encode emoji
or things will break.

## License

Copyright 2021 Joe Shaw

`iridium-phonebook` is licensed under the MIT License.  See the LICENSE
file for details.
