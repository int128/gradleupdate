package com.example

import groovy.transform.Canonical
import groovyx.gaelyk.datastore.Entity
import groovyx.gaelyk.datastore.Key

@Entity
@Canonical
class Item {

    @Key String id

    String name

}
