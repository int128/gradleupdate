package entity

import groovy.transform.Canonical
import groovyx.gaelyk.datastore.Entity
import groovyx.gaelyk.datastore.Indexed
import groovyx.gaelyk.datastore.Key

@Entity
@Canonical
class PullRequest {

    @Key String url
    @Indexed String htmlUrl
    @Indexed String gradleVersion
    @Indexed Date createdAt
    @Indexed String repo
    @Indexed String owner
    @Indexed int ownerId
    @Indexed boolean merged

}
