import infrastructure.GitHub
import util.CrossOriginPolicy

import static com.google.appengine.api.utils.SystemProperty.Environment.Value.Development

CrossOriginPolicy.allowOrigin(response, headers)

assert params.fullName
assert params.gradleVersion
assert app.env.name == Development || headers.Authorization

final gitHub = GitHub.authorizationOrDefault(headers.Authorization)
final repo = gitHub.getRepository(params.fullName)

if (repo.permissions?.admin) {
    log.info("Queue updating Gradle of the user repository: $params.fullName")
    defaultQueue.add(
            url: '/internal/update-gradle-of-user-repository.groovy',
            params: [fullName: params.fullName, gradleVersion: params.gradleVersion])
} else {
    response.sendError(404, 'No Admin Permission')
}
