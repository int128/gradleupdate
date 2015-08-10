import groovy.json.JsonSlurper

final eventType = headers.'X-GitHub-Event'

switch (eventType) {
    case 'watch':
        final json = new JsonSlurper().parse(request.inputStream)
        assert json instanceof Map
        assert json.action == 'started'
        assert json.sender.login
        log.info("Queue updating repositories of stargazer ${json.sender.login}")
        defaultQueue.add(
                url: '/internal/got-star/0.groovy',
                params: [stargazer: json.sender.login])
        break

    default:
        log.warning("Received unknown event type: $eventType")
        response.sendError(404, 'Event not handled')
        break
}
