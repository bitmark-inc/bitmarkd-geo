module github.com/araujobsd/bitmarkd-geo

go 1.12

replace (
	github.com/testcontainers/testcontainer-go v0.0.2 => github.com/testcontainers/testcontainers-go v0.0.2
	github.com/testcontainers/testcontainer-go v0.0.4 => github.com/testcontainers/testcontainers-go v0.0.4
)

require (
	github.com/NYTimes/gziphandler v1.1.1
	github.com/Wessie/appdirs v0.0.0-20141031215813-6573e894f8e2 // indirect
	github.com/flopp/go-coordsparser v0.0.0-20160810104536-845bca739e26 // indirect
	github.com/flopp/go-staticmaps v0.0.0-20180404185116-320790ed5329
	github.com/fogleman/gg v1.3.0
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/geo v0.0.0-20190507233405-a0e886e97a51
	github.com/micro/go-micro v1.5.0
	github.com/mmcloughlin/globe v0.0.0-20180909115233-4175779440e5
	github.com/nytimes/gziphandler v1.1.1
	github.com/tidwall/pinhole v0.0.0-20170713004337-171cd602c428 // indirect
	github.com/tkrajina/gpxgo v1.0.1 // indirect
)

replace github.com/NYTimes/gziphandler v1.1.1 => github.com/nytimes/gziphandler v1.1.1

replace github.com/nytimes/gziphandler v1.1.1 => github.com/NYTimes/gziphandler v1.1.1

replace (
	github.com/golang/lint v0.0.0-20190313153728-d0100b6bd8b3 => golang.org/x/lint v0.0.0-20190313153728-d0100b6bd8b3
	github.com/golang/lint v0.0.0-20190409202823-959b441ac422 => golang.org/x/lint v0.0.0-20190409202823-959b441ac422
)

replace golang.org/x/lint v0.0.0-20190409202823-959b441ac422 => golang.org/x/lint v0.0.0-20190409202823-959b441ac422
