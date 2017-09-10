package org.hidetake.gradleupdate.domain

class Branch(val name: String) {
    val ref = "refs/heads/$name"
}
