package model

import groovyx.gaelyk.datastore.Entity
import groovyx.gaelyk.datastore.Key
import groovyx.gaelyk.datastore.Unindexed

@Entity
class Credential {
    @Key String service
    @Unindexed String token
    @Unindexed String clientId
    @Unindexed String clientSecret
}
