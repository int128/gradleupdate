import gradle.TemplateRepository

import static util.RequestUtil.relativePath

assert params.gradle_version

final templateRepository = new TemplateRepository()

log.info("Requesting bump Gradle version of the template")
templateRepository.bumpVersion(params.gradle_version)

log.info("Queue updating the template repository to $params.gradle_version")
defaultQueue.add(
        url: relativePath(request, '1-wait-for-template.groovy'),
        params: [gradle_version: params.gradle_version],
        countdownMillis: 3 * 60 * 1000)
