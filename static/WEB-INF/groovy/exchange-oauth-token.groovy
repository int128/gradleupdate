import groovy.json.JsonBuilder
import infrastructure.GitHub
import util.CrossOriginPolicy

CrossOriginPolicy.allowOrigin(response, headers)

assert params.code, 'code parameter should be given'

final exchanged = GitHub.exchangeOAuthToken(params.code)

assert exchanged
assert !exchanged.error
assert exchanged.access_token
assert exchanged.scope

response.contentType = 'application/json'
println new JsonBuilder({
    token exchanged.access_token
    scope exchanged.scope
})
