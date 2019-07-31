#!/bin/bash

clear
clear
cp -a pfcp/*/* pfcp
go build main.go
rm pfcp/*.go
#./main
