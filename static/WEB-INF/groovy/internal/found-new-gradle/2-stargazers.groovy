import gradle.Stargazers

import static util.RequestUtil.relativePath

final gradleVersion = params.gradle_version
assert gradleVersion instanceof String

final stargazers = new Stargazers().fetch()

stargazers.each { stargazer ->
    log.info("Queue updating stargazer ${stargazer.login}")
    defaultQueue.add(
            url: relativePath(request, '3-stargazer.groovy'),
            params: [stargazer: stargazer.login, gradle_version: gradleVersion])
}
