package infrastructure

import wslite.rest.Response

class RestAPIException extends RuntimeException {

    final Response response

    def RestAPIException(Response response) {
        super(formatResponse(response))
        this.response = response
    }

    static String formatResponse(Response response) {
        def request = response.request
        """$response.statusCode $response.statusMessage

$request.method $request.url
$request.contentAsString

${response.headers.collect { k, v -> "$k: $v" }.join('\n')}
$response.contentAsString
"""
    }

}
