# LND in WebAssembly

This will be a workspace for cross-compiling `lnd`, the Lightning Labs lightning node implementation to WebAssembly, and ultimately to be run in a web browser.

This repository was initialized by setting `GOPATH` to this directory and running `go get -v -d github.com/lightningnetwork/lnd`

The `hide_git` script then renames all the recursive `.git` files in subdirectories to `.git2` so they won't be treated as submodules. 

Then arbitrary changes can be made to arbitrary dependencies with impunity. 

Rebasing and upgrading is left as an exercise to the reader. I have no idea how to do it. 
