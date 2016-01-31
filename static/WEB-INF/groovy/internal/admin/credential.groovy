import model.Credential

if (request.method == 'POST') {
    assert params.key
    assert params.value
    final credential = new Credential()
    credential.service = params.key
    credential.secret = params.value
    credential.save()
    response.sendRedirect(request.requestURL as String)
} else {
    html.html {
        head {
            link(rel: 'stylesheet', href: '/bootstrap.min.css')
        }
        body(class: 'container') {
            h1('Credential Setup')

            h2('Get')
            p { a href: '/_ah/admin/datastore', '/_ah/admin/datastore' }

            h2('Save')
            Credential.CredentialKey.values().each { key ->
                form(method: 'POST') {
                    h3(key.name())
                    input type: 'hidden', name: 'key', value: key.name()
                    input type: 'text', name: 'value'
                    input type: 'submit'
                }
            }
        }
    }
}
