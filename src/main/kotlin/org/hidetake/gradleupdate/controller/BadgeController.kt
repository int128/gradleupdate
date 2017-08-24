package org.hidetake.gradleupdate.controller

import org.hidetake.gradleupdate.service.GradleUpdateService
import org.hidetake.gradleupdate.view.BadgeSvg
import org.slf4j.LoggerFactory
import org.springframework.http.ResponseEntity
import org.springframework.stereotype.Controller
import org.springframework.web.bind.annotation.ExceptionHandler
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.PathVariable
import org.springframework.web.bind.annotation.ResponseBody
import java.io.IOException
import javax.servlet.http.HttpServletRequest

@Controller
@ResponseBody
class BadgeController(val gradleUpdateService: GradleUpdateService) {
    private val log = LoggerFactory.getLogger(javaClass)

    @GetMapping("/{owner}/{repo}/status.svg")
    fun get(@PathVariable owner: String, @PathVariable repo: String) =
        BadgeSvg.render(gradleUpdateService.getStatus("$owner/$repo"))

    @ExceptionHandler
    fun error(e: IOException, request: HttpServletRequest): ResponseEntity<String> {
        log.warn("Error while processing ${request.requestURL}", e)
        return BadgeSvg.render(null)
    }
}
