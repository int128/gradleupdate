import gradle.Repository

assert params.full_name

final repository = new Repository(params.full_name)
final state = repository.checkIfGradleWrapperIsLatest(params.branch ?: 'master')
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
