package org.hidetake.gradleupdate.controller

import org.hidetake.gradleupdate.service.BadgeService
import org.slf4j.LoggerFactory
import org.springframework.stereotype.Controller
import org.springframework.ui.Model
import org.springframework.web.bind.annotation.ExceptionHandler
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.PathVariable
import java.io.IOException

@Controller
class BadgeController(val badgeService: BadgeService) {
    private val log = LoggerFactory.getLogger(javaClass)

    @GetMapping("/{owner}/{repo}/status.svg")
    fun get(@PathVariable owner: String, @PathVariable repo: String, model: Model): String {
        val gradleWrapper = badgeService.findGradleWrapper("$owner/$repo")
        model.addAttribute("status", gradleWrapper?.version)
        return "badge"
    }

    @ExceptionHandler
    fun error(ioException: IOException, model: Model): String {
        log.warn("Error while generating badge", ioException)
        model.addAttribute("status", "unknown")
        return "badge"
    }
}
