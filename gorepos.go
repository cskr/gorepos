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
	"github.com/howeyc/fsnotify"
	"io"
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
		fmt.Fprintln(os.Stderr, "Reading package list failed:", err)
		os.Exit(2)
	}

	fmt.Printf("Serving package(s) on %s...\n", *addr)
	err = http.ListenAndServe(*addr, pl)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Server failed to start:", err)
		os.Exit(3)
	}
}

type PackageList struct {
	packages map[string]Package
	mx       sync.RWMutex
	file     string
}

func NewPackageList(pkgFile string) (pl *PackageList, err error) {
	pl = &PackageList{file: pkgFile}

	err = pl.loadPackages()
	if err != nil {
		return nil, err
	}

	go func() {
		err := pl.watch()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Watching package list failed: ", err)
		}
	}()

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

	pkgs := make(map[string]Package)
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

func (pl *PackageList) watch() error {
	for {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			return err
		}

		err = watcher.Watch(pl.file)
		if err != nil {
			return err
		}

		select {
		case ev := <-watcher.Event:
			if ev.IsModify() {
				pl.loadPackages()
			}

		case err = <-watcher.Error:
			return err
		}

		err = watcher.Close()
		if err != nil {
			return err
		}
	}

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
		if pkg, ok := pl.packages[r.URL.Path]; ok {
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
