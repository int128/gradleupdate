import model.CurrentGradleVersion
import util.CrossOriginPolicy

CrossOriginPolicy.allowAnyOrigin(response)

final version = CurrentGradleVersion.get('rc')
if (version == null) {
    forward '/get-stable-version.groovy'
} else {
    response.contentType = 'text/plain'
    print version.version
}
