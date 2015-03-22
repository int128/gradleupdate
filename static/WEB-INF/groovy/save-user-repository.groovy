import groovy.json.JsonSlurper
import infrastructure.GitHub
import service.GitHubRepositoryService
import util.CrossOriginPolicy

import static com.google.appengine.api.utils.SystemProperty.Environment.Value.Development

CrossOriginPolicy.allowOrigin(response, headers)

assert params.fullName
assert app.env.name == Development || headers.Authorization

final json = new JsonSlurper().parse(request.inputStream)
assert json.autoUpdate instanceof Boolean

final gitHub = GitHub.authorizationOrDefault(headers.Authorization)
final service = new GitHubRepositoryService(gitHub)

if (service.saveMetadata(params.fullName, json.autoUpdate)) {
    response.status = 204
} else {
    response.sendError 404, 'No permission'
}
