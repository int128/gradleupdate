import groovy.json.JsonSlurper
import infrastructure.GitHubWebhook

final eventType = headers.'X-GitHub-Event'
final delivery = headers.'X-GitHub-Delivery'
final signature = headers.'X-Hub-Signature'
final payload = request.inputStream.bytes
assert eventType
assert delivery
assert signature
assert payload

assert new GitHubWebhook().validate(signature, payload)

log.info("Delivery: $delivery")
log.info("Event type: $eventType")

final json = new JsonSlurper().parse(payload)
assert json instanceof Map

if (eventType == 'watch' && json.action == 'started') {
    assert json.sender.login
    log.info("Queue updating repositories of stargazer ${json.sender.login}")
    defaultQueue.add(
            url: '/internal/api/got-star/0.groovy',
            params: [stargazer: json.sender.login])

} else if (eventType == 'ping') {
    log.info("Ping: ${json.zen}")

} else {
    log.warning('Unknown event')
    log.info("Payload: ${new String(payload)}")
    response.sendError(404)
}
