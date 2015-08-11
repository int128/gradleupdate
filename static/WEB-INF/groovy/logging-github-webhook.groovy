import infrastructure.GitHubWebhook

final eventType = headers.'X-GitHub-Event'
final signature = headers.'X-Hub-Signature'
assert eventType instanceof String
assert signature instanceof String

final payload = request.inputStream.bytes
assert payload

log.info("Event type is $eventType")
log.info("Payload length is ${payload.length}")
log.info("Signature is $signature")

assert new GitHubWebhook().validate(signature, payload)

log.info('Signature is valid')
