{{ define "base" }}

<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="UTF-8" />
        <link rel="icon" href="data:,">
        <title>{{ template "title" . }}</title>
        <style>{{ template "style" . }}</style>
    </head>
    <body>

        {{ template "menu" . }}

        <hr/>

        {{ template "content" . }}

        <hr/>

        {{ template "footer" . }}

    </body>
</html>

{{ end }}