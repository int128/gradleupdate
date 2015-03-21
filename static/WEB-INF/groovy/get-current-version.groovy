import model.CurrentGradleVersion
import util.CrossOrigin

CrossOrigin.sendAccessControlAllowOriginForAny(response)

assert params.label

final version = CurrentGradleVersion.get(params.label)
if (version == null) {
    response.sendError(404)
}

response.contentType = 'text/plain'
print version.version
