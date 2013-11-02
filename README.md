gorepos allows you to,

* `go get` private packages (even over SSH).
* Change your code hosting location without changing imports.
* Change VCS without changing imports.
* Provide import paths on your domain without managing the VCS yourself.
* Automatically re-load package list on changes.

For example, you can use gitolite/gitorious to manage all your package repositories and provide authenticated access to them without losing the convenience of `go get`.

While `go get` has limitations in accessing private repos directly, it provides a nice discovery mechanism through `<meta>` tags which gorepos uses to provide convenient access.

# Configuration

Ensure $GOPATH/bin is part of your $PATH. Install gorepos with,

	go get github.com/tuxychandru/gorepos

The package list file must contain a line for each available package in the form,

	<path> <vcs> <repo-root> [<godoc-url>]

# Example

To provide access to a private git bitbucket repository on the import path `go.mycompany.com/mylib`, put this line in a package list file named `pkgs`,

	/mylib git ssh://git@bitbucket.org/mycompany/mylib

Start gorepos on the server at `go.mycompany.com` with,

	sudo gorepos -a go.mycompany.com:80 -p pkgs

Now you can `import "go.mycompany.com/mylib"` in your program and `go get` it.

## License

Copyright (c) 2013, Chandra Sekar S  
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.
2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR
ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
