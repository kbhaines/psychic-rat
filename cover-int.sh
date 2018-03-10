go test -coverprofile=prof.out -coverpkg psychic-rat/web,psychic-rat/web/admin,psychic-rat/web/pub,psychic-rat/sess,psychic-rat/sqldb psychic-rat/tests
go tool cover -html=prof.out
