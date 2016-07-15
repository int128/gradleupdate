import domain.GradleUpdate
import util.RequestUtil

import static entity.PullRequestForLatestGradleWrapperTaskState.State

assert params.full_name

final ownerToken = RequestUtil.token(headers)
assert ownerToken

def pullRequestForLatestGradleWrapper = new GradleUpdate().pullRequestForLatestGradleWrapper(params.full_name)
pullRequestForLatestGradleWrapper.checkOwnership(ownerToken)

datastore.withTransaction {
    pullRequestForLatestGradleWrapper.reportProgress(State.Queued)
    defaultQueue.add(
            url: '/api/create-pull-request-for-update-gradle.task.groovy',
            params: params)
}
