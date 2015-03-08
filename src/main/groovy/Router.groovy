import groovy.transform.CompileStatic
import groovy.util.logging.Log

import javax.servlet.http.HttpServlet
import javax.servlet.http.HttpServletRequest
import javax.servlet.http.HttpServletResponse

@Log
@CompileStatic
class Router extends HttpServlet {
    void doGet(HttpServletRequest request, HttpServletResponse response) {
        switch (request.pathInfo) {
            case '/feed':
                new FeedController(request, response).stableVersions()
                break

            default:
                super.doGet(request, response)
        }
    }
}
