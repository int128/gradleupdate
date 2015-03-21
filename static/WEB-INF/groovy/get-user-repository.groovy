import groovy.json.JsonBuilder
import infrastructure.GitHub
import model.GitHubRepository
import util.CrossOrigin

CrossOrigin.sendAccessControlAllowOrigin(response, headers)

assert params.fullName
assert headers.Authorization

final github = new GitHub(Authorization: headers.Authorization)
final repo = github.getRepository(params.fullName)

if (!repo.permissions.admin) {
    response.sendError 404, 'No permission'
}

final metadata = GitHubRepository.get(params.fullName)

response.contentType = 'application/json'

println new JsonBuilder({
    autoUpdate metadata?.autoUpdate ?: false
})
