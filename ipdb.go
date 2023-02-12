package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/netip"
	"os"
	"strconv"
	"strings"
)

type Company struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	Document string `json:"document"`
}

type IPRange struct {
	Prefix  netip.Prefix
	Company *Company
}

type IPStruct struct {
	Key       int64
	Separator string
	Children  map[uint64]*IPStruct
	Ranges    []IPRange
}

type IPDatabase struct {
	Companies map[string]*Company
	Ipv4      *IPStruct
	Ipv6      *IPStruct
}

// func main() {
// 	db := NewDatabase()
// 	db.LoadFromFile("data/nicbr-asn-blk-latest.txt")

// 	// check some addresses
// 	company, err := db.Search("45.180.216.1")
// 	if err != nil {
// 		fmt.Printf("err: %v\n", err)
// 		return
// 	}
// 	fmt.Println("The Company was found!:", company)
// }

func NewDatabase() *IPDatabase {
	return &IPDatabase{
		Companies: make(map[string]*Company, 0),
		Ipv4:      CreateIPStruct("."),
		Ipv6:      CreateIPStruct(":"),
	}
}

func (d *IPDatabase) LoadFromFile(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	return d.LoadFromReader(f)
}

func (d *IPDatabase) LoadFromReader(reader io.Reader) error {
	csvReader := csv.NewReader(reader)
	csvReader.Comma = '|'
	csvReader.FieldsPerRecord = -1
	records, err := csvReader.ReadAll()
	if err != nil {
		return err
	}
	for _, fields := range records {
		d.processLine(fields)
	}
	return nil
}

func (d *IPDatabase) processLine(fields []string) {
	//Field 0: Company code
	//Field 1: Company name
	//Field 2: Company cnpj
	//Field 3..: IP Ranges
	var (
		code    string = fields[0]
		company *Company
		found   bool
	)
	if company, found = d.Companies[code]; !found {
		//company need to be created
		company = &Company{
			Code:     code,
			Name:     fields[1],
			Document: fields[2],
		}
		d.Companies[code] = company
	}
	// percorre os campos e processa
	for _, ipRange := range fields[3:] {
		d.processRange(company, ipRange)
	}
}

func (d *IPDatabase) Search(address string) (*Company, error) {
	//convert ip address to object
	ipAddr, err := netip.ParseAddr(address)
	if err != nil {
		return nil, fmt.Errorf("invalid ip address")
	}
	// now make the search
	var ipStruct *IPStruct
	if ipAddr.Is4() {
		ipStruct = d.Ipv4
	} else {
		ipStruct = d.Ipv6
	}
	// now use struct to process thing
	return ipStruct.Search(ipAddr)
}

func (d *IPDatabase) processRange(company *Company, ipRange string) {
	ipPrefix, err := netip.ParsePrefix(ipRange)
	if err != nil {
		fmt.Println("Invalid IP Prefix", ipRange)
		return
	}
	var ipStruct *IPStruct
	if ipPrefix.Addr().Is4() {
		ipStruct = d.Ipv4
	} else {
		ipStruct = d.Ipv6
	}
	// now use struct to process thing
	ipStruct.Add(ipPrefix, company)

}

func CreateIPStruct(Separator string) *IPStruct {
	return &IPStruct{
		Separator: Separator,
		Children:  make(map[uint64]*IPStruct, 0),
		Ranges:    make([]IPRange, 0, 10),
	}
}

func (st *IPStruct) Search(ipAddr netip.Addr) (*Company, error) {
	// nao achou, tenta navegar
	parts := strings.Split(ipAddr.String(), st.Separator)
	return st.traverse(ipAddr, parts)
}

func (st *IPStruct) traverse(ipAddr netip.Addr, parts []string) (*Company, error) {
	// check the current ranges
	for _, ipRange := range st.Ranges {
		if ipRange.Prefix.Contains(ipAddr) {
			return ipRange.Company, nil
		}
	}
	// go deep one step
	if len(parts) > 1 {
		//acessa
		ip_id, err := strconv.ParseUint(parts[0], 16, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot convert ip address to number")
		}
		v, found := st.Children[ip_id]
		if found {
			return v.traverse(ipAddr, parts[1:])
		}
	}
	return nil, fmt.Errorf("IP address not found in this database")
}

func (st *IPStruct) Add(prefix netip.Prefix, company *Company) {
	parts := strings.Split(prefix.Addr().String(), st.Separator)
	st.parsePrefix(prefix, company, parts)
}

func (st *IPStruct) parsePrefix(prefix netip.Prefix, company *Company, parts []string) {
	if len(parts) == 0 {
		return
	}
	if parts[0] == "" || parts[0] == "0" {
		//fim da linha/adiciona aqui
		st.Ranges = append(st.Ranges, IPRange{
			Company: company,
			Prefix:  prefix,
		})
		return
	}
	// navega na estrutura
	ip_id, err := strconv.ParseUint(parts[0], 16, 64)
	if err != nil {
		panic("Cannot convert string to uint number!")
	}
	pStruct, exists := st.Children[ip_id]
	if !exists {
		//cria a estrutura
		pStruct = CreateIPStruct(st.Separator)
		pStruct.Key = int64(ip_id)
		st.Children[ip_id] = pStruct
	}
	pStruct.parsePrefix(prefix, company, parts[1:])
}

func (c *Company) String() string {
	return fmt.Sprintf("{Code:%s, Name:%s, Document:%s}", c.Code, c.Name, c.Document)
}

/*
func main() {
	cogen := &Company{ // pb == &Student{"Bob", 8}
		Code: "AS174",
		Name: "COGENT BRASIL TELECOMUNICAÇÕES LTDA.",
		Cnpj: "29.484.413/0001-70",
	}

	ipv6 := &IPStruct{
		Key: "",
	}

	range1 := &IPRange{
		Prefix: netip.MustParsePrefix("2804:5330::/32"),
		Company: cogen,
	}

	part1, err := strconv.ParseInt("2804", 16, 64)
	if err != nil {
		panic(err)
	}

}
*/
