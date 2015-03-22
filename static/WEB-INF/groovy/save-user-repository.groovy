import groovy.json.JsonBuilder
import groovy.json.JsonSlurper
import infrastructure.GitHub
import model.GitHubRepository
import util.CrossOriginPolicy

CrossOriginPolicy.allowOrigin(response, headers)

assert params.fullName
assert headers.Authorization

final json = new JsonSlurper().parse(request.inputStream)
assert json.autoUpdate instanceof Boolean

final github = new GitHub(Authorization: headers.Authorization)
final repo = github.getRepository(params.fullName)

if (!repo.permissions.admin) {
    response.sendError 404, 'No permission'
}

final metadata = new GitHubRepository(fullName: params.fullName, autoUpdate: json.autoUpdate)
metadata.save()

response.contentType = 'application/json'

println new JsonBuilder({
    autoUpdate metadata?.autoUpdate ?: false
})
