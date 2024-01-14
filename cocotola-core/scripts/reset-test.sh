migrate -database 'postgres://user:password@127.0.0.1:5432/postgres?sslmode=disable' -source file://src/sqls/postgres/ drop
migrate -database 'postgres://user:password@127.0.0.1:5432/postgres?sslmode=disable' -source file://src/sqls/postgres/ up

migrate -database 'mysql://user:password@tcp(127.0.0.1:3306)/development' -source file://src/sqls/mysql/ drop
migrate -database 'mysql://user:password@tcp(127.0.0.1:3306)/development' -source file://src/sqls/mysql/ up