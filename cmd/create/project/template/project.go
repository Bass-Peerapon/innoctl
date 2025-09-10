package template

import _ "embed"

//go:embed project/go.mod.tmpl
var GoMod []byte

//go:embed project/.dockerignore.tmpl
var DockerIgnore []byte

//go:embed project/.gitignore.tmpl
var GitIgnore []byte

//go:embed project/main.go.tmpl
var Main []byte

//go:embed project/.env.example.tmpl
var EnvExample []byte

//go:embed project/Dockerfile-development.tmpl
var DockerfileDevelopment []byte

//go:embed project/Dockerfile-production.tmpl
var DockerfileProduction []byte

//go:embed project/docker-compose.yml.tmpl
var DockerCompose []byte

//go:embed project/Makefile.tmpl
var Makefile []byte

//go:embed project/README.md.tmpl
var Readme []byte

//go:embed project/sonar-project.properties.tmpl
var SonarProperties []byte

//go:embed project/constants/constants.go.tmpl
var Constants []byte

//go:embed project/constants/enum.go.tmpl
var Enum []byte

//go:embed project/constants/error.go.tmpl
var Error []byte

//go:embed project/middleware/middleware.go.tmpl
var Middleware []byte

//go:embed project/middleware/jwt.go.tmpl
var JWT []byte

//go:embed project/middleware/openapi.go.tmpl
var OpenAPI []byte

//go:embed project/middleware/tracer.go.tmpl
var Tracer []byte

//go:embed project/utils/opentracing/init.go.tmpl
var OpenTracingInit []byte

//go:embed project/utils/redis/client.go.tmpl
var RedisClient []byte

//go:embed project/utils/pagination/pagination.go.tmpl
var Pagination []byte
