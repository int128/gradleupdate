import domain.GradleUpdate

import static util.RequestUtil.relativePath

assert params.gradle_version

def gradleUpdate = new GradleUpdate()

log.info("Requesting bump Gradle version of the template")
gradleUpdate.gradleWrapperTemplateRepository.updateAsync(params.gradle_version)

log.info("Queue updating the template repository to $params.gradle_version")
defaultQueue.add(
        url: relativePath(request, '1-wait-for-template.task.groovy'),
        params: [gradle_version: params.gradle_version],
        countdownMillis: 3 * 60 * 1000)
