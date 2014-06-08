<!doctype html>
<html>
 <head>
  <meta charset="UTF-8">
  <!-- Latest compiled and minified CSS -->
  <link rel="stylesheet"
  href="//netdna.bootstrapcdn.com/bootstrap/3.1.1/css/bootstrap.min.css">
  <!-- Optional theme -->
  <link rel="stylesheet" href="//netdna.bootstrapcdn.com/bootstrap/3.1.1/css/bootstrap-theme.min.css">
 </head>
 <body>
  <div class="container">
   <div class="page-header">
    <h2>identity <small>services accessible with this token</small></h2>
   </div>
   <div class="panel panel-primary">
    <div class="panel-heading">
     <h3 class="panel-title">Current token</h3>
    </div>
    <div class="panel-body">
     <p>Identity: {{ .User.Username }}</p>
     <p>Expire: {{ .User.Expire }}</p>
     <p>Issued by: {{ .User.Realm }}</p>
    </div>
   </div>
   <div class="panel panel-primary">
    <div class="panel-heading">
     <h3 class="panel-title">Available services</h3>
    </div>
    <ul class="list-group">
     {{range .Realms}}
      <li class="list-group-item">{{ .Name }}</li>
     {{end}}
    </ul>
  </div>
 </body>
</html>
