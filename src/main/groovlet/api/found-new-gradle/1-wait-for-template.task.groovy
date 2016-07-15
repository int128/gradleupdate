import domain.GradleUpdate

import static util.RequestUtil.relativePath

assert params.gradle_version

def gradleUpdate = new GradleUpdate()
assert gradleUpdate.gradleWrapperTemplateRepository.gradleWrapper.status.checkUpToDate()

log.info("The template is up-to-date $params.gradle_version, so queue updating stargazers")
defaultQueue.add(
        url: relativePath(request, '2-repositories.task.groovy'),
        params: [gradle_version: params.gradle_version])
