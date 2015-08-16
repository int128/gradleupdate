final e = request.'javax.servlet.error.exception'
assert e instanceof AssertionError

if (headers.'X-AppEngine-TaskName') {
    log.severe('Assertion error occurred in the task, maybe a bug')
    response.sendError(500)
} else {
    response.sendError(400)
}
