import groovy.json.JsonBuilder
import infrastructure.GitHubOAuth

assert params.code, 'code parameter should be given'
assert headers.origin, 'origin header should be given'

final exchanged = new GitHubOAuth().exchange(params.code)

assert exchanged
assert !exchanged.error
assert exchanged.access_token
assert exchanged.scope

response.headers.'Access-Control-Allow-Origin' =
    headers.origin.matches(/http:\/\/localhost(\:\d+)?/) ? headers.origin : 'https://gradleupdate.github.io'

response.contentType = 'application/json'
def json = new JsonBuilder()
json {
    token exchanged.access_token
    scope exchanged.scope
}
println json
