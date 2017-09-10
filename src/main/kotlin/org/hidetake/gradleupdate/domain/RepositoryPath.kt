package org.hidetake.gradleupdate.domain

class RepositoryPath(val owner: String, val name: String)  {
    val fullName = "$owner/$name"
}
