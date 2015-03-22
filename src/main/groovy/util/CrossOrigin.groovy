package util

import javax.servlet.http.HttpServletResponse

class CrossOrigin {

    static sendAccessControlAllowOrigin(HttpServletResponse response, Map<String, String> headers) {
        def origin = headers.find { it.key.equalsIgnoreCase('origin') }?.value
        response.headers.'Access-Control-Allow-Origin' =
            origin?.matches(/http:\/\/localhost(\:\d+)?/) ? origin : 'https://gradleupdate.github.io'
    }

    static sendAccessControlAllowOriginForAny(HttpServletResponse response) {
        response.headers.'Access-Control-Allow-Origin' = '*'
    }

    static sendAccessControlAllowHeaders(HttpServletResponse response, String[] headers) {
        response.headers.'Access-Control-Allow-Headers' = headers.join(',')
    }

}
