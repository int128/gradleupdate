package org.hidetake.gradleupdate.view

import org.hidetake.gradleupdate.domain.GradleWrapperStatus
import org.springframework.http.HttpHeaders
import org.springframework.http.HttpStatus
import org.springframework.http.ResponseEntity

object BadgeSvg {
    fun render(status: GradleWrapperStatus?) =
        when (status?.latest) {
            true  -> render(rightMessage = status.gradleWrapper.version, rightFill = "#4c1")
            false -> render(rightMessage = status.gradleWrapper.version, rightFill = "#e05d44")
            else  -> render(rightMessage = "unknown", rightFill = "#9f9f9f")
        }

    fun notFound() = render(rightMessage = "unknown", rightFill = "#9f9f9f")

    fun render(
        leftMessage: String = "Gradle",
        leftWidth: Int = 47,
        leftFill: String = "#555",
        rightMessage: String,
        rightWidth: Int = rightMessage.length * rightMessage.length / 4 + rightMessage.length * 5 + 12,
        rightFill: String
    ) = ResponseEntity("""
            <svg xmlns="http://www.w3.org/2000/svg" width="${leftWidth + rightWidth}" height="20">
                <linearGradient id="b" x2="0" y2="100%">
                    <stop offset="0" stop-color="#bbb" stop-opacity=".1"/>
                    <stop offset="1" stop-opacity=".1"/>
                </linearGradient>
                <clipPath id="a">
                    <rect width="${leftWidth + rightWidth}" height="20" rx="3" fill="#fff"/>
                </clipPath>
                <g clip-path="url(#a)">
                    <path fill="$leftFill" d="M0 0h${leftWidth}v20H0z"/>
                    <path fill="$rightFill" d="M${leftWidth} 0h${rightWidth}v20H${leftWidth}z"/>
                    <path fill="url(#b)" d="M0 0h${leftWidth + rightWidth}v20H0z"/>
                </g>
                <g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11">
                    <text x="${leftWidth / 2.0}" y="15" fill="#010101" fill-opacity=".3">$leftMessage</text>
                    <text x="${leftWidth / 2.0}" y="14">$leftMessage</text>
                    <text x="${leftWidth + rightWidth / 2.0}" y="15" fill="#010101" fill-opacity=".3">$rightMessage</text>
                    <text x="${leftWidth + rightWidth / 2.0}" y="14">$rightMessage</text>
                </g>
            </svg>
        """.replaceIndent(),
        HttpHeaders().apply {
            add(HttpHeaders.CONTENT_TYPE, "image/svg+xml")
        },
        HttpStatus.OK
    )
}
