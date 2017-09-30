package impl

import "psychic-rat/api"
import "psychic-rat/repo"
import irepo "psychic-rat/impl/repo"
import iapi "psychic-rat/impl/api"

func init() {
	repos = irepo.GetRepos()
	apis = iapi.GetApis(repos)
}

var apis api.Api
var repos repo.Repos

func GetApi() api.Api {
	return apis
}

func GetRepos() repo.Repos {
	return repos
}
