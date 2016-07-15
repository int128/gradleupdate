package domain

import groovy.transform.Immutable

@Immutable
class GradleWrapperStatus {

    final GradleVersion currentVersion

    boolean checkUpToDate() {
        currentVersion == new LatestGradle().version
    }

}
