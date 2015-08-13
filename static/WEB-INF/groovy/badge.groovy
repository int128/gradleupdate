import gradle.Repository
import infrastructure.GitHub

final fullName = params.full_name
assert fullName instanceof String

final branch = params.branch
assert branch instanceof String

final gitHub = new GitHub()
final repository = new Repository(fullName, gitHub)

log.info("Checking if repository $fullName is latest")
final state = repository.checkIfGradleWrapperIsLatest(branch)
switch (state) {
    case Repository.GradleWrapperState.UP_TO_DATE:
        forward("/util/svg-badge.groovy?fill=#4c1&message=${state.currentVersion}")
        break
    case Repository.GradleWrapperState.OUT_OF_DATE:
        forward("/util/svg-badge.groovy?fill=#e05d44&message=${state.currentVersion}")
        break
    default:
        forward("/util/svg-badge.groovy?fill=#9f9f9f&message=unknown")
        break
}
