package org.hidetake.gradleupdate.controller

import org.hidetake.gradleupdate.service.BadgeService
import org.hidetake.gradleupdate.view.BadgeColor
import org.hidetake.gradleupdate.view.BadgeSvg
import org.slf4j.LoggerFactory
import org.springframework.stereotype.Controller
import org.springframework.web.bind.annotation.*
import java.io.IOException
import javax.servlet.http.HttpServletRequest
import javax.servlet.http.HttpServletResponse

@Controller
@ResponseBody
class BadgeController(val badgeService: BadgeService) {
    private val log = LoggerFactory.getLogger(javaClass)

    @GetMapping("/{owner}/{repo}/status.svg")
    fun get(@PathVariable owner: String, @PathVariable repo: String): String =
        badgeService.findGradleWrapper("$owner/$repo")?.let {
            BadgeSvg(rightMessage = it.version, rightFill = BadgeColor.GREEN).render()
        } ?: BadgeSvg(rightMessage = "unknown", rightFill = BadgeColor.SILVER).render()

    @ExceptionHandler
    fun error(e: IOException, request: HttpServletRequest): String {
        log.warn("Error while processing ${request.requestURL}", e)
        return BadgeSvg(rightMessage = "unknown", rightFill = BadgeColor.SILVER).render()
    }

    @ModelAttribute
    fun contentType(response: HttpServletResponse) {
        response.contentType = "image/svg+xml"
    }
}
