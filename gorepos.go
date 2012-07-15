/*
 * Copyright (C) 2012 Chandra Sekar S
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
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
		fmt.Fprintln(os.Stderr, "Reading package list failed:", err)
		os.Exit(2)
	}

	fmt.Printf("Serving %d package(s) on %s...\n", len(pl.Packages), *addr)
	err = http.ListenAndServe(*addr, pl)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Server failed to start:", err)
		os.Exit(3)
	}
}

type PackageList struct {
	Packages map[string]Package
}

func NewPackageList(pkgFile string) (pl *PackageList, err error) {
	f, err := os.Open(pkgFile)
	if err != nil {
		return
	}
	in := bufio.NewReader(f)

	pkgs := make(map[string]Package)
	for {
		ln, _, err := in.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		pkg := NewPackage(string(ln))
		pkgs[pkg.Path] = pkg
	}
	return &PackageList{pkgs}, nil
}

func (pl *PackageList) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		indexTmpl.Execute(w, map[string]interface{}{
			"host": r.Host,
			"pkgs": pl.Packages,
		})
	} else {
		if pkg, ok := pl.Packages[r.URL.Path]; ok {
			pkgTmpl.Execute(w, map[string]interface{}{
				"host": r.Host,
				"pkg":  pkg,
			})
		} else {
			http.NotFound(w, r)
		}
	}
}

type Package struct {
	Path, Vcs, Repo string
}

func NewPackage(line string) Package {
	fields := strings.Fields(line)
	return Package{fields[0], fields[1], fields[2]}
}
