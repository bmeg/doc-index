

How to Develop
--------------

To install doc-index and work on the code, use the included builder. The commands
to setup a workspace are:

```
mkdir doc-index-ws
cd doc-index-ws
curl -O https://raw.githubusercontent.com/bmeg/doc-index/master/contrib/Makefile
make download
```

The source code and dependencies should now all be installed in this workspace.

To build utils
```
make
```

To run unit tests

```
make test
```