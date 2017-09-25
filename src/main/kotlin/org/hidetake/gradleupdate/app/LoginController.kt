package org.hidetake.gradleupdate.app

import org.springframework.stereotype.Controller
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.RequestParam
import org.springframework.web.servlet.ModelAndView
import org.springframework.web.servlet.support.ServletUriComponentsBuilder

@Controller
class LoginController(private val service: GitHubOAuthService) {
    @GetMapping("/login")
    fun authorize() = ModelAndView(
        "redirect:${service.getRedirectURL()}",
        service.createAuthorizationParameters(
            ServletUriComponentsBuilder.fromCurrentRequest()
                .replacePath("/login/auth")
                .replaceQuery(null)
                .build().toUriString()))

    @GetMapping("/login/auth")
    fun accessToken(@RequestParam state: String, @RequestParam code: String): String {
        service.continueAuthorization(state, code)
        return "redirect:/my"
    }
}
