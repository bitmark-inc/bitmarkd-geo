/*-
 * Copyright 2019 Bitmark, Inc.
 * Copyright 2019 by Marcelo Araujo <araujo@FreeBSD.org>
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted providing that the following conditions
 * are met:
 * 1. Redistributions of source code must retain the above copyright
 *    notice, this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright
 *    notice, this list of conditions and the following disclaimer in the
 *    documentation and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE AUTHOR ``AS IS'' AND ANY EXPRESS OR
 * IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED.  IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY
 * DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
 * OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT,
 * STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING
 * IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 *
 */

package geolocation

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
	"github.com/mmcloughlin/globe"
)

var (
	ndDistribution = make(map[string]int)
	green          = color.NRGBA{0x00, 0x54, 0x3c, 60}
	red            = color.NRGBA{255, 0, 0, 255}
	plat           float64
	plon           float64
	imgPath        = "webserver/mysite/img/"
	WasItRotate    = 0
)

// FlatMap - Create the context for the flatmap
func FlatMap() (flatmap *sm.Context) {
	flatmap = sm.NewContext()
	flatmap.SetSize(800, 600)
	flatmap.SetCenter(s2.LatLngFromDegrees(8.7832, 34.5085))

	return flatmap
}

// FlatMapAddMarker - Add markers in the flatmap
func FlatMapAddMarker(flatmap *sm.Context, lat float64, lon float64) {
	skey := fmt.Sprintf("%f,%f", lat, lon)
	if _, ok := ndDistribution[skey]; ok {
		ndDistribution[skey] += 1
	} else {
		ndDistribution[skey] = 1
	}

	plat = lat
	plon = lon

	flatmap.AddMarker(sm.NewMarker(s2.LatLngFromDegrees(lat, lon),
		green, 01.0+float64(ndDistribution[skey])))
}

// FlatMapRender - Render the image for the flat map
func FlatMapRender(flatmap *sm.Context) (err error) {
	img, err := flatmap.Render()
	if err != nil {
		return err
	}

	err = gg.SavePNG(imgPath+"flatmap.png", img)
	if err != nil {
		return err
	}

	return nil
}

// GlobeMap - Create the context for the globemap
func GlobeMap() (globemap *globe.Globe) {
	globemap = globe.New()
	globemap.DrawGraticule(10.0)
	globemap.DrawLandBoundaries()

	return globemap
}

// GlobeMapAddMarker - Add markers in the globemap
func GlobeMapAddMarker(globemap *globe.Globe, lat float64, lon float64) {
	skey := fmt.Sprintf("%f,%f", lat, lon)
	if _, ok := ndDistribution[skey]; ok {
		ndDistribution[skey] += 1
	} else {
		ndDistribution[skey] = 1
	}

	precision, _ := strconv.ParseFloat(fmt.Sprintf("%0.2f",
		float64(ndDistribution[skey])/60), 64)

	globemap.DrawDot(lat, lon, 0.01+precision, globe.Color(green))

	if plat != 0 && plon != 0 {
		globemap.DrawLine(lat, lon, plat, plon, globe.Color(red))
	}

	plat = lat
	plon = lon
}

// GlobeMapRender - Render the image for the globe map
func GlobeMapRender(globemap *globe.Globe, myWanLat float64, myWanLon float64) {
	if WasItRotate == 0 {
		globemap.CenterOn(myWanLat, myWanLon)
		WasItRotate = 1
	}
	globemap.CenterOn(myWanLat, myWanLon)
	globemap.SavePNG(imgPath+"globe.png", 600)
}
