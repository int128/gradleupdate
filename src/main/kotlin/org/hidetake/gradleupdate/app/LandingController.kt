package org.hidetake.gradleupdate.app

import org.springframework.stereotype.Controller
import org.springframework.validation.BindingResult
import org.springframework.validation.annotation.Validated
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.PostMapping
import org.springframework.web.bind.annotation.RequestMapping
import javax.validation.constraints.NotNull

@Controller
@RequestMapping("/landing")
class LandingController {
    @GetMapping
    fun get() = "redirect:/"

    @PostMapping
    fun post(@Validated repositoryForm: RepositoryForm, bindingResult: BindingResult) =
        when {
            bindingResult.hasErrors() -> "/"
            else -> repositoryForm.extractOwnerRepo()?.let { "$it/status" } ?: "/"
        }.let { uri ->
            "redirect:$uri"
        }

    data class RepositoryForm(@NotNull var url: String = "") {
        fun extractOwnerRepo(): String? = Regex("""[^/]+/[^/]+$""").find(url)?.value
    }
}
