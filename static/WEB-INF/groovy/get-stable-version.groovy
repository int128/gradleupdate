import model.CurrentGradleVersion
import util.CrossOrigin

CrossOrigin.sendAccessControlAllowOriginForAny(response)

final version = CurrentGradleVersion.get('stable')
if (version == null) {
    response.sendError(404)
} else {
    response.contentType = 'text/plain'
    print version.version
}
