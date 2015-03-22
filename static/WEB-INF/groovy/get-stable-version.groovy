import model.CurrentGradleVersion
import util.CrossOriginPolicy

CrossOriginPolicy.allowAnyOrigin(response)

final version = CurrentGradleVersion.get('stable')
if (version == null) {
    response.sendError(404)
} else {
    response.contentType = 'text/plain'
    print version.version
}
