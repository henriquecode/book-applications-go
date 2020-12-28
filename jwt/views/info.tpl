<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Document</title>
</head>
<body>
    {{ if .Success }}
        <p style="color: red;">{{ .Message }}</p>
    {{ else }}
        <p>Usuário: {{ .Data.Name }}</p>
        <p>Idade: {{ .Data.Age }}</p>
        <p>Email: {{ .Data.Email }}</p>
        <p>{{ .Data.Book.Title }}</p>    
    {{ end }}
</body>
</html>