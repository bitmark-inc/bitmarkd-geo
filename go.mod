module github.com/araujobsd/bitmarkd-geo

go 1.12

replace (
	github.com/testcontainers/testcontainer-go v0.0.2 => github.com/testcontainers/testcontainers-go v0.0.2
	github.com/testcontainers/testcontainer-go v0.0.4 => github.com/testcontainers/testcontainers-go v0.0.4
)

require (
	github.com/NYTimes/gziphandler v1.1.1
	github.com/OneOfOne/xxhash v1.2.5 // indirect
	github.com/Wessie/appdirs v0.0.0-20141031215813-6573e894f8e2 // indirect
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/bradfitz/gomemcache v0.0.0-20190329173943-551aad21a668 // indirect
	github.com/coreos/etcd v3.3.13+incompatible // indirect
	github.com/dgryski/go-sip13 v0.0.0-20190329191031-25c5027a8c7b // indirect
	github.com/flopp/go-coordsparser v0.0.0-20160810104536-845bca739e26 // indirect
	github.com/flopp/go-staticmaps v0.0.0-20190722115053-456a5d548ba1
	github.com/fogleman/gg v1.3.0
	github.com/go-redsync/redsync v1.2.0 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/geo v0.0.0-20190916061304-5b978397cfec
	github.com/micro/go-micro v1.10.0
	github.com/mmcloughlin/globe v0.0.0-20190613033401-6df15c77b44b
	github.com/nytimes/gziphandler v1.1.1
	github.com/pborman/uuid v1.2.0 // indirect
	github.com/prometheus/tsdb v0.8.0 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/stretchr/testify v1.4.0 // indirect
	github.com/tidwall/pinhole v0.0.0-20170713004337-171cd602c428 // indirect
	github.com/tkrajina/gpxgo v1.0.1 // indirect
	go.etcd.io/etcd v3.3.13+incompatible // indirect
	golang.org/x/image v0.0.0-20190910094157-69e4b8554b2a // indirect
	golang.org/x/sys v0.0.0-20190924062700-2aa67d56cdd7 // indirect
	gopkg.in/bsm/ratelimit.v1 v1.0.0-20160220154919-db14e161995a // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/redis.v3 v3.6.4 // indirect
)

replace github.com/NYTimes/gziphandler v1.1.1 => github.com/nytimes/gziphandler v1.1.1

replace github.com/nytimes/gziphandler v1.1.1 => github.com/NYTimes/gziphandler v1.1.1

replace (
	github.com/golang/lint v0.0.0-20190313153728-d0100b6bd8b3 => golang.org/x/lint v0.0.0-20190313153728-d0100b6bd8b3
	github.com/golang/lint v0.0.0-20190409202823-959b441ac422 => golang.org/x/lint v0.0.0-20190409202823-959b441ac422
)

replace golang.org/x/lint v0.0.0-20190409202823-959b441ac422 => golang.org/x/lint v0.0.0-20190409202823-959b441ac422
