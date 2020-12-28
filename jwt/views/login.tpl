<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Document</title>
</head>
<body>
    <form action="login" method="post">
        <input type="email" name="user" placeholder="Usuário" />
        <input type="password" name="password" placeholder="Senha" />
        <button>Logar</button>
        {{ if not .Success }}
            <p style="color: red;">{{ .Message }}</p>
        {{ end }}
    </form>
</body>
</html>