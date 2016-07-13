import model.PullRequest

assert params.url
assert params.createdAt
assert params.ownerId

final pullRequest = new PullRequest(params + [
        createdAt: Date.parse("yyyy-MM-dd'T'HH:mm:ss", params.createdAt),
        ownerId: params.ownerId as int,
])

log.info("Saving $pullRequest")
pullRequest.save()
