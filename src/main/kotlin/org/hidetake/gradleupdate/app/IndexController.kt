package org.hidetake.gradleupdate.app

import org.springframework.stereotype.Controller
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.servlet.ModelAndView
import java.security.Principal

@Controller
@RequestMapping("/")
class IndexController(private val loginUserService: LoginUserService) {
    @GetMapping
    fun get(principal: Principal?) = when (principal) {
        null -> ModelAndView("index")
        else -> ModelAndView("my",
            mapOf(
                "repositories" to loginUserService.getRepositories()
            )
        )
    }
}
