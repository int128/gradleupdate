import gradle.VersionWatcher

import static util.RequestUtil.relativePath

assert params.stargazer

log.info("Fetching version of the latest Gradle")
final gradleVersion = new VersionWatcher().fetchStableVersion()

log.info("Queue updating stargazer $params.stargazer")
defaultQueue.add(
        url: relativePath(request, '1-stargazer.groovy'),
        params: params + [gradle_version: gradleVersion])
