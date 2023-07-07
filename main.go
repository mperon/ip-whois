package main

import (
	"log"
	"net/http"
	"sync"
	"time"
)

const PREFIX string = "IPWHOIS"
const DEFAULT_UPDATE_URL string = "http://ftp.registro.br/pub/numeracao/origin/nicbr-asn-blk-latest.txt"
const DEFAULT_UPDATE_INTERVAL string = "24h"
const MINIMAL_INTERVAL time.Duration = time.Duration(12 * time.Hour)

type CachedDatabase struct {
	db             *IPDatabase
	date           time.Time
	lock           sync.RWMutex
	updateURL      string
	updateInterval time.Duration
}

func NewCachedDatabase() *CachedDatabase {
	cdb := &CachedDatabase{}
	// sets update interval
	interval, err := time.ParseDuration(SafeGetEnvStr(PREFIX+"_UPDATE_INTERVAL", DEFAULT_UPDATE_INTERVAL))
	if err != nil {
		interval = time.Duration(24 * time.Hour)
	}
	if interval < MINIMAL_INTERVAL {
		interval = MINIMAL_INTERVAL
	}
	cdb.updateInterval = interval
	// sets update url
	cdb.updateURL = SafeGetEnvStr(PREFIX+"_UPDATE_URL", DEFAULT_UPDATE_URL)
	//generate a date befor update interval, to force-it
	cdb.date = time.Now().Add(-(cdb.updateInterval + time.Duration(600*time.Second)))
	cdb.lock.Lock()
	cdb.db = NewDatabase()
	cdb.lock.Unlock()
	return cdb
}

func (cdb *CachedDatabase) LoadDatabase() bool {
	today := time.Now()
	nextUpdate := cdb.date.Add(cdb.updateInterval)
	if nextUpdate.After(today) {
		return false
	}
	// try to download new database
	log.Println("Running update...")
	resp, err := http.Get(cdb.updateURL)
	if err != nil {
		log.Printf("was not possible to download from database url.. %s", cdb.updateURL)
		log.Println(err)
		return false
	}
	defer resp.Body.Close()
	newDb := NewDatabase()
	err = newDb.LoadFromReader(resp.Body)
	if err != nil {
		log.Printf("was not possible to load data into database.. %s", cdb.updateURL)
		log.Println(err)
		return false
	}
	// data loaded
	cdb.lock.Lock()
	cdb.date = time.Now()
	cdb.db = newDb
	cdb.lock.Unlock()
	return true
}

func main() {
	cdb := NewCachedDatabase()
	//loading database
	log.Println("Loading Database...")
	cdb.LoadDatabase()
	go doInterval(cdb)
	InitRoutes(cdb)
}

func doInterval(cdb *CachedDatabase) {
	// each 2 hour
	for range time.Tick(time.Hour * 2) {
		// do the interval task
		cdb.LoadDatabase()
	}
}
