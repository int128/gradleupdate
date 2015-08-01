import static util.RequestUtil.relativePath

final fullName = params.full_name
assert fullName instanceof String

log.info("Queue sending a pull request for the latest Gradle into ${fullName}")
defaultQueue.add(
        url: relativePath(request, '1-fork.groovy'),
        params: [full_name: fullName])
