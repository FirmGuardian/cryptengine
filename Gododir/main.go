/**
 * This file is much like a gulpfile. I chose Godo mainly in keeping with the
 * theme of "keep it Go." I doesn't have dependencies like Node, npm, Gulp,
 * and a bunch of other Javascript stuff.
 *
 * Please see documentation here: https://github.com/go-godo/godo
 */

package main

import (
	do "gopkg.in/godo.v2"
)

func tasks(p *do.Project) {
	do.Env = `GOPATH=.vendor::$GOPATH`

	p.Task("default", do.S{"build"}, nil)

	// Buildy bits
	p.Task("build", do.S{"clean", "buildall"}, nil)

	p.Task("buildall", do.P{"build_darwin", "build_windows", "build_rpi"}, nil)

	p.Task("build_windows", do.P{"build_win32", "build_win64"}, nil)

	p.Task("build_darwin", nil, func(c *do.Context) {
		c.Run("GOOS=darwin GOARCH=amd64 gb build all")
	}).Src("src/**/*.go")

	p.Task("build_win32", nil, func(c *do.Context) {
		c.Run("GOOS=windows GOARCH=386 gb build all")
	}).Src("src/**/*.go")

	p.Task("build_win64", nil, func(c *do.Context) {
		c.Run("GOOS=windows GOARCH=amd64 gb build all")
	}).Src("src/**/*.go")

	p.Task("build_rpi", nil, func(c *do.Context) {
		c.Run("GOOS=linux GOARCH=arm gb build all")
	}).Src("src/**/*.go")

	// Vendor Updates
	p.Task("blind_update", nil, func(c *do.Context) {
		c.Run("gb vendor update --all")
	})

	// Bin and pkg dir cleanup (think `make clean`)
	p.Task("clean", nil, func(c *do.Context) {
		c.Run("rm -Rf pkg/* && rm -f bin/*")
	})
}

func main() {
	do.Godo(tasks)
}
