import groovy.json.JsonBuilder
import infrastructure.GitHub
import service.GitHubRepositoryService
import util.CrossOriginPolicy

CrossOriginPolicy.allowOrigin(response, headers)

assert params.fullName
assert headers.Authorization
assert headers.Authorization.startsWith('token ')

final gitHub = new GitHub(Authorization: headers.Authorization)
final service = new GitHubRepositoryService(gitHub)
final metadata = service.queryMetadata(params.fullName)

if (metadata == null) {
    response.sendError(404, 'No Permission')
}

response.contentType = 'application/json'
println new JsonBuilder(metadata)
