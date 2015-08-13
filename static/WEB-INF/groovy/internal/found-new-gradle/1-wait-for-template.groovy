import gradle.TemplateRepository

import static util.RequestUtil.relativePath

final gradleVersion = params.gradle_version
assert gradleVersion instanceof String

final templateRepository = new TemplateRepository()
final current = templateRepository.fetchGradleWrapperVersion('master')

if (current == gradleVersion) {
    log.info("The template is up-to-date $current, so queue updating stargazers")
    defaultQueue.add(
            url: relativePath(request, '2-stargazers.groovy'),
            params: [gradle_version: gradleVersion])
} else {
    log.info("The template is still old $current while expected $gradleVersion, retrying")
    response.sendError(503)
}
