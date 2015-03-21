import model.CurrentGradleVersion
import util.CrossOrigin

CrossOrigin.sendAccessControlAllowOriginForAny(response)

final version = CurrentGradleVersion.get('rc')
if (version == null) {
    forward '/get-stable-version.groovy'
} else {
    response.contentType = 'text/plain'
    print version.version
}
