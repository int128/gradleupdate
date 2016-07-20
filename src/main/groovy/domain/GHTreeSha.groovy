package domain

import groovy.transform.Immutable

@Immutable
class GHTreeSha {

    final String value

    @Override
    String toString() {
        "Tree($value)"
    }

}
