import gradle.TemplateRepository

import static util.RequestUtil.relativePath

assert params.gradle_version

final templateRepository = new TemplateRepository()
final current = templateRepository.fetchGradleWrapperVersion('master')

if (current == params.gradle_version) {
    log.info("The template is up-to-date $current, so queue updating stargazers")
    defaultQueue.add(
            url: relativePath(request, '2-stargazers.groovy'),
            params: [gradle_version: params.gradle_version])
} else {
    log.info("The template is still old $current while expected $params.gradle_version, retrying")
    response.sendError(503)
}
