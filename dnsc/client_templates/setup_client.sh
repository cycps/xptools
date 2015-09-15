#!/usr/bin/env bash

cp head /etc/resolveconf/resolv.conf.d/
resolvconf -u
