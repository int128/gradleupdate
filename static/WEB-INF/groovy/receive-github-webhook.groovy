import groovy.json.JsonSlurper

final eventType = headers.'X-GitHub-Event'

switch (eventType) {
    case 'watch':
        final json = new JsonSlurper().parse(request.inputStream)
        assert json instanceof Map
        assert json.action == 'started'
        assert json.sender.login
        log.info("Queue updating repositories of user: ${json.sender.login}")
        defaultQueue.add(
                url: '/internal/update-gradle-of-user.groovy',
                params: [user: json.sender.login])
        break

    default:
        log.warning("Received unknown event type: $eventType")
        response.sendError(404, 'Event not handled')
        break
}
