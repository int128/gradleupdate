import gradle.TemplateRepository
import infrastructure.GitHub

assert params.gradleVersion
assert params.next

final gitHub = new GitHub()
final templateRepository = new TemplateRepository(gitHub)
final current = templateRepository.queryGradleWrapperVersion()

log.info("Current Gradle version of the template: $current")
log.info("Expected Gradle version: ${params.gradleVersion}")

if (current == params.gradleVersion) {
    log.info("Queue: ${params.next}")
    defaultQueue.add(url: params.next)
} else {
    response.sendError(503, "Gradle version expected ${params.gradleVersion} but current $current")
}
