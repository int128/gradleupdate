import groovy.json.JsonSlurper
import infrastructure.GitHubWebhook

final eventType = headers.'X-GitHub-Event'
final signature = headers.'X-Hub-Signature'
final payload = request.inputStream.bytes
assert eventType instanceof String
assert signature instanceof String
assert payload

assert new GitHubWebhook().validate(signature, payload)

final json = new JsonSlurper().parse(payload)
assert json instanceof Map

if (eventType == 'watch' && json.action == 'started') {
    assert json.sender.login
    log.info("Queue updating repositories of stargazer ${json.sender.login}")
    defaultQueue.add(
            url: '/internal/got-star/0.groovy',
            params: [stargazer: json.sender.login])
} else {
    log.warning("Unknown event type: $eventType")
    log.info("Payload: ${new String(payload)}")
    response.sendError(404)
}
