import gradle.TemplateRepository
import infrastructure.GitHub

import static util.RequestUtil.relativePath

final gradleVersion = params.gradle_version
assert gradleVersion instanceof String

final gitHub = new GitHub()
final templateRepository = new TemplateRepository(gitHub)

log.info("Requesting bump Gradle version of the template")
templateRepository.bumpTemplate(gradleVersion)

log.info("Queue updating the template repository to $gradleVersion")
defaultQueue.add(
        url: relativePath(request, '1-wait-for-template.groovy'),
        params: [gradle_version: gradleVersion],
        countdownMillis: 3 * 60 * 1000)
