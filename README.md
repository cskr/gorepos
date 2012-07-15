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
