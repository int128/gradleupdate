import groovy.json.JsonSlurper

final eventType = headers.'X-GitHub-Event'

switch (eventType) {
    case 'WatchEvent':
        final json = new JsonSlurper().parse(request.inputStream)
        assert json instanceof Map
        assert json.action == 'started'
        // TODO: handle event
        log.info json.toString()
        break

    default:
        log.warning("Received unknown event type: $eventType")
        response.sendError(404, 'Event not handled')
        break
}
