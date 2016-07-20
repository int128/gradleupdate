package domain

import groovy.transform.Immutable

@Immutable
class GHBlobSha {

    final String value

    @Override
    String toString() {
        value
    }

}
