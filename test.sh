#!/usr/bin/env bash
curl 0:4444 -X POST -d '
{
    "addresses": [
        { "ip": "45.180.216.33" },
        { "ip": "45.188.140.1" },
        { "ip": "45.188.140.32342" }
    ]
}'
