// Copyright 2013, Chandra Sekar S.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the README.md file.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

func main() {
	addr := flag.String("a", ":9090", "address to listen on (host:port)")
	pkgFile := flag.String("p", "", "package list")
	help := flag.Bool("help", false, "print usage")

	flag.Parse()
	if *help {
		flag.Usage()
		return
	}

	if *pkgFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	pl, err := NewPackageList(*pkgFile)
	if err != nil {
		log.Fatalln("Reading package list failed:", err)
	}

	fmt.Printf("Serving package(s) on %s...\n", *addr)
	err = http.ListenAndServe(*addr, pl)
	if err != nil {
		log.Fatalln("Server failed to start:", err)
	}
}

type PackageList struct {
	packages map[string]*Package
	mx       sync.RWMutex
	file     string
}

func NewPackageList(pkgFile string) (pl *PackageList, err error) {
	pl = &PackageList{file: pkgFile}

	err = pl.loadPackages()
	if err != nil {
		return nil, err
	}

	return pl, nil
}

func (pl *PackageList) loadPackages() error {
	pl.mx.Lock()
	defer pl.mx.Unlock()

	f, err := os.Open(pl.file)
	if err != nil {
		return err
	}
	defer f.Close()
	in := bufio.NewReader(f)

	pkgs := make(map[string]*Package)
	for {
		ln, _, err := in.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		line := string(ln)
		if len(strings.TrimSpace(line)) > 0 {
			pkg := NewPackage(line)
			pkgs[pkg.Path] = pkg
		}
	}
	pl.packages = pkgs

	return nil
}

func (pl *PackageList) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pl.mx.RLock()
	defer pl.mx.RUnlock()
	if r.URL.Path == "/" {
		indexTmpl.Execute(w, map[string]interface{}{
			"host": r.Host,
			"pkgs": pl.packages,
		})
	} else {
		if pkg, ok := pl.getPackage(r.URL.Path); ok {
			if r.FormValue("go-get") == "1" || pkg.Doc == "" {
				pkgTmpl.Execute(w, map[string]interface{}{
					"host": r.Host,
					"pkg":  pkg,
				})
			} else {
				http.Redirect(w, r, pkg.Doc, http.StatusFound)
			}
		} else {
			http.NotFound(w, r)
		}
	}
}

func (pl *PackageList) getPackage(path string) (*Package, bool) {
	if pkg, ok := pl.packages[path]; ok {
		return pkg, ok
	}

	for prefix := path; prefix != ""; prefix = prefix[:strings.LastIndex(prefix, "/")] {
		if pkg, ok := pl.packages[prefix]; ok {
			return pkg, ok
		}
	}

	return nil, false
}

type Package struct {
	Path, Vcs, Repo, Doc string
}

func NewPackage(line string) *Package {
	fields := strings.Fields(line)
	doc := ""
	if len(fields) > 3 {
		doc = fields[3]
	}

	return &Package{fields[0], fields[1], fields[2], doc}
}
