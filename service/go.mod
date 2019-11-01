module dbutil

go 1.13

require (
	github.com/go-sql-driver/mysql v1.4.1
	sparrow/sparrow v0.0.0-00010101000000-000000000000
)

replace sparrow/sparrow => ../../../go/src/gitlab.geax.io/sparrow/sparrow
