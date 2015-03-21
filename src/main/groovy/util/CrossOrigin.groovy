package util

import javax.servlet.http.HttpServletResponse

class CrossOrigin {

    static sendAccessControlAllowOrigin(HttpServletResponse response, Map headers) {
        assert headers.Origin
        response.headers.'Access-Control-Allow-Origin' =
            headers.Origin.matches(/http:\/\/localhost(\:\d+)?/) ? headers.Origin : 'https://gradleupdate.github.io'
    }

    static sendAccessControlAllowHeaders(HttpServletResponse response, String[] headers) {
        response.headers.'Access-Control-Allow-Headers' = headers.join(',')
    }

}
