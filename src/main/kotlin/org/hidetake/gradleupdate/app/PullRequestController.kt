package org.hidetake.gradleupdate.app

import org.springframework.stereotype.Controller
import org.springframework.web.bind.annotation.PathVariable
import org.springframework.web.bind.annotation.PostMapping
import org.springframework.web.bind.annotation.ResponseBody

@Controller
class PullRequestController(private val gradleUpdateService: GradleUpdateService) {
    @PostMapping("/{owner}/{repo}/pullRequest")
    @ResponseBody //TODO: model and view
    fun create(@PathVariable owner: String, @PathVariable repo: String) =
        gradleUpdateService.createPullRequestForLatestGradleWrapper("$owner/$repo")
}
