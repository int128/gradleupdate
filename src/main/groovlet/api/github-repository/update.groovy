import infrastructure.GitHub
import model.Credential

assert params.full_name
assert params.gradle_version

assert headers.Authorization?.startsWith('token ')
final ownerToken = headers.Authorization.substring('token '.length())

final gitHubByOwner = new GitHub(new Credential(secret: ownerToken))
final repositoryByOwner = gitHubByOwner.fetch(params.full_name)
if (repositoryByOwner == null) {
    response.sendError(404, "Owner does not have permission to access $params.full_name")
    return
}
assert repositoryByOwner.permissions
if (!repositoryByOwner.permissions.push) {
    response.sendError(403, "Owner does not have permission to push $params.full_name")
    return
}

final gitHubByGradleUpdate = new GitHub()
final repositoryByGradleUpdate = gitHubByGradleUpdate.fetch(params.full_name)
if (repositoryByGradleUpdate == null) {
    response.sendError(404, "Gradle Update does not have permission to access $params.full_name")
    return
}
assert repositoryByGradleUpdate.permissions
if (!repositoryByGradleUpdate.permissions.pull) {
    response.sendError(403, "Gradle Update does not have permission to pull $params.full_name")
    return
}

log.info("Queue updating the repository $repositoryByGradleUpdate.full_name")
defaultQueue.add(
        url: '/internal/api/pull-request-for-gradle/0.groovy',
        params: [
                full_name: repositoryByGradleUpdate.full_name,
                branch: repositoryByGradleUpdate.default_branch,
                gradle_version: params.gradle_version,
        ])
