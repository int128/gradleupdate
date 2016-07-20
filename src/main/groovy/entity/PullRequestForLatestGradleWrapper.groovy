package entity

import groovy.transform.Canonical
import groovyx.gaelyk.datastore.Entity
import groovyx.gaelyk.datastore.Indexed
import groovyx.gaelyk.datastore.Key

@Entity
@Canonical
class PullRequestForLatestGradleWrapper {

    @Key String url
    @Indexed String fullName
    @Indexed Date createdAt
    @Indexed String fromVersion
    @Indexed String toVersion

}
