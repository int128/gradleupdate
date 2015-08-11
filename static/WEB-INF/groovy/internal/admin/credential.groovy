import model.Credential

if (request.method == 'POST') {
    final credential = new Credential()
    params.each { k, v -> credential[k as String] = v }
    credential.save()
    response.sendRedirect(request.requestURL as String)
} else {
    html.html {
        body {
            p { a href: '/_ah/admin/datastore', '/_ah/admin/datastore' }
            form(method: 'POST') {
                ['service', 'secret'].each { key ->
                    label(key) {
                        input type: 'text', name: key
                    }
                }
                input type: 'submit'
            }
        }
    }
}
