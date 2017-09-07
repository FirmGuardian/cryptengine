/**
 * This file is much like a gulpfile. I chose Godo mainly in keeping with the
 * theme of "keep it Go." I doesn't have dependencies like Node, npm, Gulp,
 * and a bunch of other Javascript stuff.
 *
 * Please see documentation here: https://github.com/go-godo/godo
 */

package main

import (
	//"fmt"
	do "gopkg.in/godo.v2"
)

func tasks(p *do.Project) {
	do.Env = `GOPATH=.vendor::$GOPATH`

	p.Task("default", do.S{"buildall"}, nil)

	// Buildy bits
	p.Task("buildall", do.P{"build_darwin", "build_windows"}, nil)

	p.Task("build_windows", do.P{"build_win32", "build_win64"}, nil)

	p.Task("build_darwin", nil, func(c *do.Context) {
		c.Run("GOOS=darwin GOARCH=amd64 gb build")
	}).Src("src/**/*.go")

	p.Task("build_win32", nil, func(c *do.Context) {
		c.Run("GOOS=windows GOARCH=386 gb build")
	}).Src("src/**/*.go")

	p.Task("build_win64", nil, func(c *do.Context) {
		c.Run("GOOS=windows GOARCH=amd64 gb build")
	}).Src("src/**/*.go")

	// Vendor Updates
	p.Task("blind_update", nil, func(c *do.Context) {
		c.Run("gb vendor update --all")
	})
}

func main() {
	do.Godo(tasks)
}
