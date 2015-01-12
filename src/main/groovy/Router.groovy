import groovy.transform.CompileStatic
import groovy.util.logging.Log

import javax.servlet.http.HttpServlet
import javax.servlet.http.HttpServletRequest
import javax.servlet.http.HttpServletResponse

@Log
@CompileStatic
class Router extends HttpServlet {
    void doGet(HttpServletRequest request, HttpServletResponse response) {
        log.info("GET request to ${request.pathInfo}")
    }

    void doPost(HttpServletRequest request, HttpServletResponse response) {
        log.info("POST request to ${request.pathInfo}")
    }
}
