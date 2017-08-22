package org.hidetake.gradleupdate.view

class BadgeSvg(
    val leftMessage: String = "Gradle",
    val leftWidth: Int = 47,
    val leftFill: BadgeColor = BadgeColor.GREY,
    val rightMessage: String,
    val rightWidth: Int = rightMessage.length * rightMessage.length / 4 + rightMessage.length * 5 + 12,
    val rightFill: BadgeColor
) {
    fun render() = """
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
"""
}
