import groovy.json.JsonBuilder
import infrastructure.GitHubOAuth

assert params.code, 'code parameter should be given'

final oauth = new GitHubOAuth()
final exchanged = oauth.exchangeCodeAndToken(params.code)

assert exchanged
assert !exchanged.error
assert exchanged.access_token
assert exchanged.scope

response.contentType = 'application/json'
println new JsonBuilder({
    token exchanged.access_token
    scope exchanged.scope
})
