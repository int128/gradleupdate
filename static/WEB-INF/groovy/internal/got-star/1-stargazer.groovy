import gradle.Repositories

assert params.gradle_version
assert params.stargazer

final repositories = new Repositories(params.stargazer)

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

assert page.current instanceof List
page.current.findAll { repo ->
    if (repo.fork) {
        log.info("Repository $repo.full_name is a fork, skip")
        false
    }
    true
}.each { repo ->
    log.info("Queue updating the repository $repo.full_name")
    defaultQueue.add(
            url: '/internal/pull-request-for-gradle/0.groovy',
            params: [
                    full_name: repo.full_name,
                    branch: repo.default_branch,
                    gradle_version: params.gradle_version,
            ])
}
