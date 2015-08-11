import groovyx.net.http.HttpResponseException

final e = request.'javax.servlet.error.exception'
assert e instanceof HttpResponseException

e.printStackTrace()

log.info """${e.response.statusLine}
[Headers]
${e.response.allHeaders.join('\n')}
[Body]
${e.response.data}"""

response.sendError(500)