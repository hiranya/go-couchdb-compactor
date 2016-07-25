# go-couchdb-compactor
A database compactor for CouchDB servers

## Usage
go run compactor.go -s http://localhost:6901 -u admin -p p@ssw0rd

```
./go-couchdb-compactor-mac-386 --help
Usage of ./go-couchdb-compactor-mac-386:
  -c int
    	Concurrency level required for compaction (default 5)
  -p string
    	Password to access the CouchDB server
  -s string
    	CouchDB server url. Defaults to http://localhost:5984 (default "http://localhost:5984")
  -u string
    	Username to access the CouchDB server
```

## Notes
Check ./binaries directory for pre-built binaries for Linux and Mac
