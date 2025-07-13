module tests

go 1.22

require (
	github.com/aws/aws-sdk-go-v2 v1.36.5
	github.com/aws/aws-sdk-go-v2/config v1.29.17
	github.com/aws/aws-sdk-go-v2/service/s3 v1.83.0
)

// 本地依赖app-server包
replace app-server => ../app-server
require app-server v0.0.0-00010101000000-000000000000
