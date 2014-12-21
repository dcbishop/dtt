DoTheThing
==========

[![Build Status](https://drone.io/github.com/dcbishop/dtt/status.png)](https://drone.io/github.com/dcbishop/dtt/latest)
[![License](http://img.shields.io/badge/license-CC0-red.svg)](https://creativecommons.org/publicdomain/zero/1.0/legalcode)

Moves files to a target directories based on predefined regex rules.

It's useful for moving files when your done with them, such as putting a frequently download files into an archival location after finishing looking at them.

Installation
------------
    go get -u github.com/dcbishop/dtt

Configuration
-------------
DoTheThing uses a configuration file stored in $XDG_CONFIG_HOME (~/.config/dothething/rules.yaml).

File entries are based on Go's [regexp](http://golang.org/pkg/regexp/syntax/).

    ---
    - file: SomeFile
      move: /mnt/somewhere/else

    - file: (?i)CaseInsensative
      move: /mnt/target/location

    - file: Dots.Are.Any.Characters.Could.Be.Spaces.UnderScore.OrActualDots
      move: /mnt/target/location
    
Usage
-----
    dtt SomeFiles*
