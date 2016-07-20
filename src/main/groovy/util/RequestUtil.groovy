package util

import javax.servlet.http.HttpServletRequest

class RequestUtil {

    static String dirName(String uri) {
        assert uri
        uri.substring(0, uri.lastIndexOf('/'))
    }

    static String dirName(HttpServletRequest request) {
        assert request
        dirName(request.requestURI)
    }

    static String relativePath(HttpServletRequest request, String path) {
        "${dirName(request)}/$path"
    }

    static String token(Map<String, String> headers) {
        if (headers.Authorization?.startsWith('token ')) {
            headers.Authorization.substring('token '.length())
        } else {
            null
        }
    }

}
