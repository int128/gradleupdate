package org.hidetake.gradleupdate.app

import org.springframework.stereotype.Controller
import org.springframework.web.bind.annotation.PathVariable
import org.springframework.web.bind.annotation.PostMapping

@Controller
class PullRequestController(private val gradleUpdateService: GradleUpdateService) {
    @PostMapping("/{owner}/{repo}/pullRequest")
    fun create(@PathVariable owner: String, @PathVariable repo: String): String {
        gradleUpdateService.createPullRequestForLatestGradleWrapper("$owner/$repo")
        return "redirect:/$owner/$repo/status"
    }
}
