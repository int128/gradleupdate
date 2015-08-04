import gradle.TemplateRepository
import infrastructure.GitHub

import static util.RequestUtil.relativePath

final fullName = params.full_name
assert fullName instanceof String

final forkName = params.fork_name
assert forkName instanceof String

final intoBranch = params.into_branch
assert intoBranch instanceof String

final forkOwner = params.fork_owner
assert forkOwner instanceof String

final gitHub = new GitHub()
final templateRepository = new TemplateRepository(gitHub)

log.info("Creating a tree on $forkName")
final tree = templateRepository.createTreeWithGradleWrapper(forkName)
final gradleVersion = templateRepository.queryGradleWrapperVersion()
final gradleBranch = "gradle-$gradleVersion"

log.info("Creating a branch $gradleBranch on $forkName")
gitHub.createBranch(forkName, gradleBranch, intoBranch, "Gradle $gradleVersion", tree)

final from = "$forkOwner:$gradleBranch"
log.info("Queue sending a pull request from $from into $fullName:$intoBranch")
defaultQueue.add(
        url: relativePath(request, '5-pull-request.groovy'),
        params: [
                full_name: fullName,
                into_branch: intoBranch,
                from: from,
                gradle_version: gradleVersion
        ],
        countdownMillis: 1000)
