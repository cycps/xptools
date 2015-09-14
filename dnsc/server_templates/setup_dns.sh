#!/usr/bin/env bash

cp named.conf /etc/bind/
cp named.conf.local /etc/bind
cp named.conf.options /etc/bind
cp keys.conf /etc/bind
mkdir -p /etc/bind/zones
cp db.{{.Xpname}}.cypress.net /etc/bind/zones/
cp usr.sbin.named /etc/apparmor.d/local/

sudo service bind9 restart
