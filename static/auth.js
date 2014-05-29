/** Javascript front-end for go-auth */

/* Instantiate with authentication host and application realm */
function Auth(host, realm) {
    this.host = host;
    this.realm = realm;
}

/**
 * Trigger this script upon submit:
 * <form ... onsubmit="event.preventDefault(); return auth.login(this);">
 *
 * Expects to find input 'username' and 'password' in form.
 */
Auth.prototype.login = function(form) {

    var auth = this;

    var error = document.getElementById('error-msg');
    var xmlhttp = new XMLHttpRequest();

    var username = form.elements["username"].value;
    var password = form.elements["password"].value;

    xmlhttp.open("POST", auth.host+"/auth?username="+username+"&password="+password + "&realm="+auth.realm, true);

    xmlhttp.onreadystatechange = function() {
        if (xmlhttp.readyState == XMLHttpRequest.DONE) {
            if (xmlhttp.status == 200) {
                auth.setCookie(xmlhttp.responseText);
                form.submit();
            } else if(xmlhttp.status == 401) {
                error.innerHTML = "Feil brukarnamn og/eller passord?";
            } else {
                error.innerHTML = "Noko gjekk galt, meld ifrå til Kjetil!";
            }
            error.style.display = 'block';
        }
    }
    xmlhttp.send();
    return false;
}

/** Stores the received cookie */
Auth.prototype.setCookie = function(cookie) {
    document.cookie = 'token=' + cookie + '; path=/';
}
