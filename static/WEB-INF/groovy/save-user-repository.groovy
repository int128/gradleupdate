import groovy.json.JsonBuilder
import groovy.json.JsonSlurper
import infrastructure.GitHub
import model.GitHubRepository
import util.CrossOriginPolicy

import static com.google.appengine.api.utils.SystemProperty.Environment.Value.Development

CrossOriginPolicy.allowOrigin(response, headers)

assert params.fullName
assert app.env.name == Development || headers.Authorization

final json = new JsonSlurper().parse(request.inputStream)
assert json.autoUpdate instanceof Boolean

final gitHub = GitHub.authorizationOrDefault(headers.Authorization)

final repo = gitHub.getRepository(params.fullName)

if (!repo.permissions.admin) {
    response.sendError 404, 'No permission'
}

final metadata = new GitHubRepository(fullName: params.fullName, autoUpdate: json.autoUpdate)
metadata.save()

response.contentType = 'application/json'

println new JsonBuilder({
    autoUpdate metadata?.autoUpdate ?: false
})
