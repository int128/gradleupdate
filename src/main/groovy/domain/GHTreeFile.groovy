package domain

import groovy.transform.Immutable

@Immutable
class GHTreeFile {

    String path
    String mode
    GHBlobSha sha

    Map<String, String> asMap() {
        [path: path, mode: mode, type: 'blob', sha: sha.value]
    }

}
