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

    static Credential getOrCreate(String serviceName) {
        def credential = Credential.get(serviceName)
        if (credential == null) {
            credential = new Credential()
            credential.service = serviceName
            credential.save()
        }
        credential
    }

}
