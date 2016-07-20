package entity

import groovy.transform.Canonical
import groovyx.gaelyk.datastore.Entity
import groovyx.gaelyk.datastore.Indexed
import groovyx.gaelyk.datastore.Key
import groovyx.gaelyk.datastore.Unindexed

@Entity
@Canonical
class PullRequestForLatestGradleWrapperTaskState {

    @Key String fullName
    @Indexed State state
    @Indexed Date lastUpdated
    @Unindexed String message

    static enum State {
        Queued,
        Updating,
        AlreadyLatest,
        Done,
    }

}
