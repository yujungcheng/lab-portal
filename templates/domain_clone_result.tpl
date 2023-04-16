{{ define "content" }}
<div>

<p>
    <h3>
        => Cloned domain ... {{ len . }} domain(s).
    </h3>
</p>

{{ range $key, $value := . }}

  Domain <b>{{ $key }}</b>: {{ $value }}<br/>

{{ end }}

</div>
{{ end }}