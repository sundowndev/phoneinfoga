/*
Package dorkgen is a Go package to generate dork requests for popular search engines such as Google, DuckDuckGo and Bing.
It allows you to define requests programmatically and convert them into string.
You can use it as following:

package main

import "github.com/sundowndev/dorkgen"

func main() {
	dork := &dorkgen.GoogleSearch{}
	// dork := &dorkgen.DuckDuckGo{}
	// dork := &dorkgen.Bing{}

	dork.Site("example.com").Intext("text").ToString()
	// returns: site:example.com "text"
}
*/
package dorkgen
