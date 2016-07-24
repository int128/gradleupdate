package domain

import groovy.transform.Immutable

@Immutable
class GHTreeSha {

    final String value

    @Override
    String toString() {
        value.take(8)
    }

}
