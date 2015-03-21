import groovy.json.JsonBuilder
import infrastructure.GitHub
import util.CrossOrigin

CrossOrigin.sendAccessControlAllowOrigin(response, headers)

assert params.code, 'code parameter should be given'

final exchanged = GitHub.exchangeOAuthToken(params.code)

assert exchanged
assert !exchanged.error
assert exchanged.access_token
assert exchanged.scope

response.contentType = 'application/json'

def json = new JsonBuilder()
json {
    token exchanged.access_token
    scope exchanged.scope
}
println json
