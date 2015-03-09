import groovyx.gaelyk.datastore.Entity
import groovyx.gaelyk.datastore.Key

@Entity
class CurrentGradleVersion {
    @Key String label
    String version
}
