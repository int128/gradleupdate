import gradle.VersionWatcher

import static util.RequestUtil.relativePath

final stargazer = params.stargazer
assert stargazer instanceof String

log.info("Fetching version of the latest Gradle")
final gradleVersion = new VersionWatcher().fetchStableVersion()

log.info("Queue updating stargazer $stargazer")
defaultQueue.add(
        url: relativePath(request, '1-stargazer.groovy'),
        params: params + [gradle_version: gradleVersion])
