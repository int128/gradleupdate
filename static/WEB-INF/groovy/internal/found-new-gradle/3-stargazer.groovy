import gradle.Repository

final gradleVersion = params.gradle_version
assert gradleVersion instanceof String

final stargazer = params.stargazer
assert stargazer instanceof String

Repository.fetchRepositories(stargazer).each { repo ->
    log.info("Queue updating the repository $repo.full_name")
    defaultQueue.add(
            url: '/internal/pull-request-for-gradle/0.groovy',
            params: [
                    full_name: repo.full_name,
                    branch: repo.default_branch,
                    gradle_version: gradleVersion,
            ])
}
