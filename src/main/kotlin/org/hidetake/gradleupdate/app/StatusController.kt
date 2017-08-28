package org.hidetake.gradleupdate.app

import org.springframework.stereotype.Controller
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.PathVariable
import org.springframework.web.servlet.ModelAndView

@Controller
class StatusController(val service: GradleUpdateService) {
    @GetMapping("/{owner}/{repo}/status")
    fun get(@PathVariable owner: String, @PathVariable repo: String) =
        service.getRepository("$owner/$repo").let {
            ModelAndView("status", mapOf(
                "repository" to it
            ))
        }
}
