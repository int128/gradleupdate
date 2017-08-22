package org.hidetake.gradleupdate.view

enum class BadgeColor(val code: String) {
    GREEN("#4c1"),
    RED("#e05d44"),
    SILVER("#9f9f9f"),
    GREY("#555");

    override fun toString(): String = code
}
