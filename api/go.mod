module github.com/arghanil/golang-redis-compose

go 1.14

require (
    github.com/swaggo/files v0.0.0-20190110041405-30649e0721f8 // indirect
    github.com/swaggo/http-swagger v0.0.0-20190324132102-654001218d89
    github.com/swaggo/swag v1.4.1
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-redis/redis/v8 v8.0.0-beta.7 // indirect
	github.com/gorilla/handlers v1.5.0
	github.com/gorilla/mux v1.8.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.6.0
	github.com/speps/go-hashids v2.0.0+incompatible
)

replace github.com/arghanil/golang-redis-compose/api/docs => ./docs
