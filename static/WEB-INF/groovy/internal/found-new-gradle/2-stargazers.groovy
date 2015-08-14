import gradle.Stargazers

import static util.RequestUtil.relativePath

final gradleVersion = params.gradle_version
assert gradleVersion instanceof String

final stargazers = new Stargazers()

final page
final next = params.next
if (next instanceof String) {
    page = stargazers.fetchNext(next)
} else {
    page = stargazers.fetchFirst()
}

final nextPage = page.rel.next
if (nextPage) {
    log.info("Queue next page of $stargazers: $nextPage")
    defaultQueue.add(
            url: request.requestURI,
            params: params + [next: nextPage])
} else {
    log.info("Now last page of $stargazers")
}

page.current.each { stargazer ->
    log.info("Queue updating stargazer ${stargazer.login}")
    defaultQueue.add(
            url: relativePath(request, '3-stargazer.groovy'),
            params: [stargazer: stargazer.login, gradle_version: gradleVersion])
}
