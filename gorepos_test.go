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
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestIndex(t *testing.T) {
	list, err := generateList()
	if err != nil {
		t.Errorf("Error generating package list: %s", err)
	}

	pl, err := NewPackageList(list)
	if err != nil {
		t.Errorf("Error reading package list: %s", err)
		return
	}

	r, _ := http.NewRequest("GET", "http://example.com/", nil)
	w := httptest.NewRecorder()
	pl.ServeHTTP(w, r)

	b := new(bytes.Buffer)
	indexTmpl.Execute(b, map[string]interface{}{
		"host": "example.com",
		"pkgs": [...]map[string]string{
			{
				"Path": "/lib1",
			},
			{
				"Path": "/lib2",
			},
			{
				"Path": "/lib3",
			},
		},
	})

	if w.Body.String() != b.String() {
		t.Errorf("Body = %s, want %s", w.Body.String(), b.String())
	}
}

func TestPkg(t *testing.T) {
	list, err := generateList()
	if err != nil {
		t.Errorf("Error generating package list: %s", err)
	}

	pl, err := NewPackageList(list)
	if err != nil {
		t.Errorf("Error reading package list: %s", err)
		return
	}

	body, expected := invokePkg(pl, "lib1", "git", "ssh://git@bitbucket.org/user1/lib1")
	if body != expected {
		t.Errorf("Body = %s, want %s", body, expected)
		return
	}

	body, expected = invokePkg(pl, "lib2", "hg", "ssh://hg@bitbucket.org/user2/lib2")
	if body != expected {
		t.Errorf("Body = %s, want %s", body, expected)
		return
	}

	body, expected = invokePkg(pl, "lib3", "git", "ssh://git@go.mydomain.com/lib3")
	if body != expected {
		t.Errorf("Body = %s, want %s", body, expected)
	}
}

func TestReload(t *testing.T) {
	list, err := generateList()
	if err != nil {
		t.Errorf("Error generating package list: %s", err)
	}

	pl, err := NewPackageList(list)
	if err != nil {
		t.Errorf("Error reading package list: %s", err)
		return
	}

	time.Sleep(100 * time.Millisecond)
	err = appendList(list, "/lib4 git ssh://git@go.mydomain.com/lib4")
	if err != nil {
		t.Errorf("Error appending item: %s", err)
		return
	}

	time.Sleep(100 * time.Millisecond)
	body, expected := invokePkg(pl, "lib4", "git", "ssh://git@go.mydomain.com/lib4")
	if body != expected {
		t.Errorf("Body = %s, want %s", body, expected)
	}
}

func generateList() (fname string, err error) {
	fname = fmt.Sprintf("%s%ctest_list", os.TempDir(), os.PathSeparator)
	f, err := os.Create(fname)
	if err != nil {
		return "", err
	}
	defer f.Close()

	fmt.Fprintln(f, "/lib1 git ssh://git@bitbucket.org/user1/lib1")
	fmt.Fprintln(f, "/lib2 hg ssh://hg@bitbucket.org/user2/lib2")
	fmt.Fprintln(f, "/lib3 git ssh://git@go.mydomain.com/lib3")
	return fname, nil
}

func appendList(list, line string) error {
	f, err := os.OpenFile(list, os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintln(f, line)
	return nil
}

func invokePkg(pl *PackageList, pkg, vcs, repo string) (body, expected string) {
	r, _ := http.NewRequest("GET", "http://example.com/"+pkg, nil)
	w := httptest.NewRecorder()
	pl.ServeHTTP(w, r)

	b := new(bytes.Buffer)
	pkgTmpl.Execute(b, map[string]interface{}{
		"host": "example.com",
		"pkg": map[string]string{
			"Path": "/" + pkg,
			"Vcs":  vcs,
			"Repo": repo,
		},
	})

	return w.Body.String(), b.String()
}
