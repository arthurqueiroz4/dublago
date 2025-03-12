package singleton

import (
	"tradutor-dos-crias/pipeline"
	"tradutor-dos-crias/user"
)

var (
	UserService = &user.UserService{}
	Pipeline    = &pipeline.Pipeline{}
)
