import groovyx.net.http.HttpResponseException

if (!(request.'javax.servlet.error.exception' instanceof HttpResponseException)) {
    response.sendError 404
}

final HttpResponseException e = request.'javax.servlet.error.exception'

log.warning """API returned the error response:
$e.response.statusLine
${e.response.allHeaders.join('\n')}

$e.response.data"""

response.sendError 500