#!/bin/bash
# Author: Ropon
# Email: luopeng@codoon.com

ServiceName=ops_golang
LogDir=/var/log/go_log/

exec /opt/${ServiceName} >>${LogDir}/${ServiceName}.out 2>&1
