module github.com/aruruka/go-google-s3-uploader/auth-server

go 1.24.2

require (
	github.com/aruruka/go-google-s3-uploader/shared v0.0.0-00010101000000-000000000000
	github.com/coreos/go-oidc/v3 v3.9.0
	golang.org/x/oauth2 v0.15.0
)

replace github.com/aruruka/go-google-s3-uploader/shared => ../shared

require (
	cloud.google.com/go/compute v1.20.1 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	github.com/go-jose/go-jose/v3 v3.0.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
)
