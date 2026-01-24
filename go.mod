module seaotterms-api

go 1.24.0

require (
	github.com/gofiber/fiber/v2 v2.52.10
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
	github.com/sirupsen/logrus v1.9.4
	golang.org/x/crypto v0.47.0
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	gorm.io/driver/postgres v1.6.0
	gorm.io/gorm v1.31.1
	seaotterms-db v0.0.0
)

replace seaotterms-db => ../seaotterms-db

require (
	github.com/andybalholm/brotli v1.2.0 // indirect
	github.com/clipperhouse/stringish v0.1.1 // indirect
	github.com/clipperhouse/uax29/v2 v2.3.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.8.0 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/compress v1.18.3 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.19 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.69.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
)
