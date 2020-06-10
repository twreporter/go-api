module twreporter.org/go-api

go 1.12

require (
	cloud.google.com/go v0.52.0 // indirect
	github.com/algolia/algoliasearch-client-go v0.0.0-20181217121925-1c0b06b9e47a
	github.com/auth0/go-jwt-middleware v0.0.0-20170425171159-5493cabe49f7
	github.com/aws/aws-sdk-go v1.18.2
	github.com/cheekybits/is v0.0.0-20150225183255-68e9c0620927 // indirect
	github.com/codegangsta/negroni v1.0.0 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/erikstmartin/go-testdb v0.0.0-20160219214506-8d10e4a1bae5 // indirect
	github.com/gin-contrib/cors v0.0.0-20170708080947-567de1916927
	github.com/gin-contrib/sessions v0.0.0-20180827025425-58cbcf30135c
	github.com/gin-gonic/gin v1.5.0
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/go-sql-driver/mysql v1.4.1
	github.com/gofrs/uuid v3.2.0+incompatible // indirect
	github.com/golang-migrate/migrate/v4 v4.6.1
	github.com/google/uuid v1.1.1
	github.com/gorilla/mux v1.7.2 // indirect
	github.com/gorilla/sessions v1.1.1 // indirect
	github.com/jinzhu/copier v0.0.0-20180308034124-7e38e58719c3
	github.com/jinzhu/gorm v1.9.2
	github.com/jinzhu/inflection v0.0.0-20170102125226-1c35d901db3d // indirect
	github.com/jinzhu/now v1.0.1 // indirect
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/kidstuff/mongostore v0.0.0-20180412085134-db2a8b4fac1f // indirect
	github.com/kr/pretty v0.2.0 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/lib/pq v1.1.1 // indirect
	github.com/matryer/try v0.0.0-20161228173917-9ac251b645a2 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.4.2
	github.com/smartystreets/goconvey v0.0.0-20190330032615-68dc04aab96a // indirect
	github.com/spf13/viper v1.3.2
	github.com/stretchr/testify v1.4.0
	github.com/twreporter/go-api v4.0.0+incompatible
	github.com/twreporter/logformatter v0.0.0-20200211094126-60fe42618206
	go.mongodb.org/mongo-driver v1.1.0
	golang.org/x/crypto v0.0.0-20200210222208-86ce3cb69678
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2 // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sys v0.0.0-20200202164722-d101bd2416d5 // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	google.golang.org/genproto v0.0.0-20200211035748-55294c81d784 // indirect
	google.golang.org/grpc v1.27.1 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/go-playground/validator.v8 v8.18.2
	gopkg.in/go-playground/validator.v9 v9.31.0 // indirect
	gopkg.in/guregu/null.v3 v3.4.0
	gopkg.in/matryer/try.v1 v1.0.0-20150601225556-312d2599e12e
	gopkg.in/mgo.v2 v2.0.0-20160818020120-3f83fa500528
	gopkg.in/yaml.v2 v2.2.8 // indirect
)

replace (
	github.com/coreos/bbolt v1.3.4 => go.etcd.io/bbolt v1.3.4
	go.etcd.io/bbolt v1.3.4 => github.com/coreos/bbolt v1.3.4
)
