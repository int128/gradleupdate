package model

import groovy.transform.Canonical
import groovyx.gaelyk.datastore.Entity
import groovyx.gaelyk.datastore.Key

@Entity
@Canonical
class CurrentGradleVersion {
    @Key String label
    String version
}
