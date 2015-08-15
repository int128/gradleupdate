import gradle.Repositories

final gradleVersion = params.gradle_version
assert gradleVersion instanceof String

final stargazer = params.stargazer
assert stargazer instanceof String

final repositories = new Repositories(stargazer)

final page
final next = params.next
if (next instanceof String) {
    page = repositories.fetchNext(next)
} else {
    page = repositories.fetchFirst()
}

final nextPage = page.rel.next
if (nextPage) {
    log.info("Queue next page of $repositories: $nextPage")
    defaultQueue.add(
            url: request.requestURI,
            params: params + [next: nextPage])
} else {
    log.info("Now last page of $repositories")
}

page.current.each { repo ->
    log.info("Queue updating the repository $repo.full_name")
    defaultQueue.add(
            url: '/internal/pull-request-for-gradle/0.groovy',
            params: [
                    full_name: repo.full_name,
                    branch: repo.default_branch,
                    gradle_version: gradleVersion,
            ])
}
