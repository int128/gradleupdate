import domain.GradleUpdate

assert params.full_name

new GradleUpdate().pullRequestForLatestGradleWrapper(params.full_name).create()
