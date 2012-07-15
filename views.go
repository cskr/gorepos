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

import "html/template"

var (
	indexTmpl = template.Must(template.New("index").Parse(indexView))
	pkgTmpl   = template.Must(template.New("pkg").Parse(pkgView))
)

const indexView = `
<html>
	<head>
		<title>gorepos - Packages</title>
	</head>
	<body>
		<h1>Available Packages</h1>
		<ul>
			{{$host := .host}}
			{{range .pkgs}}
				<li><a href="{{.Path}}">{{$host}}{{.Path}}</a></li>
			{{end}}
		</ul>
	</body>
</html>
`

const pkgView = `
<html>
	<head>
		<meta name="go-import" content="{{.host}}{{.pkg.Path}} {{.pkg.Vcs}} {{.pkg.Repo}}">
		<title>gorepos - go.tuxychandru.com:9090/goaes</title>
	</head>
	<body>
		<h1>{{.host}}{{.pkg.Path}}</h1>
		<span style="font-weight: bold">VCS:</span>{{.pkg.Vcs}}git<br>
		<span style="font-weight: bold">Repo-Root:</span> {{.pkg.Repo}}
	</body>
</html>
`
