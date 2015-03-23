import groovy.json.JsonBuilder
import model.GitHubRepository

final entities = datastore.execute {
    select all from 'GitHubRepository'
    where pullRequestOnStableRelease == true
}.collect {
    it as GitHubRepository
}

response.contentType = 'application/json'
println new JsonBuilder(entities)
