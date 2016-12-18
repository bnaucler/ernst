package main

import (
	"fmt"
	"strings"
	"strconv"
	"math/rand"
	"encoding/json"
	"github.com/thoj/go-ircevent"
	"github.com/boltdb/bolt"
	elib "github.com/bnaucler/ernst/elib"
)

func csetlist(event *irc.Event, settings *elib.Settings) (resp string) {

	resp = fmt.Sprintf("%v: rate: %d/%d, kdel: %d/%d, randel: %d/%d, dnrmem: %d/%d",
		event.Nick, settings.Rate, elib.Ratemax,
		settings.Kdel, elib.Kdelmax, settings.Randel, elib.Randelmax,
		settings.Dnrmem, elib.Dnrmemmax)

	return
}

func csetshow(event *irc.Event, settings *elib.Settings, setvar string) (resp string) {

	if setvar == "rate" {
		resp = fmt.Sprintf("%v: %v: %d/%d", event.Nick,
			setvar, settings.Rate, elib.Ratemax)
	} else if setvar == "kdel" {
		resp = fmt.Sprintf("%v: %v: %d/%d", event.Nick,
			setvar, settings.Kdel, elib.Kdelmax)
	} else if setvar == "randel" {
		resp = fmt.Sprintf("%v: %v: %d/%d", event.Nick,
			setvar, settings.Randel, elib.Randelmax)
	} else if setvar == "dnrmem" {
		resp = fmt.Sprintf("%v: %v: %d/%d", event.Nick,
			setvar, settings.Dnrmem, elib.Dnrmemmax)
	} else {
		resp = fmt.Sprintf("%v: %v", event.Nick, erresp)
	}

	return
}

func csetset(event *irc.Event, settings *elib.Settings, setvar,
	setval string) (resp string, dbchange bool) {

	nval, nerr := strconv.Atoi(setval)

	if setvar == "rate"  && nerr == nil && nval > -1 && nval <= elib.Ratemax {
		if settings.Rate != nval { dbchange = true }
		settings.Rate = nval
		resp = fmt.Sprintf("%v: %s: %d/%d", event.Nick, setvar,
			settings.Rate, elib.Ratemax)

	} else if setvar == "kdel" && nerr == nil && nval > -1 && nval <= elib.Kdelmax {
		if settings.Kdel != nval { dbchange = true }
		settings.Kdel = nval
		resp = fmt.Sprintf("%v: %s: %d/%d", event.Nick, setvar,
			settings.Kdel, elib.Kdelmax)

	} else if setvar == "randel" && nerr == nil && nval > -1 && nval <= elib.Randelmax {
		if settings.Randel != nval { dbchange = true }
		settings.Randel = nval
		resp = fmt.Sprintf("%v: %s: %d/%d", event.Nick, setvar,
			settings.Randel, elib.Randelmax)

	} else if setvar == "dnrmem" && nerr == nil && nval > -1 && nval <= elib.Dnrmemmax {
		if settings.Dnrmem != nval { dbchange = true }
		settings.Dnrmem = nval
		resp = fmt.Sprintf("%v: %s: %d/%d", event.Nick, setvar,
			settings.Dnrmem, elib.Dnrmemmax)

	} else {
		resp = fmt.Sprintf("%v: %v", event.Nick, erresp)
	}

	return
}

func cset(irccon *irc.Connection, db *bolt.DB, rnd *rand.Rand,
	event *irc.Event, settings *elib.Settings) bool {

	ssp := strings.Split(event.Arguments[1], " ")

	var (
		setvar, setval, resp string
		dbchange bool
	)

	if len(ssp) > 1 { setvar = strings.ToLower(ssp[1]) }
	if len(ssp) > 2 { setval = ssp[2] }

	if len(ssp) == 1 { resp = csetlist(event, settings)
	} else if len(ssp) == 2 { resp = csetshow(event, settings, setvar)
	} else if len(ssp) == 3 {
		resp, dbchange = csetset (event, settings, setvar, setval)
	}

	if len(resp) != 0 {
		if dbchange {
			s, err:= json.Marshal(settings)
			elib.Cherr(err)
			err = elib.Wrdb(db, 0, s)
			elib.Cherr(err)
		}
		irccon.Privmsg(settings.Channel, resp)
	}

	return true
}
