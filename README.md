[![GoDoc](https://godoc.org/github.com/bitmark-inc/bitmarkd-geo/plugins?status.svg)](https://godoc.org/github.com/bitmark-inc/bitmarkd-geo/)
[![GitHub issues](https://img.shields.io/github/issues/bitmark-inc/bitmarkd-geo.svg)](https://github.com/bitmark-inc/bitmarkd-geo/issues)
[![GitHub forks](https://img.shields.io/github/forks/bitmark-inc/bitmarkd-geo.svg)](https://github.com/bitmark-inc/bitmarkd-geo/network)
[![Go Report Card](https://goreportcard.com/badge/github.com/bitmark-inc/bitmarkd-geo)](https://goreportcard.com/report/github.com/bitmark-inc/bitmarkd-geo)
[![CircleCI](https://circleci.com/gh/bitmark-inc/bitmarkd-geo.svg?style=svg)](https://circleci.com/gh/bitmark-inc/bitmarkd-geo)

bitmarkd-geo
================
This software maps all Bitmark Inc. nodes in the live network.
It is used at: [Bitmark-Nodes](https://nodes.bitmark.com)

## Build and run instructions
1) `make release`
2) `cd standalone && make release`
3) `cp <binary> /usr/local/bin/bitmarkd-geo-cmd`
4) `cd ../`
5) `cp config/example/config.yaml /tmp/`
6) `daemon ./bitmarkdgeo-freebsd-amd64-1.0`

## Copyright and licensing
Distributed under [2-Clause BSD License](https://github.com/araujobsd/aws-icinga2-generator/blob/master/LICENSE).
