import groovy.transform.CompileStatic
import groovy.transform.TupleConstructor

import javax.servlet.http.HttpServletRequest
import javax.servlet.http.HttpServletResponse

@CompileStatic
@TupleConstructor
class FeedController {

    HttpServletRequest request
    HttpServletResponse response

    final FeedService service = new FeedService()

    void stableVersions() {
        response.setContentType('text/xml')
        response.setCharacterEncoding('UTF-8')
        response.writer << service.stableVersions
    }

    void allVersions() {
        response.setContentType('text/xml')
        response.setCharacterEncoding('UTF-8')
        response.writer << service.allVersions
    }

}
