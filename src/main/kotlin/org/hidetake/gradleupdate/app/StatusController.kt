package org.hidetake.gradleupdate.app

import org.springframework.stereotype.Controller
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.PathVariable
import org.springframework.web.servlet.ModelAndView

@Controller
class StatusController(private val service: GradleUpdateService) {
    @GetMapping("/{owner}/{repo}/status")
    fun get(@PathVariable owner: String, @PathVariable repo: String) =
        ModelAndView("status", mapOf(
            "repository" to service.getRepository("$owner/$repo"),
            "status" to service.getGradleWrapperVersionStatus("$owner/$repo"),
            "pullRequest" to service.findPullRequestForLatestGradleWrapper("$owner/$repo")
        ))
}
