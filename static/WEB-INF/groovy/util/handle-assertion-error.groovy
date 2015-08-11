final e = request.'javax.servlet.error.exception'
assert e instanceof AssertionError

response.sendError(400)