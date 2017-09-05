package org.hidetake.gradleupdate.app

import org.springframework.stereotype.Controller
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.PathVariable
import org.springframework.web.servlet.ModelAndView

@Controller
class StatusController(val service: GradleUpdateService) {
    @GetMapping("/{owner}/{repo}/status")
    fun get(@PathVariable owner: String, @PathVariable repo: String): ModelAndView {
        val repository = service.getRepository("$owner/$repo")
        val pullRequest = service.findPullRequestForLatestGradleWrapper("$owner/$repo")
        return ModelAndView("status", mapOf(
            "repository" to repository,
            "pullRequest" to pullRequest
        ))
    }
}
