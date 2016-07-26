import infrastructure.GitHubOAuth

assert params.code
assert params.state
assert params.redirect_uri

final exchanged = GitHubOAuth.exchangeCodeAndToken(params.code, params.state, params.redirect_uri)
assert exchanged

response.setStatus(exchanged.statusCode)
response.contentType = exchanged.contentType
println exchanged.contentAsString
