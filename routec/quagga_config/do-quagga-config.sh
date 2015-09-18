#/usr/bin/env bash

export GOPATH=/proj/cypress/go
export PATH=$PATH:$GOPATH/bin

quagga_config /proj/cypress/exp/{{.User}}-{{.XP}}/{{.XP}}.route/{{.Id}}.rc.json
