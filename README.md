# IP-Whois - IP Whois Database Server Search

Build on golang, this is an web application to identify owners of ip addresses on registro.br database.

The software uses public registro.br data to identify.

Writting in go, is very lightweight and fast.

## Example:
```bash
curl serverhost:port -X POST -d '
{
    "addresses": [
        { "ip": "45.180.216.33" },
        { "ip": "45.188.140.1" },
        { "ip": ""45.188.140.32342"" }
    ]
}'
```

The result will be:
```json
{
   "results":{
      "45.180.216.33":{
         "code":"AS269194",
         "name":"STAR1 INTERNET",
         "document":"20.241.468/0001-85"
      },
      "45.188.140.1":{
         "code":"AS269532",
         "name":"ROUTINGER INTERNET E TI",
         "document":"32.915.048/0001-16"
      }
   },
   "errors":{
      "45.188.140.32342":"invalid ip address: 45.188.140.32342"
   }
}
```
