package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type StatusCount struct {
	Requests  int
	Addresses int
	Success   int
	NotFound  int
	sync.RWMutex
}

type IPQuery struct {
	Addresses []IPEntry `json:"addresses" binding:"required"`
}

type IPEntry struct {
	Address string `json:"ip" uri:"address" binding:"required"`
}

type IPResult struct {
	Results map[string]*Company `json:"results"`
	Errors  map[string]string   `json:"errors"`
}

func InitRoutes(cdb *CachedDatabase) {
	statusCount := StatusCount{
		Requests:  0,
		Addresses: 0,
	}
	r := gin.Default()

	// https://github.com/gin-gonic/gin/issues/3336
	r.RemoteIPHeaders = append(r.RemoteIPHeaders, "True-Client-IP")
	r.ForwardedByClientIP = true
	r.SetTrustedProxies([]string{"127.0.0.1", "::1"})

	//listen
	r.POST("/", func(c *gin.Context) {
		var f IPQuery
		if err := c.BindJSON(&f); err != nil {
			Error(c, fmt.Errorf("invalid body"), http.StatusBadRequest)
			return
		}
		// checa limites
		if len(f.Addresses) > 50 {
			Error(c, fmt.Errorf("too many addresses"), http.StatusBadRequest)
		}
		// armazena os enderecos aqui
		reply := IPResult{
			Results: make(map[string]*Company),
			Errors:  make(map[string]string),
		}
		addrCount := len(f.Addresses)
		okCount := 0
		errCount := 0
		//para cada ip, faz uma consulta e retorna o resultado
		for _, address := range f.Addresses {
			//verifica se o ip ja foi processado
			if _, found := reply.Results[address.Address]; found {
				log.Printf("ip already searched: %v", address.Address)
				continue
			}
			// nao foi, procura no banco de dados
			cdb.lock.RLock()
			company, err := cdb.db.Search(address.Address)
			cdb.lock.RUnlock()
			if err != nil {
				errCount += 1
				reply.Errors[address.Address] = err.Error()
			} else {
				okCount += 1
				reply.Results[address.Address] = company
			}
		}
		statusCount.Lock()
		statusCount.Addresses += addrCount
		statusCount.Success += okCount
		statusCount.NotFound += errCount
		statusCount.Requests++
		statusCount.Unlock()
		c.JSON(http.StatusOK, reply)
	})

	r.GET("/:address", func(c *gin.Context) {
		addrCount := 1
		okCount := 0
		errCount := 0
		// armazena os enderecos aqui
		reply := IPResult{
			Results: make(map[string]*Company),
			Errors:  make(map[string]string),
		}
		var address IPEntry
		if err := c.BindUri(&address); err != nil {
			Error(c, fmt.Errorf("invalid parameter"), http.StatusBadRequest)
			return
		}
		cdb.lock.RLock()
		company, err := cdb.db.Search(address.Address)
		cdb.lock.RUnlock()
		if err != nil {
			errCount += 1
			reply.Errors[address.Address] = err.Error()
		} else {
			okCount += 1
			reply.Results[address.Address] = company
		}
		statusCount.Lock()
		statusCount.Addresses += addrCount
		statusCount.Success += okCount
		statusCount.NotFound += errCount
		statusCount.Requests++
		statusCount.Unlock()
		c.JSON(http.StatusOK, reply)
	})

	r.GET("/status", func(c *gin.Context) {
		statusCount.RLock()
		defer statusCount.RUnlock()
		c.IndentedJSON(http.StatusOK, gin.H{
			"status": gin.H{
				"running":   true,
				"requests":  statusCount.Requests,
				"addresses": statusCount.Addresses,
				"success":   statusCount.Success,
				"notfound":  statusCount.NotFound,
			},
			"database": gin.H{
				"companies":   len(cdb.db.Companies),
				"last_update": cdb.date,
				"next_update": cdb.date.Add(cdb.updateInterval),
			},
			"source": gin.H{
				"url":      cdb.updateURL,
				"interval": cdb.updateInterval.String(),
			},
		})
	})
	r.Run("127.0.0.1:" + SafeGetEnvStr(PREFIX+"_PORT", "4444"))
}

func Error(c *gin.Context, err error, code int) {
	if code == 0 {
		code = 500
	}
	c.JSON(code, gin.H{"error": err.Error()})
}
