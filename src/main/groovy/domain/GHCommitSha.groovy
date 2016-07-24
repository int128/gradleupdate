package domain

import groovy.transform.Immutable

@Immutable
class GHCommitSha {

    final String value

    @Override
    String toString() {
        value.take(8)
    }

}
