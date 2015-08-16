import gradle.Repository

final fullName = params.full_name
final branch = params.branch ?: 'master'
assert fullName instanceof String
assert branch instanceof String

final repository = new Repository(fullName)
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
