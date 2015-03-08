import groovy.transform.CompileStatic
import groovy.transform.TupleConstructor

import javax.servlet.http.HttpServletRequest
import javax.servlet.http.HttpServletResponse

@CompileStatic
@TupleConstructor
class InternalController {

    HttpServletRequest request
    HttpServletResponse response

    final FeedService service = new FeedService()

    void purgeVersionCache() {
        service.purgeCache()
    }

}
