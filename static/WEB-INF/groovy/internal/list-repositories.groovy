import groovy.json.JsonBuilder
import service.GitHubRepositoryService

final service = new GitHubRepositoryService()
final repositories = service.listPullRequestOnStableRelease()

response.contentType = 'application/json'
println new JsonBuilder(repositories)
