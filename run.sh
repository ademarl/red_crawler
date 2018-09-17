rm bin/red_crawler 2>/dev/null
go build -o bin/red_crawler src/basic_crawler.go src/paper.go src/parsing.go src/share_database.go src/red_crawler.go
./bin/red_crawler
