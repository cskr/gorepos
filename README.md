gorepos allows you to,

* `go get` private packages (even over SSH).
* Change your code hosting location without changing imports.
* Change VCS without changing imports.
* Provide import paths on your domain without managing the VCS yourself.

For example, you can use gitolite/gitorious to manage all your package repositories and provide authenticated access to them without losing the convenience of `go get`.

While `go get` has limitations in accessing private repos directly, it provides a nice discovery mechanism through `<meta>` tags which gorepos uses to provide convenient access.

# Configuration

Ensure $GOPATH/bin is part of your $PATH. Install gorepos with,

	go get github.com/tuxychandru/gorepos

The package list file must contain a line for each available package in the form,

	<path> <vcs> <repo-root>

# Example

To provide access to a private git bitbucket repository on the import path `go.mycompany.com/mylib`, put this line in a package list file named `pkgs`,

	/mylib git ssh://git@bitbucket.org/mycompany/mylib

Start gorepos on the server at `go.mycompany.com` with,

	sudo gorepos -a go.mycompany.com:80 -p pkgs

Now you can `import "go.mycompany.com/mylib"` in your program and `go get` it.

## License

Copyright (C) 2012 Chandra Sekar S

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
