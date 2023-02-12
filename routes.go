package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IPQuery struct {
	Addresses []IPEntry `json:"addresses" binding:"required"`
}

type IPEntry struct {
	Address string `json:"ip" binding:"required"`
}

type IPResult struct {
	Results map[string]*Company `json:"results"`
	Errors  map[string]string   `json:"errors"`
}

func InitRoutes(cdb *CachedDatabase) {
	r := gin.Default()
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
				reply.Errors[address.Address] = err.Error()
			} else {
				reply.Results[address.Address] = company
			}
		}
		c.JSON(http.StatusOK, reply)
	})
	r.Run(":4000")
}

func Error(c *gin.Context, err error, code int) {
	if code == 0 {
		code = 500
	}
	c.JSON(code, gin.H{"error": err.Error()})
}
